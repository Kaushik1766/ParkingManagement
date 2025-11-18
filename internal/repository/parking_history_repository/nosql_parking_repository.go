package parkinghistoryrepository

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Kaushik1766/ParkingManagement/internal/config"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

const (
	userPKPrefix    = "USER#"
	parkingSKPrefix = "PARKING#"
)

var (
	errParkingNotFound      = errors.New("parkingrepo: parking record not found")
	errVehicleAlreadyParked = errors.New("parkingrepo: vehicle is already parked")
)

// NOSQLParkingRepository persists parking history entries in DynamoDB following the documented single-table design.
type NOSQLParkingRepository struct {
	client         *dynamodb.Client
	buildingCache  map[string]string
	buildingCacheM sync.RWMutex
}

// NewNOSQLParkingRepository builds a ParkingHistoryStorage backed by DynamoDB.
func NewNOSQLParkingRepository(client *dynamodb.Client) *NOSQLParkingRepository {
	return &NOSQLParkingRepository{
		client:        client,
		buildingCache: map[string]string{},
	}
}

func (repo *NOSQLParkingRepository) AddParking(vehicle models.Vehicle) (string, error) {
	if repo.client == nil {
		return "", errors.New("parkingrepo: dynamodb client is nil")
	}

	if vehicle.UserID == uuid.Nil {
		return "", errors.New("parkingrepo: vehicle is not linked to a user")
	}

	if vehicle.AssignedBuildingID == uuid.Nil || vehicle.AssignedFloorNumber == 0 || vehicle.AssignedSlotNumber == 0 {
		return "", errors.New("parkingrepo: vehicle is not assigned to a slot")
	}

	normalizedPlate := normalizePlate(vehicle.NumberPlate)
	ctx := context.Background()

	// Ensure the same vehicle is not already parked.
	if _, err := repo.findActiveParkingByNumberPlate(ctx, normalizedPlate); err == nil {
		return "", errVehicleAlreadyParked
	} else if err != nil && !errors.Is(err, errParkingNotFound) {
		return "", err
	}

	ticketID := uuid.NewString()
	now := time.Now().UTC()

	buildingName, err := repo.getBuildingName(ctx, vehicle.AssignedBuildingID.String())
	if err != nil {
		return "", err
	}

	item := map[string]types.AttributeValue{
		"PK":             &types.AttributeValueMemberS{Value: userPK(vehicle.UserID.String())},
		"SK":             &types.AttributeValueMemberS{Value: buildSortKey(now, ticketID)},
		"ParkingId":      &types.AttributeValueMemberS{Value: ticketID},
		"UserId":         &types.AttributeValueMemberS{Value: vehicle.UserID.String()},
		"UserEmail":      &types.AttributeValueMemberS{Value: strings.ToLower(strings.TrimSpace(vehicle.User.Email))},
		"Numberplate":    &types.AttributeValueMemberS{Value: normalizedPlate},
		"BuildingId":     &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", vehicle.AssignedBuildingID.String())},
		"BuildingUuid":   &types.AttributeValueMemberS{Value: vehicle.AssignedBuildingID.String()},
		"BuildingName":   &types.AttributeValueMemberS{Value: buildingName},
		"FloorNumber":    &types.AttributeValueMemberN{Value: strconv.Itoa(vehicle.AssignedFloorNumber)},
		"SlotNumber":     &types.AttributeValueMemberN{Value: strconv.Itoa(vehicle.AssignedSlotNumber)},
		"SlotId":         &types.AttributeValueMemberS{Value: fmt.Sprintf("SLOT#F%d#%d", vehicle.AssignedFloorNumber, vehicle.AssignedSlotNumber)},
		"VehicleType":    &types.AttributeValueMemberS{Value: vehicle.VehicleType.String()},
		"StartTime":      &types.AttributeValueMemberS{Value: now.Format(time.RFC3339)},
		"StartTimestamp": &types.AttributeValueMemberN{Value: strconv.FormatInt(now.UnixMilli(), 10)},
	}

	_, err = repo.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(config.DynamoDBTable),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(PK) AND attribute_not_exists(SK)"),
	})
	if err != nil {
		return "", err
	}

	return ticketID, nil
}

func (repo *NOSQLParkingRepository) Unpark(ticketID string) error {
	ctx := context.Background()

	item, err := repo.findParkingByTicketID(ctx, ticketID)
	if err != nil {
		return err
	}

	return repo.closeParking(ctx, item.PK, item.SK)
}

func (repo *NOSQLParkingRepository) UnparkByNumberPlate(numberplate string) error {
	ctx := context.Background()
	normalizedPlate := normalizePlate(numberplate)

	item, err := repo.findActiveParkingByNumberPlate(ctx, normalizedPlate)
	if err != nil {
		return err
	}

	return repo.closeParking(ctx, item.PK, item.SK)
}

