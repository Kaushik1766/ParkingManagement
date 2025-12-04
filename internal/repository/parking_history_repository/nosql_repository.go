package parkinghistoryrepository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/Kaushik1766/ParkingManagement/internal/config"
	"github.com/Kaushik1766/ParkingManagement/internal/constants"
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
	userCtx := ctx.Value(constants.User).(models.UserJwt)
	userId := userCtx.ID

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

	// Use TransactWriteItems for atomic operation
	transactItems := []types.TransactWriteItem{
		// 1. Put parking history record
		{
			Put: &types.Put{
				TableName: aws.String(config.DynamoDBTable),
				Item: map[string]types.AttributeValue{
					"PK":          &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userId)},
					"SK":          &types.AttributeValueMemberS{Value: fmt.Sprintf("PARKING#%d", timestamp)},
					"ParkingId":   &types.AttributeValueMemberS{Value: parkingID.String()},
					"Numberplate": &types.AttributeValueMemberS{Value: vehicle.NumberPlate},
					"BuildingId":  &types.AttributeValueMemberS{Value: vehicle.AssignedBuildingID.String()},
					"FloorNumber": &types.AttributeValueMemberN{Value: strconv.Itoa(vehicle.AssignedFloorNumber)},
					"SlotId":      &types.AttributeValueMemberN{Value: strconv.Itoa(vehicle.AssignedSlotNumber)},
					"StartTime":   &types.AttributeValueMemberN{Value: strconv.FormatInt(timestamp, 10)},
					"VehicleType": &types.AttributeValueMemberS{Value: vehicle.VehicleType.String()},
				},
			},
		},
		// 2. Update vehicle IsParked status
		{
			Update: &types.Update{
				TableName: aws.String(config.DynamoDBTable),
				Key: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userId)},
					"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("VEHICLE#%s", vehicle.NumberPlate)},
				},
				UpdateExpression: aws.String("SET IsParked = :isParked"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":isParked": &types.AttributeValueMemberBOOL{Value: true},
				},
			},
		},
		// 3. Decrement floor AvailableSlots
		{
			Update: &types.Update{
				TableName: aws.String(config.DynamoDBTable),
				Key: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", vehicle.AssignedBuildingID.String())},
					"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("FLOORINFO#%d", vehicle.AssignedFloorNumber)},
				},
				UpdateExpression: aws.String("SET AvailableSlots = AvailableSlots - :decrement"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":decrement": &types.AttributeValueMemberN{Value: "1"},
				},
			},
		},
		// 4. Decrement building AvailableSlots
		{
			Update: &types.Update{
				TableName: aws.String(config.DynamoDBTable),
				Key: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: "BUILDING"},
					"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", vehicle.AssignedBuildingID.String())},
				},
				UpdateExpression: aws.String("SET AvailableSlots = AvailableSlots - :decrement"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":decrement": &types.AttributeValueMemberN{Value: "1"},
				},
			},
		},
		// 5. Mark slot as occupied
		{
			Update: &types.Update{
				TableName: aws.String(config.DynamoDBTable),
				Key: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", vehicle.AssignedBuildingID.String())},
					"SK": &types.AttributeValueMemberS{Value: slotKey},
				},
				UpdateExpression: aws.String("SET IsOccupied = :isOccupied"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":isOccupied": &types.AttributeValueMemberBOOL{Value: true},
				},
			},
		},
	}

	_, err = nosqlpr.client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	})
	if err != nil {
		log.Println("Error in transaction:", err.Error())
		return "", errors.New("error adding parking history")
	}

	return parkingID.String(), nil
}

