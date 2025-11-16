package officerepository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/Kaushik1766/ParkingManagement/internal/config"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type NOSQLOfficeRepository struct {
	client *dynamodb.Client
}

func NewNOSQLOfficeRepository(client *dynamodb.Client) *NOSQLOfficeRepository {
	return &NOSQLOfficeRepository{
		client: client,
	}
}

func (nosqlor *NOSQLOfficeRepository) AddOffice(officeName string, buildingID string, floorNumber int) error {
	_, err := uuid.Parse(buildingID)
	if err != nil {
		log.Println(err.Error())
		return errors.New("invalid building ID")
	}
	// office := models.Office{
	// 	BuildingID:  buildingUUID,
	// 	FloorNumber: floorNumber,
	// 	OfficeName:  officeName,
	// }

	_, err = nosqlor.client.UpdateItem(context.Background(), &dynamodb.UpdateItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", buildingID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("FLOORINFO#%d", floorNumber)},
		},
		UpdateExpression: aws.String("ADD Office :val"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":val": &types.AttributeValueMemberSS{Value: []string{officeName}},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return errors.New("error adding office")
	}

	return nil
}

func (nosqlor *NOSQLOfficeRepository) DeleteOffice(officeId string) error {
	officeUUID, err := uuid.Parse(officeId)
	if err != nil {
		return err
	}

	// if officeName == constants.AdminOffice {
	// 	return errors.New("officerepo: cannot delete admin office")
	// }
	return nosqlor.db.Where("office_id = ?", officeUUID).Delete(&models.Office{}).Error
}

func (nosqlor *NOSQLOfficeRepository) GetBuildingAndFloorByOffice(officeName string) (uuid.UUID, int, error) {
	panic("not implemented coz not used")
}

func (nosqlor *NOSQLOfficeRepository) GetOfficesByBuilding(buildingID string) ([]models.Office, error) {
	var offices []models.Office
	_, err := uuid.Parse(buildingID)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("invalid building id")
	}

	// err = nosqlor.db.Where("building_id = ? AND office_name <> ?", buildingUUID, constants.AdminOffice).Find(&offices).Error
	items, err := nosqlor.client.Query(context.Background(), &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :prefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":     &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", buildingID)},
			":prefix": &types.AttributeValueMemberS{Value: "FLOORINFO#"},
		},
	})
	if err != nil {
		return nil, err
	}

	for _, item := range items.Items {
		var office models.Office
		office.BuildingID = uuid.MustParse(buildingID)
		office.FloorNumber, _ = strconv.Atoi(item["FloorNumber"].(*types.AttributeValueMemberN).Value)
		office.OfficeName = item["Office"].(*types.AttributeValueMemberS).Value
		offices = append(offices, office)
	}
	return offices, nil
}

func (nosqlor *NOSQLOfficeRepository) GetAllOffices() ([]models.Office, error) {
	var offices []models.Office

	buildings, err := nosqlor.client.Query(context.Background(), &dynamodb.QueryInput{
		TableName:        aws.String(config.DynamoDBTable),
		FilterExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "BUILDING"},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	for _, building := range buildings.Items {
		buildingID := building["BuildingID"].(*types.AttributeValueMemberS).Value

	}

	// err = nosqlor.db.Where("building_id = ? AND office_name <> ?", buildingUUID, constants.AdminOffice).Find(&offices).Error
	items, err := nosqlor.client.Query(context.Background(), &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :prefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":     &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", buildingID)},
			":prefix": &types.AttributeValueMemberS{Value: "FLOORINFO#"},
		},
	})
	if err != nil {
		return nil, err
	}

	for _, item := range items.Items {
		var office models.Office
		office.BuildingID = uuid.MustParse(buildingID)
		office.FloorNumber, _ = strconv.Atoi(item["FloorNumber"].(*types.AttributeValueMemberN).Value)
		office.OfficeName = item["Office"].(*types.AttributeValueMemberS).Value
		offices = append(offices, office)
	}
	return offices, nil
}

func (nosqlor *NOSQLOfficeRepository) GetOfficeByName(officeName string) (models.Office, error) {
	// if officeName == constants.AdminOffice {
	// 	return models.Office{}, errors.New("officerepo: cannot get admin office by name")
	// }
	var office models.Office
	err := nosqlor.db.Where("office_name = ?", officeName).First(&office).Error
	return office, err
}
