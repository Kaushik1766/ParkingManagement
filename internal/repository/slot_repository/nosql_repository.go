package slotrepository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Kaushik1766/ParkingManagement/internal/config"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type NOSQLSlotRepository struct {
	client *dynamodb.Client
}

func NewNOSQLSlotRepository(client *dynamodb.Client) *NOSQLSlotRepository {
	return &NOSQLSlotRepository{
		client: client,
	}
}

func (nosqlsr *NOSQLSlotRepository) AddSlot(ctx context.Context, buildingId uuid.UUID, floorNumber, slotNumber int, slotType vehicletypes.VehicleType) error {
	item := map[string]types.AttributeValue{
		"PK":         &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", buildingId.String())},
		"SK":         &types.AttributeValueMemberS{Value: fmt.Sprintf("FLOOR#%d#SLOT#%d", floorNumber, slotNumber)},
		"SlotNumber": &types.AttributeValueMemberN{Value: strconv.Itoa(slotNumber)},
		"SlotType":   &types.AttributeValueMemberS{Value: slotType.String()},
	}

	_, err := nosqlsr.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Item:      item,
	})
	if err != nil {
		log.Println(err.Error())
		return errors.New("error adding slot")
	}

	return nil
}

func (nosqlsr *NOSQLSlotRepository) DeleteSlot(ctx context.Context, buildingId uuid.UUID, floorNumber, slotNumber int) error {
	_, err := nosqlsr.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", buildingId.String())},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("FLOOR#%d#SLOT#%d", floorNumber, slotNumber)},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return errors.New("error deleting slot")
	}

	return nil
}

func (nosqlsr *NOSQLSlotRepository) GetSlotsByFloor(ctx context.Context, buildingId uuid.UUID, floorNumber int) ([]models.Slot, error) {
	var slots []models.Slot

	queryRes, err := nosqlsr.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", buildingId.String())},
			":sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("FLOOR#%d#SLOT#", floorNumber)},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("error fetching slots")
	}

	// Get all active parkings to determine which slots have parked vehicles
	activeParkings := make(map[string]models.Vehicle)

	scanRes, err := nosqlsr.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(config.DynamoDBTable),
		FilterExpression: aws.String("attribute_not_exists(EndTime) AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":sk": &types.AttributeValueMemberS{Value: "PARKING#"},
		},
	})
	if err == nil {
		for _, item := range scanRes.Items {
			if item["BuildingId"] != nil && item["FloorNumber"] != nil && item["SlotId"] != nil {
				bId := item["BuildingId"].(*types.AttributeValueMemberS).Value
				fNum, _ := strconv.Atoi(item["FloorNumber"].(*types.AttributeValueMemberN).Value)
				sNum, _ := strconv.Atoi(item["SlotId"].(*types.AttributeValueMemberN).Value)

				if bId == buildingId.String() && fNum == floorNumber {
					slotKey := fmt.Sprintf("%s_%d_%d", bId, fNum, sNum)

					var vehicle models.Vehicle
					vehicle.NumberPlate = item["Numberplate"].(*types.AttributeValueMemberS).Value

					// Get user email from PK
					userPK := item["PK"].(*types.AttributeValueMemberS).Value
					vehicle.UserEmail = userPK[5:] // Remove "USER#" prefix

					// Fetch user details
					userRes, err := nosqlsr.client.Query(ctx, &dynamodb.QueryInput{
						TableName:              aws.String(config.DynamoDBTable),
						KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
						ExpressionAttributeValues: map[string]types.AttributeValue{
							":pk": &types.AttributeValueMemberS{Value: userPK},
							":sk": &types.AttributeValueMemberS{Value: "PROFILE#"},
						},
					})
					if err == nil && len(userRes.Items) > 0 {
						vehicle.User.Name = userRes.Items[0]["Username"].(*types.AttributeValueMemberS).Value
						vehicle.User.Email = userRes.Items[0]["Email"].(*types.AttributeValueMemberS).Value
					}

					// Parse parking history
					if item["StartTime"] != nil {
						startTimeUnix, _ := strconv.ParseInt(item["StartTime"].(*types.AttributeValueMemberN).Value, 10, 64)
						startTime := time.Unix(startTimeUnix, 0)
						vehicle.ParkingHistory = []models.ParkingHistory{
							{
								ParkingID: uuid.MustParse(item["ParkingId"].(*types.AttributeValueMemberS).Value),
								StartTime: startTime,
							},
						}
					}

					activeParkings[slotKey] = vehicle
				}
			}
		}
	}

	for _, item := range queryRes.Items {
		var slot models.Slot
		slot.BuildingID = buildingId
		slot.FloorNumber = floorNumber
		slot.SlotNumber, _ = strconv.Atoi(item["SlotNumber"].(*types.AttributeValueMemberN).Value)

		slotTypeStr := item["SlotType"].(*types.AttributeValueMemberS).Value
		if slotTypeStr == "TwoWheeler" {
			slot.SlotType = vehicletypes.TwoWheeler
		} else {
			slot.SlotType = vehicletypes.FourWheeler
		}

		// Check if slot has parked vehicle
		slotKey := fmt.Sprintf("%s_%d_%d", buildingId.String(), floorNumber, slot.SlotNumber)
		if vehicle, ok := activeParkings[slotKey]; ok {
			slot.Vehicles = []models.Vehicle{vehicle}
		}

		slots = append(slots, slot)
	}

	return slots, nil
}

