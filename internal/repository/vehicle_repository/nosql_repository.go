package vehiclerepository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/Kaushik1766/ParkingManagement/internal/config"
	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type NOSQLVehicleRepository struct {
	client *dynamodb.Client
}

func NewNOSQLVehicleRepository(client *dynamodb.Client) *NOSQLVehicleRepository {
	return &NOSQLVehicleRepository{
		client: client,
	}
}

func (nosqlvr *NOSQLVehicleRepository) AddVehicle(numberplate string, userid uuid.UUID, vehicleType vehicletypes.VehicleType) (models.Vehicle, error) {
	// First get user email from userid
	userRes, err := nosqlvr.client.Query(context.Background(), &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userid.String())},
			":sk": &types.AttributeValueMemberS{Value: "PROFILE#"},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return models.Vehicle{}, errors.New("error fetching user")
	}

	if len(userRes.Items) == 0 {
		return models.Vehicle{}, errors.New("user not found")
	}

	userEmail := userRes.Items[0]["Email"].(*types.AttributeValueMemberS).Value
	userName := userRes.Items[0]["Username"].(*types.AttributeValueMemberS).Value

	vehicle := models.Vehicle{
		VehicleID:   uuid.New(),
		NumberPlate: numberplate,
		UserID:      userid,
		VehicleType: vehicleType,
		IsActive:    true,
		UserEmail:   userEmail,
	}

	item := map[string]types.AttributeValue{
		"PK":          &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userEmail)},
		"SK":          &types.AttributeValueMemberS{Value: fmt.Sprintf("VEHICLE#%s", numberplate)},
		"VehicleId":   &types.AttributeValueMemberS{Value: vehicle.VehicleID.String()},
		"Numberplate": &types.AttributeValueMemberS{Value: numberplate},
		"VehicleType": &types.AttributeValueMemberS{Value: vehicleType.String()},
		"IsActive":    &types.AttributeValueMemberBOOL{Value: true},
		"UserId":      &types.AttributeValueMemberS{Value: userid.String()},
		"Username":    &types.AttributeValueMemberS{Value: userName},
	}

	_, err = nosqlvr.client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Item:      item,
	})
	if err != nil {
		log.Println(err.Error())
		return models.Vehicle{}, errors.New("error adding vehicle")
	}

	return vehicle, nil
}

