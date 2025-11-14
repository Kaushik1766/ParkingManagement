package floorrepository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/Kaushik1766/ParkingManagement/internal/config"
	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type NOSQLFloorRepository struct {
	client *dynamodb.Client
}

func (nosqlfr *NOSQLFloorRepository) AddFloor(buildingId string, floorNumber int) error {
	item := map[string]types.AttributeValue{}

	item["PK"] = &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", buildingId)}
	item["SK"] = &types.AttributeValueMemberS{Value: fmt.Sprintf("FLOOR#%s", strconv.Itoa(floorNumber))}
	item["FloorNumber"] = &types.AttributeValueMemberN{Value: strconv.Itoa(floorNumber)}
	item["TotalSlots"] = &types.AttributeValueMemberN{Value: strconv.Itoa(len(constants.SlotLayout))}
	item["AvailableSlots"] = &types.AttributeValueMemberN{Value: strconv.Itoa(len(constants.SlotLayout))}

	_, err := nosqlfr.client.
		PutItem(context.Background(), &dynamodb.PutItemInput{
			TableName: aws.String(config.DynamoDBTable),
			Item:      item,
		})
	if err != nil {
		log.Println(err.Error())
		return errors.New("error adding floor")
	}

	// adding slots to the floor

	slots := []types.WriteRequest{}

	for i, s := range constants.SlotLayout {
		if s == '0' {
			slots = append(slots, types.WriteRequest{
				PutRequest: &types.PutRequest{
					Item: map[string]types.AttributeValue{
						"PK":         &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", buildingId)},
						"SK":         &types.AttributeValueMemberS{Value: fmt.Sprintf("FLOOR#%s#SLOT#%s", floorNumber, i)},
						"SlotNumber": &types.AttributeValueMemberN{Value: strconv.Itoa(i)},
						"SlotType":   &types.AttributeValueMemberS{Value: vehicletypes.TwoWheeler.String()},
					},
				},
			})
		} else {
			slots = append(slots, types.WriteRequest{
				PutRequest: &types.PutRequest{
					Item: map[string]types.AttributeValue{
						"PK":         &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", buildingId)},
						"SK":         &types.AttributeValueMemberS{Value: fmt.Sprintf("FLOOR#%s#SLOT#%s", floorNumber, i)},
						"SlotNumber": &types.AttributeValueMemberN{Value: strconv.Itoa(i)},
						"SlotType":   &types.AttributeValueMemberS{Value: vehicletypes.FourWheeler.String()},
					},
				},
			})
		}

		if len(slots) == 25 {
			_, err := nosqlfr.client.
				BatchWriteItem(context.Background(), &dynamodb.BatchWriteItemInput{
					RequestItems: map[string][]types.WriteRequest{
						"pm_nosql": slots,
					},
				})
			if err != nil {
				log.Println(err.Error())
				return errors.New("error adding slots in floor" + strconv.Itoa(floorNumber))
			}
			slots = []types.WriteRequest{}
		}
	}

	return nil
}

func (nosqlfr *NOSQLFloorRepository) DeleteFloor(buildingId string, floorNumber int) error {
	// use of this function will be avoided coz need to delete the slots in the floors for consistency
	_, err := nosqlfr.client.
		DeleteItem(context.Background(), &dynamodb.DeleteItemInput{
			TableName: aws.String(config.DynamoDBTable),
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", buildingId)},
				"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("FLOOR#%s", floorNumber)},
			},
		})
	if err != nil {
		log.Println(err.Error())
		return errors.New("error deleting floor")
	}

	return nil
}

func (nosqlfr *NOSQLFloorRepository) GetFloor(buildingId uuid.UUID, floorNumber int) (int, error) {
	// var floor models.Floor
	// err := nosqlfr.db.Where("building_id = ? and floor_number = ?", buildingId, floorNumber).First(&floor).Error
	// if err != nil {
	// 	return 0, err
	// }
	// return floor.FloorNumber, nil
	panic("not implemented coz not used")
}

func (nosqlfr *NOSQLFloorRepository) GetFloorsByBuildingId(buildingId string) ([]models.FloorSummary, error) {
	// buildingUUID, err := uuid.Parse(buildingId)
	// if err != nil {
	// 	return nil, err
	// }
	// var floors []models.FloorSummary
	// //err = nosqlfr.db.Where("building_id = ?", buildingUUID).Find(&floors).Error
	// err = nosqlfr.db.
	// 	Table("floors").
	// 	Select(`
	// 		floors.building_id as building_id,
	// 		floors.floor_number as floor_number,
	// 		offices.office_name as assigned_office,
	// 		count(distinct case when slots.slot_number is not null then
	// 			(slots.building_id, slots.floor_number, slots.slot_number) end) as total_slots,
	// 		count(distinct case when vehicles.assigned_slot_number is null then
	// 			(slots.building_id,slots.floor_number, slots.slot_number) end) as available_slots
	// 	`).
	// 	Joins("left join slots on floors.building_id = slots.building_id and floors.floor_number = slots.floor_number").
	// 	Joins("left join offices on offices.building_id = floors.building_id and offices.floor_number = floors.floor_number").
	// 	Joins("left join vehicles on slots.building_id = vehicles.assigned_building_id and slots.floor_number = vehicles.assigned_floor_number and slots.slot_number = vehicles.assigned_slot_number").
	// 	Where("floors.building_id = ?", buildingUUID).
	// 	Group("floors.building_id, floors.floor_number, offices.office_id").
	// 	Find(&floors).Error
	// return floors, err

	floors = make([]models.FloorSummary, 0)

	res, err := nosqlfr.client.
		Query(context.Background(), &dynamodb.QueryInput{
			TableName:              aws.String(config.DynamoDBTable),
			KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", buildingId)},
				":sk": &types.AttributeValueMemberS{Value: "FLOOR"},
			},
		})
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("error fetching floors")
	}
}

func NewNOSQLFloorRepository(client *dynamodb.Client) *NOSQLFloorRepository {
	return &NOSQLFloorRepository{
		client: client,
	}
}
