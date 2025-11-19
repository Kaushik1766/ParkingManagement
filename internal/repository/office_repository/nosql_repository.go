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

func (nosqlor *NOSQLOfficeRepository) AddOffice(ctx context.Context, officeName string, buildingID string, floorNumber int) error {
	_, err := uuid.Parse(buildingID)
	if err != nil {
		log.Println("invalid building ID:", err.Error())
		return errors.New("invalid building ID")
	}
	// office := models.Office{
	// 	BuildingID:  buildingUUID,
	// 	FloorNumber: floorNumber,
	// 	OfficeName:  officeName,
	// }

	_, err = nosqlor.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", buildingID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("FLOORINFO#%d", floorNumber)},
		},
		UpdateExpression: aws.String("SET Office = :val, OfficeId = :id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":val": &types.AttributeValueMemberS{Value: officeName},
			":id":  &types.AttributeValueMemberS{Value: uuid.NewString()},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return errors.New("error adding office")
	}

	return nil
}

func (nosqlor *NOSQLOfficeRepository) DeleteOffice(ctx context.Context, officeId string) error {
	panic("pending implementation")
}

func (nosqlor *NOSQLOfficeRepository) GetBuildingAndFloorByOffice(ctx context.Context, officeName string) (uuid.UUID, int, error) {
	panic("not implemented coz not used")
}

func (nosqlor *NOSQLOfficeRepository) GetOfficesByBuilding(ctx context.Context, buildingID string) ([]models.Office, error) {
	var offices []models.Office
	_, err := uuid.Parse(buildingID)
	if err != nil {
		log.Println("invalid building id:", err.Error())
		return nil, errors.New("invalid building id")
	}

	// err = nosqlor.db.Where("building_id = ? AND office_name <> ?", buildingUUID, constants.AdminOffice).Find(&offices).Error
	items, err := nosqlor.client.Query(ctx, &dynamodb.QueryInput{
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
		if item["Office"] == nil {
			continue
		}
		office.OfficeName = item["Office"].(*types.AttributeValueMemberS).Value
		if item["OfficeId"] == nil {
			continue
		}
		office.OfficeID = uuid.MustParse(item["OfficeId"].(*types.AttributeValueMemberS).Value)
		offices = append(offices, office)
	}
	return offices, nil
}

func (nosqlor *NOSQLOfficeRepository) GetAllOffices(ctx context.Context) ([]models.Office, error) {
	var offices []models.Office

	buildings, err := nosqlor.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "BUILDING"},
			":sk": &types.AttributeValueMemberS{Value: "BUILDING#"},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("failed to fetch buildings")
	}

	for _, building := range buildings.Items {
		buildingID := building["BuildingId"].(*types.AttributeValueMemberS).Value
		fmt.Println(buildingID)
		res, err := nosqlor.client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(config.DynamoDBTable),
			KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :prefix)"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk":     &types.AttributeValueMemberS{Value: "BUILDING#" + buildingID},
				":prefix": &types.AttributeValueMemberS{Value: "FLOORINFO#"},
			},
			ProjectionExpression: aws.String("FloorNumber, Office, OfficeId"),
		})
		if err != nil {
			return nil, err
		}

		for _, item := range res.Items {
			var office models.Office
			office.BuildingID = uuid.MustParse(buildingID)
			office.FloorNumber, _ = strconv.Atoi(item["FloorNumber"].(*types.AttributeValueMemberN).Value)

			if item["Office"] == nil {
				continue
			}
			office.OfficeName = item["Office"].(*types.AttributeValueMemberS).Value

			if item["OfficeId"] == nil {
				continue
			}
			office.OfficeID = uuid.MustParse(item["OfficeId"].(*types.AttributeValueMemberS).Value)

			offices = append(offices, office)
		}
	}

	return offices, nil
}

func (nosqlor *NOSQLOfficeRepository) GetOfficeByName(ctx context.Context, officeName string) (models.Office, error) {
	var office models.Office

	// Need to scan since offices are stored with PK=BUILDING#{id}, not PK=BUILDING
	scanRes, err := nosqlor.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(config.DynamoDBTable),
		FilterExpression: aws.String("begins_with(SK, :sk) AND Office = :officeName"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":officeName": &types.AttributeValueMemberS{Value: officeName},
			":sk":         &types.AttributeValueMemberS{Value: "FLOORINFO#"},
		},
	})
	if err != nil {
		log.Println("error in GetOfficeByName scan:", err.Error())
		return models.Office{}, err
	}

	if len(scanRes.Items) == 0 {
		log.Printf("No office found with name: %s", officeName)
		return models.Office{}, errors.New("office not found")
	}

	item := scanRes.Items[0]

	// Extract BuildingID from PK (format: BUILDING#{uuid})
	pkValue := item["PK"].(*types.AttributeValueMemberS).Value
	buildingIdStr := pkValue[9:] // Remove "BUILDING#" prefix
	office.BuildingID = uuid.MustParse(buildingIdStr)

	office.FloorNumber, _ = strconv.Atoi(item["FloorNumber"].(*types.AttributeValueMemberN).Value)
	office.OfficeName = item["Office"].(*types.AttributeValueMemberS).Value
	office.OfficeID = uuid.MustParse(item["OfficeId"].(*types.AttributeValueMemberS).Value)

	log.Printf("Found office: %s, BuildingID: %s, Floor: %d", office.OfficeName, office.BuildingID, office.FloorNumber)

	return office, nil
}