func (repo *NOSQLParkingRepository) GetParkingHistoryByNumberPlate(numberplate string, startTime, endTime time.Time) ([]models.ParkingHistoryDTO, error) {
	ctx := context.Background()
	normalizedPlate := normalizePlate(numberplate)

	stmt := fmt.Sprintf(`SELECT * FROM "%s" WHERE Numberplate = ?`, config.DynamoDBTable)

	res, err := repo.client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
		Statement: aws.String(stmt),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: normalizedPlate},
		},
	})
	if err != nil {
		return nil, err
	}

	histories := make([]models.ParkingHistoryDTO, 0, len(res.Items))

	for _, av := range res.Items {
		record, err := repo.unmarshalItem(av)
		if err != nil {
			return nil, err
		}

		if record.EndTime == nil {
			continue
		}

		dto, err := repo.toDTO(record)
		if err != nil {
			return nil, err
		}

		if dto.StartTime.Before(startTime) || dto.EndTime.After(endTime) {
			continue
		}

		histories = append(histories, dto)
	}

	sort.SliceStable(histories, func(i, j int) bool {
		return histories[i].StartTime.After(histories[j].StartTime)
	})

	return histories, nil
}

func (repo *NOSQLParkingRepository) GetParkingHistoryByUser(userID string, startTime, endTime time.Time) ([]models.ParkingHistoryDTO, error) {
	ctx := context.Background()

	res, err := repo.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk AND SK BETWEEN :start AND :end"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":    &types.AttributeValueMemberS{Value: userPK(userID)},
			":start": &types.AttributeValueMemberS{Value: buildRangeKey(startTime, false)},
			":end":   &types.AttributeValueMemberS{Value: buildRangeKey(endTime, true)},
		},
	})
	if err != nil {
		return nil, err
	}

	histories := make([]models.ParkingHistoryDTO, 0, len(res.Items))

	for _, av := range res.Items {
		record, err := repo.unmarshalItem(av)
		if err != nil {
			return nil, err
		}

		if record.EndTime == nil {
			continue
		}

		dto, err := repo.toDTO(record)
		if err != nil {
			return nil, err
		}

		histories = append(histories, dto)
	}

	sort.SliceStable(histories, func(i, j int) bool {
		return histories[i].StartTime.After(histories[j].StartTime)
	})

	return histories, nil
}

func (repo *NOSQLParkingRepository) GetActiveUserParkings(userID string) ([]models.ParkingHistoryDTO, error) {
	ctx := context.Background()

	res, err := repo.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk"),
		FilterExpression:       aws.String("attribute_not_exists(EndTime)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: userPK(userID)},
		},
	})
	if err != nil {
		return nil, err
	}

	active := make([]models.ParkingHistoryDTO, 0, len(res.Items))

	for _, av := range res.Items {
		record, err := repo.unmarshalItem(av)
		if err != nil {
			return nil, err
		}

		dto, err := repo.toDTO(record)
		if err != nil {
			return nil, err
		}

		active = append(active, dto)
	}

	sort.SliceStable(active, func(i, j int) bool {
		return active[i].StartTime.After(active[j].StartTime)
	})

	return active, nil
}

// parkingHistoryItem mirrors the attribute schema for a parking record.
type parkingHistoryItem struct {
	PK             string  `dynamodbav:"PK"`
	SK             string  `dynamodbav:"SK"`
	ParkingID      string  `dynamodbav:"ParkingId"`
	UserID         string  `dynamodbav:"UserId"`
	UserEmail      string  `dynamodbav:"UserEmail"`
	NumberPlate    string  `dynamodbav:"Numberplate"`
	BuildingID     string  `dynamodbav:"BuildingId"`
	BuildingUUID   string  `dynamodbav:"BuildingUuid"`
	BuildingName   string  `dynamodbav:"BuildingName"`
	FloorNumber    int     `dynamodbav:"FloorNumber"`
	SlotNumber     int     `dynamodbav:"SlotNumber"`
	SlotID         string  `dynamodbav:"SlotId"`
	VehicleType    string  `dynamodbav:"VehicleType"`
	StartTime      string  `dynamodbav:"StartTime"`
	StartTimestamp int64   `dynamodbav:"StartTimestamp"`
	EndTime        *string `dynamodbav:"EndTime"`
	EndTimestamp   *int64  `dynamodbav:"EndTimestamp"`
}

func (repo *NOSQLParkingRepository) findParkingByTicketID(ctx context.Context, ticketID string) (parkingHistoryItem, error) {
	stmt := fmt.Sprintf(`SELECT PK, SK, ParkingId, EndTime FROM "%s" WHERE ParkingId = ?`, config.DynamoDBTable)

	res, err := repo.client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
		Statement: aws.String(stmt),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: ticketID},
		},
	})
	if err != nil {
		return parkingHistoryItem{}, err
	}

	if len(res.Items) == 0 {
		return parkingHistoryItem{}, errParkingNotFound
	}

	return repo.unmarshalItem(res.Items[0])
}

