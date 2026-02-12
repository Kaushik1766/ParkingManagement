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
	slotKey := fmt.Sprintf("%s%d#%s%d", constants.PrefixFloor, vehicle.AssignedFloorNumber, constants.PrefixSlot, vehicle.AssignedSlotNumber)
	slotRes, err := nosqlpr.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PrefixBuilding, vehicle.AssignedBuildingID.String())},
			"SK": &types.AttributeValueMemberS{Value: slotKey},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return "", errors.New(constants.ErrCheckingSlotAvailability)
	}

	if slotRes.Item != nil && slotRes.Item["IsOccupied"] != nil {
		if slotRes.Item["IsOccupied"].(*types.AttributeValueMemberBOOL).Value {
			log.Println("slot is already occupied")
			return "", errors.New(constants.ErrSlotAlreadyOccupied)
		}
	}

	parkingID := uuid.New()
	timestamp := time.Now().Unix()

	// fetch user profile to get Username and Email
	userRes, err := nosqlpr.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PrefixUser, userId)},
			":sk": &types.AttributeValueMemberS{Value: constants.PrefixProfile},
		},
	})
	if err != nil {
		log.Println("Error fetching user profile:", err.Error())
		return "", errors.New(constants.ErrFetchingUserProfile)
	}
	if len(userRes.Items) == 0 {
		return "", errors.New(constants.ErrUserProfileNotFound)
	}
	username := userRes.Items[0]["Username"].(*types.AttributeValueMemberS).Value
	email := userRes.Items[0]["Email"].(*types.AttributeValueMemberS).Value

	// use TransactWriteItems for atomic operation
	transactItems := []types.TransactWriteItem{
		// 1 Put parking history record
		{
			Put: &types.Put{
				TableName: aws.String(config.DynamoDBTable),
				Item: map[string]types.AttributeValue{
					"PK":          &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PrefixUser, userId)},
					"SK":          &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%d", constants.PKParking, timestamp)},
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
		// 2 Update vehicle IsParked status
		{
			Update: &types.Update{
				TableName: aws.String(config.DynamoDBTable),
				Key: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PrefixUser, userId)},
					"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PKVehicle, vehicle.NumberPlate)},
				},
				UpdateExpression: aws.String("SET IsParked = :isParked"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":isParked": &types.AttributeValueMemberBOOL{Value: true},
				},
			},
		},
		// 3 Decrement floor AvailableSlots
		{
			Update: &types.Update{
				TableName: aws.String(config.DynamoDBTable),
				Key: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PrefixBuilding, vehicle.AssignedBuildingID.String())},
					"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%d", constants.PrefixFloorInfo, vehicle.AssignedFloorNumber)},
				},
				UpdateExpression: aws.String("SET AvailableSlots = AvailableSlots - :decrement"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":decrement": &types.AttributeValueMemberN{Value: "1"},
				},
			},
		},
		// 4 Decrement building AvailableSlots
		{
			Update: &types.Update{
				TableName: aws.String(config.DynamoDBTable),
				Key: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: constants.PKBuilding},
					"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PrefixBuilding, vehicle.AssignedBuildingID.String())},
				},
				UpdateExpression: aws.String("SET AvailableSlots = AvailableSlots - :decrement"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":decrement": &types.AttributeValueMemberN{Value: "1"},
				},
			},
		},
		// 5 Mark slot as occupied
		{
			Update: &types.Update{
				TableName: aws.String(config.DynamoDBTable),
				Key: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PrefixBuilding, vehicle.AssignedBuildingID.String())},
					"SK": &types.AttributeValueMemberS{Value: slotKey},
				},
				UpdateExpression: aws.String("SET IsOccupied = :isOccupied, OccupiedBy = :occupiedBy"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":isOccupied": &types.AttributeValueMemberBOOL{Value: true},
					":occupiedBy": &types.AttributeValueMemberM{Value: map[string]types.AttributeValue{
						"Username":    &types.AttributeValueMemberS{Value: username},
						"Numberplate": &types.AttributeValueMemberS{Value: vehicle.NumberPlate},
						"Email":       &types.AttributeValueMemberS{Value: email},
						"StartTime":   &types.AttributeValueMemberN{Value: strconv.FormatInt(timestamp, 10)},
					}},
				},
			},
		},
	}

	_, err = nosqlpr.client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	})
	if err != nil {
		log.Println("Error in transaction:", err.Error())
		return "", errors.New(constants.ErrAddingParkingHistory)
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
			":pk":        &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PrefixUser, userCtx.Email)},
			":sk":        &types.AttributeValueMemberS{Value: constants.PKParking},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return errors.New(constants.ErrFindingParkingRecord)
	}

	if len(queryRes.Items) == 0 {
		log.Println("parking record not found or already unparked in Unpark")
		return errors.New(constants.ErrParkingRecordNotFound)
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
		// 1 Update parking record with EndTime
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
		// 2 Update vehicle IsParked status
		{
			Update: &types.Update{
				TableName: aws.String(config.DynamoDBTable),
				Key: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: PK},
					"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PKVehicle, numberplate)},
				},
				UpdateExpression: aws.String("SET IsParked = :isParked"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":isParked": &types.AttributeValueMemberBOOL{Value: false},
				},
			},
		},
		// 3 Increment floor AvailableSlots
		{
			Update: &types.Update{
				TableName: aws.String(config.DynamoDBTable),
				Key: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PrefixBuilding, buildingID)},
					"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%d", constants.PrefixFloorInfo, floorNumber)},
				},
				UpdateExpression: aws.String("SET AvailableSlots = AvailableSlots + :increment"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":increment": &types.AttributeValueMemberN{Value: "1"},
				},
			},
		},
		// 4 Increment building AvailableSlots
		{
			Update: &types.Update{
				TableName: aws.String(config.DynamoDBTable),
				Key: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: constants.PKBuilding},
					"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PrefixBuilding, buildingID)},
				},
				UpdateExpression: aws.String("SET AvailableSlots = AvailableSlots + :increment"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":increment": &types.AttributeValueMemberN{Value: "1"},
				},
			},
		},
		// 5 Mark slot as unoccupied
		{
			Update: &types.Update{
				TableName: aws.String(config.DynamoDBTable),
				Key: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PrefixBuilding, buildingID)},
					"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%d#%s%d", constants.PrefixFloor, floorNumber, constants.PrefixSlot, slotId)},
				},
				UpdateExpression: aws.String("SET IsOccupied = :isOccupied REMOVE OccupiedBy"),
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
		return errors.New(constants.ErrUnparkingVehicle)
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
			":pk":          &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PrefixUser, userCtx.Email)},
			":sk":          &types.AttributeValueMemberS{Value: constants.PKParking},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New(constants.ErrFetchingParkingHistory)
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
			":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PrefixUser, userId)},
			":sk": &types.AttributeValueMemberS{Value: constants.PKParking},
		},
	})
	if err != nil {
		log.Println("Error fetching parking history:", err.Error())
		return history, errors.New(constants.ErrFetchingParkingHistory)
	}

	log.Printf("Found %d parking records for user %s", len(queryRes.Items), userId)
	log.Printf("Filtering parking records for time range: %d to %d (%v to %v)", startTimestamp, endTimestamp, startTime, endTime)

	for _, item := range queryRes.Items {
		if item["EndTime"] == nil {
			log.Printf("Skipping parking record without EndTime (active parking)")
			continue
		}

		if item["StartTime"] != nil {
			if startNum, ok := item["StartTime"].(*types.AttributeValueMemberN); ok && startNum != nil {
				itemStartTime, _ := strconv.ParseInt(startNum.Value, 10, 64)
				if itemStartTime < startTimestamp || itemStartTime > endTimestamp {
					log.Printf("Skipping parking record with StartTime %d (outside range %d-%d)", itemStartTime, startTimestamp, endTimestamp)
					continue
				}
			} else {
				log.Printf("Skipping parking record with unexpected StartTime type (expected number)")
				continue
			}
		}

		var dto models.ParkingHistoryDTO
		dto.TicketId = item["ParkingId"].(*types.AttributeValueMemberS).Value
		dto.NumberPlate = item["Numberplate"].(*types.AttributeValueMemberS).Value
		dto.BuildingId = item["BuildingId"].(*types.AttributeValueMemberS).Value
		if floorNum, ok := item["FloorNumber"].(*types.AttributeValueMemberN); ok && floorNum != nil {
			dto.FLoorNumber, _ = strconv.Atoi(floorNum.Value)
		}
		if slotNum, ok := item["SlotId"].(*types.AttributeValueMemberN); ok && slotNum != nil {
			dto.SlotNumber, _ = strconv.Atoi(slotNum.Value)
		}
		dto.VechicleType = item["VehicleType"].(*types.AttributeValueMemberS).Value

		if startNum, ok := item["StartTime"].(*types.AttributeValueMemberN); ok && startNum != nil {
			startTimeUnix, _ := strconv.ParseInt(startNum.Value, 10, 64)
			dto.StartTime = time.Unix(startTimeUnix, 0).Local()
		} else {
			log.Printf("Skipping parking record %s due to invalid StartTime", dto.TicketId)
			continue
		}

		if endNum, ok := item["EndTime"].(*types.AttributeValueMemberN); ok && endNum != nil {
			endTimeUnix, _ := strconv.ParseInt(endNum.Value, 10, 64)
			dto.EndTime = time.Unix(endTimeUnix, 0).Local()
		} else {
			log.Printf("Skipping parking record %s due to invalid EndTime", dto.TicketId)
			continue
		}

		buildingRes, err := nosqlpr.client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(config.DynamoDBTable),
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: constants.PKBuilding},
				"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PrefixBuilding, dto.BuildingId)},
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
			":pk":          &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PrefixUser, userCtx.ID)},
			":sk":          &types.AttributeValueMemberS{Value: constants.PKParking},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return errors.New(constants.ErrFindingParkingRecord)
	}

	if len(queryRes.Items) == 0 {
		log.Println("no active parking found for numberplate:", numberplate)
		return errors.New(constants.ErrNoActiveParkingFound)
	}

	item := queryRes.Items[0]
	pk := item["PK"].(*types.AttributeValueMemberS).Value
	timestamp := item["SK"].(*types.AttributeValueMemberS).Value
	buildingID := item["BuildingId"].(*types.AttributeValueMemberS).Value
	floorNumber, _ := strconv.Atoi(item["FloorNumber"].(*types.AttributeValueMemberN).Value)
	slotId, _ := strconv.Atoi(item["SlotId"].(*types.AttributeValueMemberN).Value)

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
		return errors.New(constants.ErrUpdatingParkingRecord)
	}

	_, err = nosqlpr.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PKVehicle, numberplate)},
		},
		UpdateExpression: aws.String("SET IsParked = :isParked"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":isParked": &types.AttributeValueMemberBOOL{Value: false},
		},
	})
	if err != nil {
		log.Println("Warning: could not update IsParked status:", err.Error())
	}

	// update floor AvailableSlots (increment by 1)
	_, err = nosqlpr.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PrefixBuilding, buildingID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%d", constants.PrefixFloorInfo, floorNumber)},
		},
		UpdateExpression: aws.String("SET AvailableSlots = AvailableSlots + :increment"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":increment": &types.AttributeValueMemberN{Value: "1"},
		},
	})
	if err != nil {
		log.Println("Warning: could not update floor AvailableSlots:", err.Error())
	}

	// update building AvailableSlots (increment by 1)
	_, err = nosqlpr.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: constants.PKBuilding},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PrefixBuilding, buildingID)},
		},
		UpdateExpression: aws.String("SET AvailableSlots = AvailableSlots + :increment"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":increment": &types.AttributeValueMemberN{Value: "1"},
		},
	})
	if err != nil {
		log.Println("Warning: could not update building AvailableSlots:", err.Error())
	}

	// mark slot as unoccupied
	_, err = nosqlpr.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(config.DynamoDBTable),
				Key: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PrefixBuilding, buildingID)},
					"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%d#%s%d", constants.PrefixFloor, floorNumber, constants.PrefixSlot, slotId)},
				},
		UpdateExpression: aws.String("SET IsOccupied = :isOccupied REMOVE OccupiedBy"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":isOccupied": &types.AttributeValueMemberBOOL{Value: false},
		},
	})
	if err != nil {
		log.Println("Warning: could not update slot status:", err.Error())
	}

	return nil
}
