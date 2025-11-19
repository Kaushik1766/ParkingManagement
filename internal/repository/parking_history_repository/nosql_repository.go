package parkinghistoryrepository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Kaushik1766/ParkingManagement/internal/config"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type NOSQLParkingRepository struct {
	client *dynamodb.Client
}

func NewNOSQLParkingRepository(client *dynamodb.Client) *NOSQLParkingRepository {
	return &NOSQLParkingRepository{
		client: client,
	}
}

func (nosqlpr *NOSQLParkingRepository) AddParking(ctx context.Context, vehicle models.Vehicle) (string, error) {
	// check if slot already occupied
	slotKey := fmt.Sprintf("FLOOR#%d#SLOT#%d", vehicle.AssignedFloorNumber, vehicle.AssignedSlotNumber)
	slotRes, err := nosqlpr.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", vehicle.AssignedBuildingID.String())},
			"SK": &types.AttributeValueMemberS{Value: slotKey},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return "", errors.New("error checking slot availability")
	}

	if slotRes.Item != nil && slotRes.Item["IsOccupied"] != nil {
		if slotRes.Item["IsOccupied"].(*types.AttributeValueMemberBOOL).Value {
			log.Println("slot is already occupied")
			return "", errors.New("parkingrepo: vehicle is already parked")
		}
	}

	parkingID := uuid.New()
	timestamp := time.Now().Unix()

	item := map[string]types.AttributeValue{
		"PK":          &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", vehicle.User.Email)},
		"SK":          &types.AttributeValueMemberS{Value: fmt.Sprintf("PARKING#%d", timestamp)},
		"ParkingId":   &types.AttributeValueMemberS{Value: parkingID.String()},
		"Numberplate": &types.AttributeValueMemberS{Value: vehicle.NumberPlate},
		"BuildingId":  &types.AttributeValueMemberS{Value: vehicle.AssignedBuildingID.String()},
		"FloorNumber": &types.AttributeValueMemberN{Value: strconv.Itoa(vehicle.AssignedFloorNumber)},
		"SlotId":      &types.AttributeValueMemberN{Value: strconv.Itoa(vehicle.AssignedSlotNumber)},
		"StartTime":   &types.AttributeValueMemberN{Value: strconv.FormatInt(timestamp, 10)},
		"VehicleType": &types.AttributeValueMemberS{Value: vehicle.VehicleType.String()},
	}

	_, err = nosqlpr.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Item:      item,
	})
	if err != nil {
		log.Println(err.Error())
		return "", errors.New("error adding parking history")
	}

	// Update IsParked status in vehicle entity
	_, err = nosqlpr.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", vehicle.User.Email)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("VEHICLE#%s", vehicle.NumberPlate)},
		},
		UpdateExpression: aws.String("SET IsParked = :isParked"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":isParked": &types.AttributeValueMemberBOOL{Value: true},
		},
	})
	if err != nil {
		log.Println("Warning: could not update IsParked status:", err.Error())
	}

	return parkingID.String(), nil
}

func (nosqlpr *NOSQLParkingRepository) Unpark(ctx context.Context, id string) error {
	var userEmail string
	var timestamp string

	scanRes, err := nosqlpr.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(config.DynamoDBTable),
		FilterExpression: aws.String("ParkingId = :parkingId AND attribute_not_exists(EndTime)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":parkingId": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return errors.New("error finding parking record")
	}

	if len(scanRes.Items) == 0 {
		log.Println("parking record not found or already unparked in Unpark")
		return errors.New("parking record not found or already unparked")
	}

	item := scanRes.Items[0]
	userEmail = item["PK"].(*types.AttributeValueMemberS).Value
	timestamp = item["SK"].(*types.AttributeValueMemberS).Value

	endTime := time.Now().Unix()
	_, err = nosqlpr.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: userEmail},
			"SK": &types.AttributeValueMemberS{Value: timestamp},
		},
		UpdateExpression: aws.String("SET EndTime = :endTime"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":endTime": &types.AttributeValueMemberN{Value: strconv.FormatInt(endTime, 10)},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return errors.New("error updating parking record")
	}

	// Update IsParked status in vehicle entity
	numberplate := item["Numberplate"].(*types.AttributeValueMemberS).Value
	_, err = nosqlpr.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: userEmail},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("VEHICLE#%s", numberplate)},
		},
		UpdateExpression: aws.String("SET IsParked = :isParked"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":isParked": &types.AttributeValueMemberBOOL{Value: false},
		},
	})
	if err != nil {
		log.Println("Warning: could not update IsParked status:", err.Error())
	}

	return nil
}

