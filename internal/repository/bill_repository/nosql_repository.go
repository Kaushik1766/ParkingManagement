package billrepository

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Kaushik1766/ParkingManagement/internal/config"
	"github.com/Kaushik1766/ParkingManagement/internal/constants"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type NOSQLBillRepository struct {
	client *dynamodb.Client
}

func NewNOSQLBillRepository(client *dynamodb.Client) *NOSQLBillRepository {
	return &NOSQLBillRepository{
		client: client,
	}
}

func (nosqlbr *NOSQLBillRepository) SaveBill(ctx context.Context, bill models.BillDTO) error {
	// bill.UserId contains the user uuid
	pk := fmt.Sprintf("%s%s", constants.PrefixUser, bill.UserId)
	sk := fmt.Sprintf("%s%d#%d", constants.PKBill, bill.BillingYear, bill.BillingMonth)

	// Convert ParkingHistory to a list of maps for DynamoDB
	var parkingHistoryItems []types.AttributeValue
	for _, ph := range bill.ParkingHistory {
		phItem := map[string]types.AttributeValue{
			"TicketId":     &types.AttributeValueMemberS{Value: ph.TicketId},
			"NumberPlate":  &types.AttributeValueMemberS{Value: ph.NumberPlate},
			"BuildingId":   &types.AttributeValueMemberS{Value: ph.BuildingId},
			"BuildingName": &types.AttributeValueMemberS{Value: ph.BuildingName},
			"FloorNumber":  &types.AttributeValueMemberN{Value: strconv.Itoa(ph.FLoorNumber)},
			"SlotNumber":   &types.AttributeValueMemberN{Value: strconv.Itoa(ph.SlotNumber)},
			"VehicleType":  &types.AttributeValueMemberS{Value: ph.VechicleType},
			"StartTime":    &types.AttributeValueMemberN{Value: strconv.FormatInt(ph.StartTime.Unix(), 10)},
		}
		if !ph.EndTime.IsZero() {
			phItem["EndTime"] = &types.AttributeValueMemberN{Value: strconv.FormatInt(ph.EndTime.Unix(), 10)}
		}
		parkingHistoryItems = append(parkingHistoryItems, &types.AttributeValueMemberM{Value: phItem})
	}

	item := map[string]types.AttributeValue{
		"PK":             &types.AttributeValueMemberS{Value: pk},
		"SK":             &types.AttributeValueMemberS{Value: sk},
		"TotalAmount":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", bill.TotalAmount)},
		"BillDate":       &types.AttributeValueMemberS{Value: bill.BillDate},
		"BillingMonth":   &types.AttributeValueMemberN{Value: strconv.Itoa(bill.BillingMonth)},
		"BillingYear":    &types.AttributeValueMemberN{Value: strconv.Itoa(bill.BillingYear)},
		"ParkingHistory": &types.AttributeValueMemberL{Value: parkingHistoryItems},
	}

	_, err := nosqlbr.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Item:      item,
	})
	if err != nil {
		log.Println("Error saving bill:", err)
		return err
	}

	return nil
}

func (nosqlbr *NOSQLBillRepository) GetBill(ctx context.Context, userId string, month, year int) (models.BillDTO, error) {
	pk := fmt.Sprintf("%s%s", constants.PrefixUser, userId)
	sk := fmt.Sprintf("%s%d#%d", constants.PKBill, year, month)

	res, err := nosqlbr.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	})
	if err != nil {
		log.Println("Error fetching bill:", err)
		return models.BillDTO{}, err
	}

	if res.Item == nil {
		return models.BillDTO{}, nil
	}

	var bill models.BillDTO
	bill.UserId = userId
	bill.BillDate = res.Item["BillDate"].(*types.AttributeValueMemberS).Value
	if res.Item["BillingMonth"] != nil {
		bill.BillingMonth, _ = strconv.Atoi(res.Item["BillingMonth"].(*types.AttributeValueMemberN).Value)
	} else {
		// Fallback if not present (legacy)
		bill.BillingMonth = month
	}
	if res.Item["BillingYear"] != nil {
		bill.BillingYear, _ = strconv.Atoi(res.Item["BillingYear"].(*types.AttributeValueMemberN).Value)
	} else {
		bill.BillingYear = year
	}
	totalAmountStr := res.Item["TotalAmount"].(*types.AttributeValueMemberN).Value
	bill.TotalAmount, _ = strconv.ParseFloat(totalAmountStr, 64)

	parkingHistoryItems := res.Item["ParkingHistory"].(*types.AttributeValueMemberL).Value
	for _, item := range parkingHistoryItems {
		phMap := item.(*types.AttributeValueMemberM).Value
		var ph models.ParkingHistoryDTO
		ph.TicketId = phMap["TicketId"].(*types.AttributeValueMemberS).Value
		ph.NumberPlate = phMap["NumberPlate"].(*types.AttributeValueMemberS).Value
		ph.BuildingId = phMap["BuildingId"].(*types.AttributeValueMemberS).Value
		ph.BuildingName = phMap["BuildingName"].(*types.AttributeValueMemberS).Value
		ph.FLoorNumber, _ = strconv.Atoi(phMap["FloorNumber"].(*types.AttributeValueMemberN).Value)
		ph.SlotNumber, _ = strconv.Atoi(phMap["SlotNumber"].(*types.AttributeValueMemberN).Value)
		ph.VechicleType = phMap["VehicleType"].(*types.AttributeValueMemberS).Value

		startTimeUnix, _ := strconv.ParseInt(phMap["StartTime"].(*types.AttributeValueMemberN).Value, 10, 64)
		ph.StartTime = time.Unix(startTimeUnix, 0).Local()

		if phMap["EndTime"] != nil {
			endTimeUnix, _ := strconv.ParseInt(phMap["EndTime"].(*types.AttributeValueMemberN).Value, 10, 64)
			ph.EndTime = time.Unix(endTimeUnix, 0).Local()
		}

		bill.ParkingHistory = append(bill.ParkingHistory, ph)
	}

	return bill, nil
}
