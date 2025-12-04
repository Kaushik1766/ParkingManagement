package vehiclerepository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/Kaushik1766/ParkingManagement/internal/config"
	"github.com/Kaushik1766/ParkingManagement/internal/constants"
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

func (nosqlvr *NOSQLVehicleRepository) AddVehicle(ctx context.Context, numberplate string, userid uuid.UUID, vehicleType vehicletypes.VehicleType) (models.Vehicle, error) {

	// check if vehicle with same numberplate already exists for the user
	existingVehicleRes, err := nosqlvr.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userid.String())},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("VEHICLE#%s", numberplate)},
		},
	})
	if err != nil {
		log.Println("Error checking existing vehicle:", err.Error())
		return models.Vehicle{}, errors.New("error adding vehicle")
	}

	if existingVehicleRes.Item != nil {
		log.Println("Vehicle with same numberplate already exists for the user")
		return models.Vehicle{}, errors.New("vehicle with same numberplate already exists")
	}

	userCtx := ctx.Value(constants.User).(models.UserJwt)
	userEmail := userCtx.Email

	vehicle := models.Vehicle{
		VehicleID:   uuid.New(),
		NumberPlate: numberplate,
		UserID:      userid,
		VehicleType: vehicleType,
		IsActive:    true,
		UserEmail:   userEmail,
	}

	// check if user has other vehicles of the same type to reuse their slot assignment
	vehiclesRes, err := nosqlvr.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		FilterExpression:       aws.String("VehicleType = :vehicleType AND attribute_exists(AssignedSlot)"),
		Limit:                  aws.Int32(1),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":          &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userid.String())},
			":sk":          &types.AttributeValueMemberS{Value: "VEHICLE#"},
			":vehicleType": &types.AttributeValueMemberS{Value: vehicleType.String()},
		},
	})
	if err != nil {
		log.Println("Error fetching user vehicles:", err.Error())
		return models.Vehicle{}, errors.New("error adding vehicle")
	}

	// vehicle of same type found
	if vehiclesRes.Count > 0 {
		item := vehiclesRes.Items[0]
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
	}

	item := map[string]types.AttributeValue{
		"PK":          &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userid.String())},
		"SK":          &types.AttributeValueMemberS{Value: fmt.Sprintf("VEHICLE#%s", numberplate)},
		"VehicleId":   &types.AttributeValueMemberS{Value: vehicle.VehicleID.String()},
		"Numberplate": &types.AttributeValueMemberS{Value: numberplate},
		"VehicleType": &types.AttributeValueMemberS{Value: vehicleType.String()},
		"IsParked":    &types.AttributeValueMemberBOOL{Value: false},
	}

	// Add AssignedSlot if it was set from reusing another vehicle's slot
	if vehicle.AssignedBuildingID != uuid.Nil {
		assignedSlot := map[string]types.AttributeValue{
			"BuildingId":  &types.AttributeValueMemberS{Value: vehicle.AssignedBuildingID.String()},
			"FloorNumber": &types.AttributeValueMemberN{Value: strconv.Itoa(vehicle.AssignedFloorNumber)},
			"SlotId":      &types.AttributeValueMemberN{Value: strconv.Itoa(vehicle.AssignedSlotNumber)},
		}
		item["AssignedSlot"] = &types.AttributeValueMemberM{Value: assignedSlot}
	}

	_, err = nosqlvr.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Item:      item,
	})
	if err != nil {
		log.Println(err.Error())
		return models.Vehicle{}, errors.New("error adding vehicle")
	}

	return vehicle, nil
}

func (nosqlvr *NOSQLVehicleRepository) RemoveVehicle(ctx context.Context, numberplate string) error {
	// checkong if vehicle is parked
	isParked, err := nosqlvr.GetParkingStatus(ctx, numberplate)
	if err != nil {
		return err
	}

	if isParked {
		log.Println("vehicle is parked")
		return errors.New("vehicle is parked please unpark it first")
	}

	// get user id from ctx
	userCtx := ctx.Value(constants.User).(models.UserJwt)
	userId := userCtx.ID

	// Delete the vehicle item directly using PK and SK
	_, err = nosqlvr.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userId)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("VEHICLE#%s", numberplate)},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return errors.New("error removing vehicle")
	}

	return nil
}

func (nosqlvr *NOSQLVehicleRepository) GetVehicleById(ctx context.Context, vehicleId uuid.UUID) (models.Vehicle, error) {
	userCtx := ctx.Value(constants.User).(models.UserJwt)
	var vehicle models.Vehicle

	// TODO: change this to get vehicle by numberplate wherever used
	scanRes, err := nosqlvr.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(config.DynamoDBTable),
		FilterExpression: aws.String("VehicleId = :vehicleId AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":vehicleId": &types.AttributeValueMemberS{Value: vehicleId.String()},
			":sk":        &types.AttributeValueMemberS{Value: "VEHICLE#"},
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
	vehicle.UserID = uuid.MustParse(userCtx.ID)
	vehicle.UserEmail = userCtx.Email

	return vehicle, nil
}

func (nosqlvr *NOSQLVehicleRepository) GetVehiclesByUserId(ctx context.Context, userId uuid.UUID) ([]models.Vehicle, error) {
	var vehicles []models.Vehicle

	userCtx := ctx.Value(constants.User).(models.UserJwt)

	vehiclesRes, err := nosqlvr.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userId)},
			":sk": &types.AttributeValueMemberS{Value: "VEHICLE#"},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("error fetching vehicles")
	}

	for _, item := range vehiclesRes.Items {
		vehicle := nosqlvr.itemToVehicle(item)
		vehicle.UserID = userId
		vehicle.UserEmail = userCtx.Email
		vehicles = append(vehicles, vehicle)
	}

	return vehicles, nil
}