func (nosqlpr *NOSQLParkingRepository) GetParkingHistoryByNumberPlate(ctx context.Context, numberplate string, startTime, endTime time.Time) ([]models.ParkingHistoryDTO, error) {
	var history []models.ParkingHistoryDTO

	startTimestamp := startTime.Unix()
	endTimestamp := endTime.Unix()

	scanRes, err := nosqlpr.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(config.DynamoDBTable),
		FilterExpression: aws.String("Numberplate = :numberplate AND StartTime >= :startTime AND EndTime <= :endTime AND attribute_exists(EndTime)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":numberplate": &types.AttributeValueMemberS{Value: numberplate},
			":startTime":   &types.AttributeValueMemberN{Value: strconv.FormatInt(startTimestamp, 10)},
			":endTime":     &types.AttributeValueMemberN{Value: strconv.FormatInt(endTimestamp, 10)},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("error fetching parking history")
	}

	for _, item := range scanRes.Items {
		var dto models.ParkingHistoryDTO
		dto.TicketId = item["ParkingId"].(*types.AttributeValueMemberS).Value
		dto.NumberPlate = item["Numberplate"].(*types.AttributeValueMemberS).Value
		dto.BuildingId = item["BuildingId"].(*types.AttributeValueMemberS).Value
		dto.FLoorNumber, _ = strconv.Atoi(item["FloorNumber"].(*types.AttributeValueMemberN).Value)
		dto.SlotNumber, _ = strconv.Atoi(item["SlotId"].(*types.AttributeValueMemberN).Value)
		dto.VechicleType = item["VehicleType"].(*types.AttributeValueMemberS).Value

		startTimeUnix, _ := strconv.ParseInt(item["StartTime"].(*types.AttributeValueMemberN).Value, 10, 64)
		dto.StartTime = time.Unix(startTimeUnix, 0).Local()

		if item["EndTime"] != nil {
			endTimeUnix, _ := strconv.ParseInt(item["EndTime"].(*types.AttributeValueMemberN).Value, 10, 64)
			dto.EndTime = time.Unix(endTimeUnix, 0).Local()
		}

		history = append(history, dto)
	}

	return history, nil
}

func (nosqlpr *NOSQLParkingRepository) GetParkingHistoryByUser(ctx context.Context, userId string, startTime, endTime time.Time) ([]models.ParkingHistoryDTO, error) {
	var history []models.ParkingHistoryDTO

	// First, get user email from userId
	userRes, err := nosqlpr.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userId)},
			":sk": &types.AttributeValueMemberS{Value: "PROFILE#"},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("error fetching user")
	}

	if len(userRes.Items) == 0 {
		log.Println("user not found in GetParkingHistoryByUser")
		return nil, errors.New("user not found")
	}

	userEmail := userRes.Items[0]["Email"].(*types.AttributeValueMemberS).Value

	// Query parking history
	startTimestamp := startTime.Unix()
	endTimestamp := endTime.Unix()

	queryRes, err := nosqlpr.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk AND SK BETWEEN :startSK AND :endSK"),
		FilterExpression:       aws.String("attribute_exists(EndTime)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":      &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userEmail)},
			":startSK": &types.AttributeValueMemberS{Value: fmt.Sprintf("PARKING#%d", startTimestamp)},
			":endSK":   &types.AttributeValueMemberS{Value: fmt.Sprintf("PARKING#%d", endTimestamp)},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("error fetching parking history")
	}

	for _, item := range queryRes.Items {
		var dto models.ParkingHistoryDTO
		dto.TicketId = item["ParkingId"].(*types.AttributeValueMemberS).Value
		dto.NumberPlate = item["Numberplate"].(*types.AttributeValueMemberS).Value
		dto.BuildingId = item["BuildingId"].(*types.AttributeValueMemberS).Value
		dto.FLoorNumber, _ = strconv.Atoi(item["FloorNumber"].(*types.AttributeValueMemberN).Value)
		dto.SlotNumber, _ = strconv.Atoi(item["SlotId"].(*types.AttributeValueMemberN).Value)
		dto.VechicleType = item["VehicleType"].(*types.AttributeValueMemberS).Value

		startTimeUnix, _ := strconv.ParseInt(item["StartTime"].(*types.AttributeValueMemberN).Value, 10, 64)
		dto.StartTime = time.Unix(startTimeUnix, 0).Local()

		if item["EndTime"] != nil {
			endTimeUnix, _ := strconv.ParseInt(item["EndTime"].(*types.AttributeValueMemberN).Value, 10, 64)
			dto.EndTime = time.Unix(endTimeUnix, 0).Local()
		}

		// Get building name
		buildingRes, err := nosqlpr.client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(config.DynamoDBTable),
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: "BUILDING"},
				"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", dto.BuildingId)},
			},
		})
		if err == nil && buildingRes.Item != nil {
			dto.BuildingName = buildingRes.Item["BuildingName"].(*types.AttributeValueMemberS).Value
		}

		history = append(history, dto)
	}

	return history, nil
}