func (repo *NOSQLParkingRepository) findActiveParkingByNumberPlate(ctx context.Context, numberplate string) (parkingHistoryItem, error) {
	stmt := fmt.Sprintf(`SELECT * FROM "%s" WHERE Numberplate = ?`, config.DynamoDBTable)

	res, err := repo.client.ExecuteStatement(ctx, &dynamodb.ExecuteStatementInput{
		Statement: aws.String(stmt),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: numberplate},
		},
	})
	if err != nil {
		return parkingHistoryItem{}, err
	}

	for _, item := range res.Items {
		record, err := repo.unmarshalItem(item)
		if err != nil {
			return parkingHistoryItem{}, err
		}

		if record.EndTime == nil {
			return record, nil
		}
	}

	return parkingHistoryItem{}, errParkingNotFound
}

func (repo *NOSQLParkingRepository) closeParking(ctx context.Context, pk, sk string) error {
	now := time.Now().UTC()

	_, err := repo.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
		UpdateExpression: aws.String("SET EndTime = :end, EndTimestamp = :ts"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":end": &types.AttributeValueMemberS{Value: now.Format(time.RFC3339)},
			":ts":  &types.AttributeValueMemberN{Value: strconv.FormatInt(now.UnixMilli(), 10)},
		},
		ConditionExpression: aws.String("attribute_not_exists(EndTime)"),
	})
	if err != nil {
		var condErr *types.ConditionalCheckFailedException
		if errors.As(err, &condErr) {
			return errors.New("parkingrepo: parking already closed")
		}
		return err
	}

	return nil
}

func (repo *NOSQLParkingRepository) getBuildingName(ctx context.Context, buildingID string) (string, error) {
	repo.buildingCacheM.RLock()
	name, ok := repo.buildingCache[buildingID]
	repo.buildingCacheM.RUnlock()
	if ok {
		return name, nil
	}

	res, err := repo.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "BUILDING"},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("BUILDING#%s", buildingID)},
		},
		ProjectionExpression: aws.String("BuildingName"),
	})
	if err != nil {
		return "", err
	}

	if res.Item == nil {
		return "", errors.New("parkingrepo: building not found")
	}

	raw, ok := res.Item["BuildingName"].(*types.AttributeValueMemberS)
	if !ok {
		return "", errors.New("parkingrepo: invalid building data")
	}

	repo.buildingCacheM.Lock()
	repo.buildingCache[buildingID] = raw.Value
	repo.buildingCacheM.Unlock()

	return raw.Value, nil
}

func (repo *NOSQLParkingRepository) unmarshalItem(av map[string]types.AttributeValue) (parkingHistoryItem, error) {
	var item parkingHistoryItem
	if err := attributevalue.UnmarshalMap(av, &item); err != nil {
		return parkingHistoryItem{}, err
	}
	return item, nil
}

func (repo *NOSQLParkingRepository) toDTO(item parkingHistoryItem) (models.ParkingHistoryDTO, error) {
	dto := models.ParkingHistoryDTO{
		TicketId:     item.ParkingID,
		NumberPlate:  item.NumberPlate,
		BuildingId:   sanitizeBuildingID(item),
		BuildingName: item.BuildingName,
		FLoorNumber:  item.FloorNumber,
		SlotNumber:   item.SlotNumber,
		VechicleType: item.VehicleType,
	}

	start, err := time.Parse(time.RFC3339, item.StartTime)
	if err != nil {
		return models.ParkingHistoryDTO{}, err
	}
	dto.StartTime = start.Local()

	if item.EndTime != nil {
		end, err := time.Parse(time.RFC3339, *item.EndTime)
		if err != nil {
			return models.ParkingHistoryDTO{}, err
		}
		dto.EndTime = end.Local()
	}

	return dto, nil
}

func sanitizeBuildingID(item parkingHistoryItem) string {
	if item.BuildingUUID != "" {
		return item.BuildingUUID
	}
	return strings.TrimPrefix(item.BuildingID, "BUILDING#")
}

func normalizePlate(numberplate string) string {
	return strings.ToUpper(strings.TrimSpace(numberplate))
}

func userPK(userID string) string {
	return userPKPrefix + userID
}

func buildSortKey(t time.Time, ticketID string) string {
	return fmt.Sprintf("%s%020d#%s", parkingSKPrefix, t.UTC().UnixMilli(), ticketID)
}

func buildRangeKey(t time.Time, upper bool) string {
	suffix := ""
	if upper {
		suffix = "~"
	}
	return fmt.Sprintf("%s%020d%s", parkingSKPrefix, t.UTC().UnixMilli(), suffix)
}
