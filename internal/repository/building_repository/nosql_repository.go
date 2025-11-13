package buildingrepository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Kaushik1766/ParkingManagement/internal/config"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
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

func (nosqlbr *NOSQLBuidlingRepository) DeleteBuildingByID(buildingID string) error {

	out, err := nosqlbr.
		client.DeleteItem(context.Background(), &dynamodb.DeleteItemInput{
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

func (nosqlbr *NOSQLBuidlingRepository) GetAllBuildingSummary() ([]models.BuildingSummary, error) {
	var buildings []models.BuildingSummary

	res, err := nosqlbr.client.ExecuteStatement(context.Background(), &dynamodb.ExecuteStatementInput{
		Statement: aws.String(`select * from "pm_nosql" where "PK"="BUILDING"`),
	})

	res1, err := nosqlbr.client.
		GetItem(context.Background(), &dynamodb.GetItemInput{
			
		})
	
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("error fetching buildings")
	}
	return buildings, nil
}

func (nosqlbr *NOSQLBuidlingRepository) GetAllBuildings() ([]models.Building, error) {
	// var buildings []models.Building
	// err := nosqlbr.db.
	// 	Where("building_name <> ?", constants.AdminBuilding).
	// 	Preload("Floors.Slots.Vehicles").
	// 	Find(&buildings).Error
	// return buildings, err
	panic("not implemented")
}

func (nosqlbr *NOSQLBuidlingRepository) GetBuildingByID(buildingID uuid.UUID) (models.Building, error) {
	// building := models.Building{}
	// err := nosqlbr.db.
	// 	Where("building_id = ? AND building_name <> ?", buildingID, constants.AdminBuilding).
	// 	First(&building).Error
	// return building, err
	panic("not implemented")
}

func (nosqlbr *NOSQLBuidlingRepository) AddBuilding(buildingName string) error {
	// if buildingName == constants.AdminBuilding {
	// 	return errors.New("buildingrepo: cannot add admin building")
	// }
	// building := models.Building{
	// 	BuildingName: buildingName,
	// 	Floors:       nil,
	// }
	// err := nosqlbr.db.Create(&building).Error
	// if err != nil {
	// 	return errors.New("building name should be unique")
	// }
	// return nil
	panic("not implemented")
}