func (nosqlpr *NOSQLParkingRepository) GetActiveUserParkings(ctx context.Context, userId string) ([]models.ParkingHistoryDTO, error) {
	var activeParkings []models.ParkingHistoryDTO

	// First, get user email from userId
	userRes, err := nosqlpr.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userId)},
			":sk": &types.AttributeValueMemberS{Value: "PROFILE#"},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("error fetching user")
	}

	if len(userRes.Items) == 0 {
		log.Println("user not found in GetActiveUserParkings")
		return nil, errors.New("user not found")
	}

	userEmail := userRes.Items[0]["Email"].(*types.AttributeValueMemberS).Value

	// Query active parkings (without EndTime)
	queryRes, err := nosqlpr.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		FilterExpression:       aws.String("attribute_not_exists(EndTime)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userEmail)},
			":sk": &types.AttributeValueMemberS{Value: "PARKING#"},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("error fetching active parkings")
	}

	for _, item := range queryRes.Items {
		var dto models.ParkingHistoryDTO
		dto.TicketId = item["ParkingId"].(*types.AttributeValueMemberS).Value
		dto.NumberPlate = item["Numberplate"].(*types.AttributeValueMemberS).Value
		dto.BuildingId = item["BuildingId"].(*types.AttributeValueMemberS).Value
		dto.FLoorNumber, _ = strconv.Atoi(item["FloorNumber"].(*types.AttributeValueMemberN).Value)
		dto.SlotNumber, _ = strconv.Atoi(item["SlotId"].(*types.AttributeValueMemberN).Value)
		dto.VechicleType = item["VehicleType"].(*types.AttributeValueMemberS).Value

		startTimeUnix, _ := strconv.ParseInt(item["StartTime"].(*types.AttributeValueMemberN).Value, 10, 64)
		dto.StartTime = time.Unix(startTimeUnix, 0).Local()

		activeParkings = append(activeParkings, dto)
	}

	return activeParkings, nil
}

func (nosqlpr *NOSQLParkingRepository) UnparkByNumberPlate(ctx context.Context, numberplate string) error {
	// Find active parking by numberplate
	scanRes, err := nosqlpr.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(config.DynamoDBTable),
		FilterExpression: aws.String("Numberplate = :numberplate AND attribute_not_exists(EndTime)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":numberplate": &types.AttributeValueMemberS{Value: numberplate},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return errors.New("error finding parking record")
	}

	if len(scanRes.Items) == 0 {
		log.Println("no active parking found for numberplate:", numberplate)
		return errors.New("no active parking found for this numberplate")
	}

	item := scanRes.Items[0]
	userEmail := item["PK"].(*types.AttributeValueMemberS).Value
	timestamp := item["SK"].(*types.AttributeValueMemberS).Value

	// Update with end time
	endTime := time.Now().Unix()
	_, err = nosqlpr.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: userEmail},
			"SK": &types.AttributeValueMemberS{Value: timestamp},
		},
		UpdateExpression: aws.String("SET EndTime = :endTime"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":endTime": &types.AttributeValueMemberN{Value: strconv.FormatInt(endTime, 10)},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return errors.New("error updating parking record")
	}

	// Update IsParked status in vehicle entity
	_, err = nosqlpr.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: userEmail},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("VEHICLE#%s", numberplate)},
		},
		UpdateExpression: aws.String("SET IsParked = :isParked"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":isParked": &types.AttributeValueMemberBOOL{Value: false},
		},
	})
	if err != nil {
		log.Println("Warning: could not update IsParked status:", err.Error())
	}

	return nil
}