func (nosqlsr *NOSQLSlotRepository) GetFreeSlotsByFloor(ctx context.Context, buildingId uuid.UUID, floorNumber int) ([]models.Slot, error) {
	allSlots, err := nosqlsr.GetSlotsByFloor(ctx, buildingId, floorNumber)
	if err != nil {
		return nil, err
	}

	var freeSlots []models.Slot
	for _, slot := range allSlots {
		if len(slot.Vehicles) == 0 {
			freeSlots = append(freeSlots, slot)
		}
	}

	return freeSlots, nil
}

func (nosqlsr *NOSQLSlotRepository) GetFreeSlotsByBuilding(ctx context.Context, buildingId uuid.UUID) ([]models.Slot, error) {
	var freeSlots []models.Slot

	// Get all floors for this building
	floorsRes, err := nosqlsr.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", buildingId.String())},
			":sk": &types.AttributeValueMemberS{Value: "FLOORINFO#"},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("error fetching floors")
	}

	// For each floor, get free slots
	for _, floorItem := range floorsRes.Items {
		floorNum, _ := strconv.Atoi(floorItem["FloorNumber"].(*types.AttributeValueMemberN).Value)
		floorFreeSlots, err := nosqlsr.GetFreeSlotsByFloor(ctx, buildingId, floorNum)
		if err == nil {
			freeSlots = append(freeSlots, floorFreeSlots...)
		}
	}

	return freeSlots, nil
}

func (nosqlsr *NOSQLSlotRepository) Save(ctx context.Context, slot models.Slot) error {
	updateExpression := "SET SlotType = :slotType"
	expressionValues := map[string]types.AttributeValue{
		":slotType": &types.AttributeValueMemberS{Value: slot.SlotType.String()},
	}

	// Add OccupiedBy if slot has vehicles
	if len(slot.Vehicles) > 0 {
		vehicle := slot.Vehicles[0]
		occupiedBy := map[string]types.AttributeValue{
			"Username":    &types.AttributeValueMemberS{Value: vehicle.User.Name},
			"Numberplate": &types.AttributeValueMemberS{Value: vehicle.NumberPlate},
			"Email":       &types.AttributeValueMemberS{Value: vehicle.User.Email},
		}
		updateExpression += ", OccupiedBy = :occupiedBy, IsOccupied = :isOccupied"
		expressionValues[":occupiedBy"] = &types.AttributeValueMemberM{Value: occupiedBy}
		expressionValues[":isOccupied"] = &types.AttributeValueMemberBOOL{Value: true}
	} else {
		updateExpression += ", IsOccupied = :isOccupied"
		expressionValues[":isOccupied"] = &types.AttributeValueMemberBOOL{Value: false}
		updateExpression += " REMOVE OccupiedBy"
	}

	_, err := nosqlsr.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", slot.BuildingID.String())},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("FLOOR#%d#SLOT#%d", slot.FloorNumber, slot.SlotNumber)},
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionValues,
	})
	if err != nil {
		log.Println(err.Error())
		return errors.New("error saving slot")
	}

	return nil
}
