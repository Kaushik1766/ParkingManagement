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

func (nosqlfr *NOSQLFloorRepository) AddFloor(ctx context.Context, buildingId string, floorNumber int) error {
	item := map[string]types.AttributeValue{}

	item["PK"] = &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", buildingId)}
	item["SK"] = &types.AttributeValueMemberS{Value: fmt.Sprintf("FLOORINFO#%s", strconv.Itoa(floorNumber))}
	item["FloorNumber"] = &types.AttributeValueMemberN{Value: strconv.Itoa(floorNumber)}
	item["TotalSlots"] = &types.AttributeValueMemberN{Value: strconv.Itoa(len(constants.SlotLayout))}
	item["AvailableSlots"] = &types.AttributeValueMemberN{Value: strconv.Itoa(len(constants.SlotLayout))}

	_, err := nosqlfr.client.
		PutItem(ctx, &dynamodb.PutItemInput{
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
						"SK":         &types.AttributeValueMemberS{Value: fmt.Sprintf("FLOOR#%d#SLOT#%d", floorNumber, i)},
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
						"SK":         &types.AttributeValueMemberS{Value: fmt.Sprintf("FLOOR#%d#SLOT#%d", floorNumber, i)},
						"SlotNumber": &types.AttributeValueMemberN{Value: strconv.Itoa(i)},
						"SlotType":   &types.AttributeValueMemberS{Value: vehicletypes.FourWheeler.String()},
					},
				},
			})
		}

		if len(slots) == 25 {
			_, err := nosqlfr.client.
				BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
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

	// update building details
	_, err = nosqlfr.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "BUILDING"},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", buildingId)},
		},
		UpdateExpression: aws.String("SET TotalFloors = TotalFloors + :inc, TotalSlots = TotalSlots + :slots, AvailableSlots = AvailableSlots + :avail"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":inc":   &types.AttributeValueMemberN{Value: "1"},
			":slots": &types.AttributeValueMemberN{Value: strconv.Itoa(len(constants.SlotLayout))},
			":avail": &types.AttributeValueMemberN{Value: strconv.Itoa(len(constants.SlotLayout))},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return errors.New("error updating building details")
	}

	return nil
}

func (nosqlfr *NOSQLFloorRepository) DeleteFloor(ctx context.Context, buildingId string, floorNumber int) error {
	// use of this function will be avoided coz need to delete the slots in the floors for consistency
	_, err := nosqlfr.client.
		DeleteItem(ctx, &dynamodb.DeleteItemInput{
			TableName: aws.String(config.DynamoDBTable),
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", buildingId)},
				"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("FLOORINFO#%d", floorNumber)},
			},
		})
	if err != nil {
		log.Println(err.Error())
		return errors.New("error deleting floor")
	}

	return nil
}

func (nosqlfr *NOSQLFloorRepository) GetFloor(ctx context.Context, buildingId uuid.UUID, floorNumber int) (int, error) {
	panic("not implemented coz not used")
}

func (nosqlfr *NOSQLFloorRepository) GetFloorsByBuildingId(ctx context.Context, buildingId string) ([]models.FloorSummary, error) {

	floors := make([]models.FloorSummary, 0)

	res, err := nosqlfr.client.
		Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(config.DynamoDBTable),
			KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", buildingId)},
				":sk": &types.AttributeValueMemberS{Value: "FLOORINFO#"},
			},
		})
	if err != nil {
		log.Println("Error details:", err.Error())
		return nil, errors.New("error fetching floors")
	}

	for _, item := range res.Items {
		var floor models.FloorSummary
		floor.BuildingID = uuid.MustParse(buildingId)
		floor.FloorNumber, _ = strconv.Atoi(item["FloorNumber"].(*types.AttributeValueMemberN).Value)
		floor.TotalSlots, _ = strconv.Atoi(item["TotalSlots"].(*types.AttributeValueMemberN).Value)
		floor.AvailableSlots, _ = strconv.Atoi(item["AvailableSlots"].(*types.AttributeValueMemberN).Value)
		if item["OfficeId"] != nil {
			officeId := item["OfficeId"].(*types.AttributeValueMemberS).Value

			officeQuery, err := nosqlfr.client.Query(ctx, &dynamodb.QueryInput{
				TableName:              aws.String(config.DynamoDBTable),
				KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":pk": &types.AttributeValueMemberS{Value: "OFFICE"},
					":sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("DETAILS#%s", officeId)},
				},
			})
			if err != nil {
				log.Println("Error details:", err.Error())
				return nil, errors.New("error fetching office")
			}
			floor.AssignedOffice = officeQuery.Items[0]["OfficeName"].(*types.AttributeValueMemberS).Value
		}
		floors = append(floors, floor)
	}
	return floors, nil
}

func NewNOSQLFloorRepository(client *dynamodb.Client) *NOSQLFloorRepository {
	return &NOSQLFloorRepository{
		client: client,
	}
}