func (nosqlvr *NOSQLVehicleRepository) GetVehicleByNumberPlate(ctx context.Context, numberplate string) (models.Vehicle, error) {
	var vehicle models.Vehicle

	userCtx := ctx.Value(constants.User).(models.UserJwt)
	userId := userCtx.ID

	getRes, err := nosqlvr.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userId)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("VEHICLE#%s", numberplate)},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return vehicle, errors.New("error fetching vehicle")
	}

	if len(getRes.Item) == 0 {
		return vehicle, errors.New("vehicle not found")
	}

	vehicle = nosqlvr.itemToVehicle(getRes.Item)
	vehicle.UserID = uuid.MustParse(userId)
	vehicle.UserEmail = userCtx.Email

	return vehicle, nil
}

func (nosqlvr *NOSQLVehicleRepository) GetVehiclesWithUnassignedSlots(ctx context.Context) (vehicles []models.Vehicle, err error) {
	// scanRes, err := nosqlvr.client.Scan(ctx, &dynamodb.ScanInput{
	// 	TableName:        aws.String(config.DynamoDBTable),
	// 	FilterExpression: aws.String("attribute_not_exists(AssignedSlot) AND begins_with(SK, :sk)"),
	// 	ExpressionAttributeValues: map[string]types.AttributeValue{
	// 		":sk": &types.AttributeValueMemberS{Value: "VEHICLE#"},
	// 	},
	// })
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return nil, errors.New("error fetching vehicles with unassigned slots")
	// }

	// for _, item := range scanRes.Items {
	// 	vehicle := nosqlvr.itemToVehicle(item)
	// 	// Fetch UserID from email
	// 	userID, err := nosqlvr.getUserIDFromEmail(ctx, vehicle.UserEmail)
	// 	if err != nil {
	// 		log.Println("Warning: could not fetch UserID:", err.Error())
	// 	} else {
	// 		vehicle.UserID = userID
	// 	}
	// 	vehicles = append(vehicles, vehicle)
	// }

	// return vehicles, nil
	panic("Not implemented yet coz not used")
}

func (nosqlvr *NOSQLVehicleRepository) GetParkingStatus(ctx context.Context, numberplate string) (bool, error) {
	// Get user email from context
	userCtx := ctx.Value(constants.User).(models.UserJwt)
	userId := userCtx.ID

	// Query parking history for this user and filter by numberplate without EndTime
	queryRes, err := nosqlvr.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		FilterExpression:       aws.String("Numberplate = :numberplate AND attribute_not_exists(EndTime)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":          &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userId)},
			":sk":          &types.AttributeValueMemberS{Value: "PARKING#"},
			":numberplate": &types.AttributeValueMemberS{Value: numberplate},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return false, errors.New("error checking parking status")
	}

	return len(queryRes.Items) > 0, nil
}

func (nosqlvr *NOSQLVehicleRepository) Save(ctx context.Context, vehicle models.Vehicle) error {
	// Get user email from context
	userCtx := ctx.Value(constants.User).(models.UserJwt)
	userId := userCtx.ID

	updateExpression := "SET VehicleType = :vehicleType"
	expressionValues := map[string]types.AttributeValue{
		":vehicleType": &types.AttributeValueMemberS{Value: vehicle.VehicleType.String()},
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

	_, err := nosqlvr.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userId)},
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

// Helper function to get UserID from email
func (nosqlvr *NOSQLVehicleRepository) getUserIDFromEmail(ctx context.Context, email string) (uuid.UUID, error) {
	userRes, err := nosqlvr.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", email)},
			":sk": &types.AttributeValueMemberS{Value: "PROFILE#"},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return uuid.Nil, errors.New("error fetching user")
	}

	if len(userRes.Items) == 0 {
		log.Println("user not found in getUserIDFromEmail")
		return uuid.Nil, errors.New("user not found")
	}

	userID := uuid.MustParse(userRes.Items[0]["Id"].(*types.AttributeValueMemberS).Value)
	return userID, nil
}

// Helper function to convert DynamoDB item to Vehicle model
func (nosqlvr *NOSQLVehicleRepository) itemToVehicle(item map[string]types.AttributeValue) models.Vehicle {
	var vehicle models.Vehicle

	vehicle.VehicleID = uuid.MustParse(item["VehicleId"].(*types.AttributeValueMemberS).Value)
	vehicle.NumberPlate = item["Numberplate"].(*types.AttributeValueMemberS).Value

	vehicleTypeStr := item["VehicleType"].(*types.AttributeValueMemberS).Value
	if vehicleTypeStr == vehicletypes.TwoWheeler.String() {
		vehicle.VehicleType = vehicletypes.TwoWheeler
	} else {
		vehicle.VehicleType = vehicletypes.FourWheeler
	}

	vehicle.IsActive = true

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
