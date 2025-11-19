package buildingrepository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/Kaushik1766/ParkingManagement/internal/config"
	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type NOSQLBuidlingRepository struct {
	client *dynamodb.Client
}

func NewNOSQLBuidlingRepository(client *dynamodb.Client) *NOSQLBuidlingRepository {
	return &NOSQLBuidlingRepository{
		client: client,
	}
}

func (nosqlbr *NOSQLBuidlingRepository) DeleteBuildingByID(ctx context.Context, buildingID string) error {

	out, err := nosqlbr.
		client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{
				Value: "BUILDING",
			},
			"SK": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("BUILDING#%s", buildingID),
			},
		},
	})

	fmt.Printf("%+v\n", out)

	return err
}

// func (nosqlbr *NOSQLBuidlingRepository) GetBuildingByName(buildingName string) (models.Building, error) {
// if buildingName == constants.AdminBuilding {
// 	return models.Building{}, errors.New("buildingrepo: cannot get admin building")
// }
// building := models.Building{}
// err := nosqlbr.db.Where("building_name = ?", buildingName).First(&building).Error
// return building, err
// 	panic("not implemented")
// }

func (nosqlbr *NOSQLBuidlingRepository) GetAllBuildingSummary(ctx context.Context) ([]models.BuildingSummary, error) {
	var buildings []models.BuildingSummary

	res, err := nosqlbr.client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
		Statement: aws.String(`SELECT * FROM "` + config.DynamoDBTable + `" WHERE PK = 'BUILDING'`),
	})
	if err != nil {
		return buildings, err
	}

	for _, item := range res.Items {
		building := models.BuildingSummary{}
		building.BuildingId = uuid.MustParse(item["BuildingId"].(*types.AttributeValueMemberS).Value)
		building.BuildingName = item["BuildingName"].(*types.AttributeValueMemberS).Value
		n, err := strconv.Atoi(item["TotalFloors"].(*types.AttributeValueMemberN).Value)
		if err != nil {
			return buildings, err
		}
		building.TotalFloors = n
		n, err = strconv.Atoi(item["TotalSlots"].(*types.AttributeValueMemberN).Value)
		if err != nil {
			return buildings, err
		}
		building.TotalSlots = n

		n, err = strconv.Atoi(item["AvailableSlots"].(*types.AttributeValueMemberN).Value)
		if err != nil {
			return buildings, err
		}
		building.AvailableSlots = n
		buildings = append(buildings, building)
	}
	return buildings, nil
}

// unused in project
func (nosqlbr *NOSQLBuidlingRepository) GetAllBuildings(ctx context.Context) ([]models.Building, error) {
	// var buildings []models.Building
	// err := nosqlbr.db.
	// 	Where("building_name <> ?", constants.AdminBuilding).
	// 	Preload("Floors.Slots.Vehicles").
	// 	Find(&buildings).Error
	// return buildings, err
	panic("not implemented")
}

func (nosqlbr *NOSQLBuidlingRepository) GetBuildingByID(ctx context.Context, buildingID uuid.UUID) (models.Building, error) {
	building := models.Building{}

	res, err := nosqlbr.client.
		GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(config.DynamoDBTable),
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{
					Value: "BUILDING",
				},
				"SK": &types.AttributeValueMemberS{
					Value: fmt.Sprintf("BUILDING#%s", buildingID.String()),
				},
			},
		})
	if err != nil {
		return building, err
	}

	if res.Item == nil {
		return building, errors.New("building not found")
	}

	building.BuildingID = buildingID
	building.BuildingName = res.Item["BuildingName"].(*types.AttributeValueMemberS).Value

	// building := models.Building{}
	// err := nosqlbr.db.
	// 	Where("building_id = ? AND building_name <> ?", buildingID, constants.AdminBuilding).
	// 	First(&building).Error
	// return building, err
	return building, nil
}

func (nosqlbr *NOSQLBuidlingRepository) AddBuilding(ctx context.Context, buildingName string) error {
	if buildingName == constants.AdminBuilding {
		return errors.New("buildingrepo: cannot add admin building")
	}
	building := models.BuildingSummary{
		BuildingName:   buildingName,
		BuildingId:     uuid.New(),
		AvailableSlots: 0,
		TotalSlots:     0,
		TotalFloors:    0,
	}

	item, err := attributevalue.MarshalMap(building)
	if err != nil {
		return err
	}

	item["PK"] = &types.AttributeValueMemberS{Value: "BUILDING"}
	item["SK"] = &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", building.BuildingId.String())}
	item["BuildingId"] = &types.AttributeValueMemberS{Value: building.BuildingId.String()}

	_, err = nosqlbr.client.
		PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(config.DynamoDBTable),
			Item:      item,
		})
	if err != nil {
		log.Println(err.Error())
		return errors.New("error adding building")
	}

	return nil
}