func (nosqlpr *NOSQLParkingRepository) Unpark(ctx context.Context, ticketId string) error {
	userCtx := ctx.Value(constants.User).(models.UserJwt)

	var PK string
	var SK string

	// get the parking record
	queryRes, err := nosqlpr.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		FilterExpression:       aws.String("ParkingId = :parkingId AND attribute_not_exists(EndTime)"),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":parkingId": &types.AttributeValueMemberS{Value: ticketId},
			":pk":        &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userCtx.Email)},
			":sk":        &types.AttributeValueMemberS{Value: "PARKING#"},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return errors.New("error finding parking record")
	}

	if len(queryRes.Items) == 0 {
		log.Println("parking record not found or already unparked in Unpark")
		return errors.New("parking record not found or already unparked")
	}

	item := queryRes.Items[0]
	PK = item["PK"].(*types.AttributeValueMemberS).Value
	SK = item["SK"].(*types.AttributeValueMemberS).Value
	buildingID := item["BuildingId"].(*types.AttributeValueMemberS).Value
	floorNumber, _ := strconv.Atoi(item["FloorNumber"].(*types.AttributeValueMemberN).Value)
	slotId, _ := strconv.Atoi(item["SlotId"].(*types.AttributeValueMemberN).Value)

	endTime := time.Now().Unix()
	numberplate := item["Numberplate"].(*types.AttributeValueMemberS).Value

	transactItems := []types.TransactWriteItem{
		// 1. Update parking record with EndTime
		{
			Update: &types.Update{
				TableName: aws.String(config.DynamoDBTable),
				Key: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: PK},
					"SK": &types.AttributeValueMemberS{Value: SK},
				},
				UpdateExpression: aws.String("SET EndTime = :endTime"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":endTime": &types.AttributeValueMemberN{Value: strconv.FormatInt(endTime, 10)},
				},
			},
		},
		// 2. Update vehicle IsParked status
		{
			Update: &types.Update{
				TableName: aws.String(config.DynamoDBTable),
				Key: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: PK},
					"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("VEHICLE#%s", numberplate)},
				},
				UpdateExpression: aws.String("SET IsParked = :isParked"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":isParked": &types.AttributeValueMemberBOOL{Value: false},
				},
			},
		},
		// 3. Increment floor AvailableSlots
		{
			Update: &types.Update{
				TableName: aws.String(config.DynamoDBTable),
				Key: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", buildingID)},
					"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("FLOORINFO#%d", floorNumber)},
				},
				UpdateExpression: aws.String("SET AvailableSlots = AvailableSlots + :increment"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":increment": &types.AttributeValueMemberN{Value: "1"},
				},
			},
		},
		// 4. Increment building AvailableSlots
		{
			Update: &types.Update{
				TableName: aws.String(config.DynamoDBTable),
				Key: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: "BUILDING"},
					"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", buildingID)},
				},
				UpdateExpression: aws.String("SET AvailableSlots = AvailableSlots + :increment"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":increment": &types.AttributeValueMemberN{Value: "1"},
				},
			},
		},
		// 5. Mark slot as unoccupied
		{
			Update: &types.Update{
				TableName: aws.String(config.DynamoDBTable),
				Key: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", buildingID)},
					"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("FLOOR#%d#SLOT#%d", floorNumber, slotId)},
				},
				UpdateExpression: aws.String("SET IsOccupied = :isOccupied"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":isOccupied": &types.AttributeValueMemberBOOL{Value: false},
				},
			},
		},
	}

	_, err = nosqlpr.client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	})
	if err != nil {
		log.Println("Error in transaction:", err.Error())
		return errors.New("error unparking vehicle")
	}

	return nil
}

func (nosqlpr *NOSQLParkingRepository) GetParkingHistoryByNumberPlate(ctx context.Context, numberplate string, startTime, endTime time.Time) ([]models.ParkingHistoryDTO, error) {
	userCtx := ctx.Value(constants.User).(models.UserJwt)

	history := []models.ParkingHistoryDTO{}

	startTimestamp := startTime.Unix()
	endTimestamp := endTime.Unix()

	queryRes, err := nosqlpr.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		FilterExpression:       aws.String("Numberplate = :numberplate AND StartTime >= :startTime AND EndTime <= :endTime AND attribute_exists(EndTime)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":numberplate": &types.AttributeValueMemberS{Value: numberplate},
			":startTime":   &types.AttributeValueMemberN{Value: strconv.FormatInt(startTimestamp, 10)},
			":endTime":     &types.AttributeValueMemberN{Value: strconv.FormatInt(endTimestamp, 10)},
			":pk":          &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userCtx.Email)},
			":sk":          &types.AttributeValueMemberS{Value: "PARKING#"},
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

		history = append(history, dto)
	}

	return history, nil
}