func (nosqlvr *NOSQLVehicleRepository) RemoveVehicle(numberplate string) error {
	// Check if vehicle is parked
	isParked, err := nosqlvr.GetParkingStatus(numberplate)
	if err != nil {
		return err
	}

	if isParked {
		return errors.New("vehicle is parked please unpark it first")
	}

	// Find the vehicle
	scanRes, err := nosqlvr.client.Scan(context.Background(), &dynamodb.ScanInput{
		TableName:        aws.String(config.DynamoDBTable),
		FilterExpression: aws.String("Numberplate = :numberplate AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":numberplate": &types.AttributeValueMemberS{Value: numberplate},
			":sk":          &types.AttributeValueMemberS{Value: "VEHICLE#"},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return errors.New("error finding vehicle")
	}

	if len(scanRes.Items) == 0 {
		return errors.New("vehicle not found")
	}

	item := scanRes.Items[0]
	userEmail := item["PK"].(*types.AttributeValueMemberS).Value
	vehicleSK := item["SK"].(*types.AttributeValueMemberS).Value

	// Mark as inactive instead of deleting
	_, err = nosqlvr.client.UpdateItem(context.Background(), &dynamodb.UpdateItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: userEmail},
			"SK": &types.AttributeValueMemberS{Value: vehicleSK},
		},
		UpdateExpression: aws.String("SET IsActive = :val"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":val": &types.AttributeValueMemberBOOL{Value: false},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return errors.New("error removing vehicle")
	}

	return nil
}

func (nosqlvr *NOSQLVehicleRepository) GetVehicleById(vehicleId uuid.UUID) (models.Vehicle, error) {
	var vehicle models.Vehicle

	scanRes, err := nosqlvr.client.Scan(context.Background(), &dynamodb.ScanInput{
		TableName:        aws.String(config.DynamoDBTable),
		FilterExpression: aws.String("VehicleId = :vehicleId AND IsActive = :active"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":vehicleId": &types.AttributeValueMemberS{Value: vehicleId.String()},
			":active":    &types.AttributeValueMemberBOOL{Value: true},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return vehicle, errors.New("error fetching vehicle")
	}

	if len(scanRes.Items) == 0 {
		return vehicle, errors.New("vehicle not found")
	}

	item := scanRes.Items[0]
	vehicle = nosqlvr.itemToVehicle(item)

	return vehicle, nil
}

func (nosqlvr *NOSQLVehicleRepository) GetVehiclesByUserId(userId uuid.UUID) ([]models.Vehicle, error) {
	var vehicles []models.Vehicle

	// First get user email
	userRes, err := nosqlvr.client.Query(context.Background(), &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userId.String())},
			":sk": &types.AttributeValueMemberS{Value: "PROFILE#"},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("error fetching user")
	}

	if len(userRes.Items) == 0 {
		return nil, errors.New("user not found")
	}

	userEmail := userRes.Items[0]["Email"].(*types.AttributeValueMemberS).Value

	// Query vehicles for this user
	vehiclesRes, err := nosqlvr.client.Query(context.Background(), &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		FilterExpression:       aws.String("IsActive = :active"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":     &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userEmail)},
			":sk":     &types.AttributeValueMemberS{Value: "VEHICLE#"},
			":active": &types.AttributeValueMemberBOOL{Value: true},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("error fetching vehicles")
	}

	for _, item := range vehiclesRes.Items {
		vehicle := nosqlvr.itemToVehicle(item)
		vehicles = append(vehicles, vehicle)
	}

	return vehicles, nil
}

func (nosqlvr *NOSQLVehicleRepository) GetVehicleByNumberPlate(numberplate string) (models.Vehicle, error) {
	var vehicle models.Vehicle

	scanRes, err := nosqlvr.client.Scan(context.Background(), &dynamodb.ScanInput{
		TableName:        aws.String(config.DynamoDBTable),
		FilterExpression: aws.String("Numberplate = :numberplate AND IsActive = :active AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":numberplate": &types.AttributeValueMemberS{Value: numberplate},
			":active":      &types.AttributeValueMemberBOOL{Value: true},
			":sk":          &types.AttributeValueMemberS{Value: "VEHICLE#"},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return vehicle, errors.New("error fetching vehicle")
	}

	if len(scanRes.Items) == 0 {
		return vehicle, errors.New("vehicle not found")
	}

	item := scanRes.Items[0]
	vehicle = nosqlvr.itemToVehicle(item)

	return vehicle, nil
}

func (nosqlvr *NOSQLVehicleRepository) GetVehiclesWithUnassignedSlots() (vehicles []models.Vehicle, err error) {
	scanRes, err := nosqlvr.client.Scan(context.Background(), &dynamodb.ScanInput{
		TableName:        aws.String(config.DynamoDBTable),
		FilterExpression: aws.String("attribute_not_exists(AssignedSlot) AND begins_with(SK, :sk) AND IsActive = :active"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":sk":     &types.AttributeValueMemberS{Value: "VEHICLE#"},
			":active": &types.AttributeValueMemberBOOL{Value: true},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("error fetching vehicles with unassigned slots")
	}

	for _, item := range scanRes.Items {
		vehicle := nosqlvr.itemToVehicle(item)
		vehicles = append(vehicles, vehicle)
	}

	return vehicles, nil
}

func (nosqlvr *NOSQLVehicleRepository) GetParkingStatus(numberplate string) (bool, error) {
	// Check if there's an active parking (without EndTime) for this numberplate
	scanRes, err := nosqlvr.client.Scan(context.Background(), &dynamodb.ScanInput{
		TableName:        aws.String(config.DynamoDBTable),
		FilterExpression: aws.String("Numberplate = :numberplate AND attribute_not_exists(EndTime) AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":numberplate": &types.AttributeValueMemberS{Value: numberplate},
			":sk":          &types.AttributeValueMemberS{Value: "PARKING#"},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return false, errors.New("error checking parking status")
	}

	return len(scanRes.Items) > 0, nil
}

func (nosqlvr *NOSQLVehicleRepository) Save(vehicle models.Vehicle) error {
	// Get user email from vehicle
	userEmail := vehicle.UserEmail
	if userEmail == "" {
		// Fetch user email if not provided
		userRes, err := nosqlvr.client.Query(context.Background(), &dynamodb.QueryInput{
			TableName:              aws.String(config.DynamoDBTable),
			KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", vehicle.UserID.String())},
				":sk": &types.AttributeValueMemberS{Value: "PROFILE#"},
			},
		})
		if err != nil {
			log.Println(err.Error())
			return errors.New("error fetching user")
		}

		if len(userRes.Items) == 0 {
			return errors.New("user not found")
		}

		userEmail = userRes.Items[0]["Email"].(*types.AttributeValueMemberS).Value
	}

	updateExpression := "SET VehicleType = :vehicleType, IsActive = :isActive"
	expressionValues := map[string]types.AttributeValue{
		":vehicleType": &types.AttributeValueMemberS{Value: vehicle.VehicleType.String()},
		":isActive":    &types.AttributeValueMemberBOOL{Value: vehicle.IsActive},
	}

	// Add assigned slot info if present
	if vehicle.AssignedBuildingID != uuid.Nil {
		updateExpression += ", AssignedSlot = :assignedSlot"
		assignedSlot := map[string]types.AttributeValue{
			"BuildingId":  &types.AttributeValueMemberS{Value: vehicle.AssignedBuildingID.String()},
			"FloorNumber": &types.AttributeValueMemberN{Value: strconv.Itoa(vehicle.AssignedFloorNumber)},
			"SlotId":      &types.AttributeValueMemberN{Value: strconv.Itoa(vehicle.AssignedSlotNumber)},
		}
		expressionValues[":assignedSlot"] = &types.AttributeValueMemberM{Value: assignedSlot}
	}

	_, err := nosqlvr.client.UpdateItem(context.Background(), &dynamodb.UpdateItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userEmail)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("VEHICLE#%s", vehicle.NumberPlate)},
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionValues,
	})
	if err != nil {
		log.Println(err.Error())
		return errors.New("error saving vehicle")
	}

	return nil
}

// Helper function to convert DynamoDB item to Vehicle model
func (nosqlvr *NOSQLVehicleRepository) itemToVehicle(item map[string]types.AttributeValue) models.Vehicle {
	var vehicle models.Vehicle

	vehicle.VehicleID = uuid.MustParse(item["VehicleId"].(*types.AttributeValueMemberS).Value)
	vehicle.NumberPlate = item["Numberplate"].(*types.AttributeValueMemberS).Value
	vehicle.UserID = uuid.MustParse(item["UserId"].(*types.AttributeValueMemberS).Value)

	vehicleTypeStr := item["VehicleType"].(*types.AttributeValueMemberS).Value
	if vehicleTypeStr == "TwoWheeler" {
		vehicle.VehicleType = vehicletypes.TwoWheeler
	} else {
		vehicle.VehicleType = vehicletypes.FourWheeler
	}

	vehicle.IsActive = item["IsActive"].(*types.AttributeValueMemberBOOL).Value

	// Extract user email from PK
	pkValue := item["PK"].(*types.AttributeValueMemberS).Value
	vehicle.UserEmail = pkValue[5:] // Remove "USER#" prefix

	// Get username if available
	if item["Username"] != nil {
		vehicle.User.Name = item["Username"].(*types.AttributeValueMemberS).Value
	}

	// Handle AssignedSlot if present
	if item["AssignedSlot"] != nil {
		assignedSlotMap := item["AssignedSlot"].(*types.AttributeValueMemberM).Value
		if buildingId, ok := assignedSlotMap["BuildingId"]; ok {
			vehicle.AssignedBuildingID = uuid.MustParse(buildingId.(*types.AttributeValueMemberS).Value)
		}
		if floorNumber, ok := assignedSlotMap["FloorNumber"]; ok {
			vehicle.AssignedFloorNumber, _ = strconv.Atoi(floorNumber.(*types.AttributeValueMemberN).Value)
		}
		if slotId, ok := assignedSlotMap["SlotId"]; ok {
			vehicle.AssignedSlotNumber, _ = strconv.Atoi(slotId.(*types.AttributeValueMemberN).Value)
		}

		// Populate AssignedSlot model
		vehicle.AssignedSlot.BuildingID = vehicle.AssignedBuildingID
		vehicle.AssignedSlot.FloorNumber = vehicle.AssignedFloorNumber
		vehicle.AssignedSlot.SlotNumber = vehicle.AssignedSlotNumber
	}

	return vehicle
}