func (nosqlpr *NOSQLParkingRepository) GetParkingHistoryByUser(ctx context.Context, userId string, startTime, endTime time.Time) ([]models.ParkingHistoryDTO, error) {
	history := []models.ParkingHistoryDTO{}

	startTimestamp := startTime.Unix()
	endTimestamp := endTime.Unix()

	queryRes, err := nosqlpr.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userId)},
			":sk": &types.AttributeValueMemberS{Value: "PARKING#"},
		},
	})
	if err != nil {
		log.Println("Error fetching parking history:", err.Error())
		return history, errors.New("error fetching parking history")
	}

	log.Printf("Found %d parking records for user %s", len(queryRes.Items), userId)
	log.Printf("Filtering parking records for time range: %d to %d (%v to %v)", startTimestamp, endTimestamp, startTime, endTime)

	for _, item := range queryRes.Items {
		if item["EndTime"] == nil {
			log.Printf("Skipping parking record without EndTime (active parking)")
			continue
		}

		if item["StartTime"] != nil {
			itemStartTime, _ := strconv.ParseInt(item["StartTime"].(*types.AttributeValueMemberN).Value, 10, 64)
			if itemStartTime < startTimestamp || itemStartTime > endTimestamp {
				log.Printf("Skipping parking record with StartTime %d (outside range %d-%d)", itemStartTime, startTimestamp, endTimestamp)
				continue
			}
		}

		var dto models.ParkingHistoryDTO
		dto.TicketId = item["ParkingId"].(*types.AttributeValueMemberS).Value
		dto.NumberPlate = item["Numberplate"].(*types.AttributeValueMemberS).Value
		dto.BuildingId = item["BuildingId"].(*types.AttributeValueMemberS).Value
		dto.FLoorNumber, _ = strconv.Atoi(item["FloorNumber"].(*types.AttributeValueMemberN).Value)
		dto.SlotNumber, _ = strconv.Atoi(item["SlotId"].(*types.AttributeValueMemberN).Value)
		dto.VechicleType = item["VehicleType"].(*types.AttributeValueMemberS).Value

		startTimeUnix, _ := strconv.ParseInt(item["StartTime"].(*types.AttributeValueMemberN).Value, 10, 64)
		dto.StartTime = time.Unix(startTimeUnix, 0).Local()

		endTimeUnix, _ := strconv.ParseInt(item["EndTime"].(*types.AttributeValueMemberN).Value, 10, 64)
		dto.EndTime = time.Unix(endTimeUnix, 0).Local()

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
	sort.Slice(history, func(i, j int) bool {
		return history[i].StartTime.After(history[j].StartTime)
	})

	return history, nil
}

func (nosqlpr *NOSQLParkingRepository) UnparkByNumberPlate(ctx context.Context, numberplate string) error {
	userCtx := ctx.Value(constants.User).(models.UserJwt)

	queryRes, err := nosqlpr.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		FilterExpression:       aws.String("Numberplate = :numberplate AND attribute_not_exists(EndTime)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":numberplate": &types.AttributeValueMemberS{Value: numberplate},
			":pk":          &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userCtx.ID)},
			":sk":          &types.AttributeValueMemberS{Value: "PARKING#"},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return errors.New("error finding parking record")
	}

	if len(queryRes.Items) == 0 {
		log.Println("no active parking found for numberplate:", numberplate)
		return errors.New("no active parking found for this numberplate")
	}

	item := queryRes.Items[0]
	pk := item["PK"].(*types.AttributeValueMemberS).Value
	timestamp := item["SK"].(*types.AttributeValueMemberS).Value
	buildingID := item["BuildingId"].(*types.AttributeValueMemberS).Value
	floorNumber, _ := strconv.Atoi(item["FloorNumber"].(*types.AttributeValueMemberN).Value)

	endTime := time.Now().Unix()
	_, err = nosqlpr.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
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

	_, err = nosqlpr.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
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

	// Update floor AvailableSlots (increment by 1)
	_, err = nosqlpr.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", buildingID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("FLOORINFO#%d", floorNumber)},
		},
		UpdateExpression: aws.String("SET AvailableSlots = AvailableSlots + :increment"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":increment": &types.AttributeValueMemberN{Value: "1"},
		},
	})
	if err != nil {
		log.Println("Warning: could not update floor AvailableSlots:", err.Error())
	}

	// Update building AvailableSlots (increment by 1)
	_, err = nosqlpr.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "BUILDING"},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", buildingID)},
		},
		UpdateExpression: aws.String("SET AvailableSlots = AvailableSlots + :increment"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":increment": &types.AttributeValueMemberN{Value: "1"},
		},
	})
	if err != nil {
		log.Println("Warning: could not update building AvailableSlots:", err.Error())
	}

	return nil
}
