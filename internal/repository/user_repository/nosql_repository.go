package userrepository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Kaushik1766/ParkingManagement/internal/config"
	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	"github.com/Kaushik1766/ParkingManagement/internal/models/enums/roles"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type NOSQLUserRepository struct {
	client *dynamodb.Client
}

func NewNOSQLUserRepository(client *dynamodb.Client) *NOSQLUserRepository {
	return &NOSQLUserRepository{
		client: client,
	}
}

func (nosqlur *NOSQLUserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User

	queryRes, err := nosqlur.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		FilterExpression:       aws.String("IsActive = :active"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":     &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", email)},
			":sk":     &types.AttributeValueMemberS{Value: "PROFILE#"},
			":active": &types.AttributeValueMemberBOOL{Value: true},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return user, errors.New("error fetching user")
	}

	if len(queryRes.Items) == 0 {
		log.Println("user not found in GetUserByEmail")
		return user, errors.New("user not found")
	}

	item := queryRes.Items[0]
	user = nosqlur.itemToUser(item)

	return user, nil
}

func (nosqlur *NOSQLUserRepository) GetUserById(ctx context.Context, id string) (models.User, error) {
	var user models.User

	// TODO: fix scan, coz db schema updated
	scanRes, err := nosqlur.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(config.DynamoDBTable),
		FilterExpression: aws.String("Id = :id AND begins_with(SK, :sk) AND IsActive = :active"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id":     &types.AttributeValueMemberS{Value: id},
			":sk":     &types.AttributeValueMemberS{Value: "PROFILE#"},
			":active": &types.AttributeValueMemberBOOL{Value: true},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return user, errors.New("error fetching user")
	}

	if len(scanRes.Items) == 0 {
		log.Println("user not found in GetUserById")
		return user, errors.New("user not found")
	}

	item := scanRes.Items[0]
	user = nosqlur.itemToUser(item)

	return user, nil
}

func (nosqlur *NOSQLUserRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	// TODO: this function isnt needed much, for billing get only user ids, and delete the get all users admin route coz its not used
	scanRes, err := nosqlur.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(config.DynamoDBTable),
		FilterExpression: aws.String("begins_with(SK, :sk) AND IsActive = :active"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":sk":     &types.AttributeValueMemberS{Value: "PROFILE#"},
			":active": &types.AttributeValueMemberBOOL{Value: true},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("error fetching users")
	}

	for _, item := range scanRes.Items {
		user := nosqlur.itemToUser(item)
		users = append(users, user)
	}

	return users, nil
}

func (nosqlur *NOSQLUserRepository) Save(ctx context.Context, user models.User) error {
	updateExpression := "SET Username = :username, #role = :role, IsActive = :active"
	expressionValues := map[string]types.AttributeValue{
		":username": &types.AttributeValueMemberS{Value: user.Name},
		":role":     &types.AttributeValueMemberS{Value: user.Role.String()},
		":active":   &types.AttributeValueMemberBOOL{Value: user.IsActive},
	}
	expressionNames := map[string]string{
		"#role": "Role", // Role is a reserved word in DynamoDB
	}

	// Update office if provided
	if user.OfficeID != uuid.Nil {
		updateExpression += ", Office = :office, OfficeId = :officeId"
		expressionValues[":office"] = &types.AttributeValueMemberS{Value: user.Office.OfficeName}
		expressionValues[":officeId"] = &types.AttributeValueMemberS{Value: user.OfficeID.String()}
	}

	// Update password if provided
	if user.Password != "" {
		updateExpression += ", PasswordHash = :password"
		expressionValues[":password"] = &types.AttributeValueMemberS{Value: user.Password}
	}

	_, err := nosqlur.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", user.Email)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("PROFILE#%s", user.UserID.String())},
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionValues,
		ExpressionAttributeNames:  expressionNames,
	})
	if err != nil {
		log.Println(err.Error())
		return errors.New("error saving user")
	}

	return nil
}

func (nosqlur *NOSQLUserRepository) CreateUser(ctx context.Context, name, email, password, officeName string, role roles.Role) error {
	var officeId uuid.UUID

	if officeName != "" {
		// TODO: checking of office existence can be done using the OFFICE PK
		scanRes, err := nosqlur.client.Scan(ctx, &dynamodb.ScanInput{
			TableName:        aws.String(config.DynamoDBTable),
			FilterExpression: aws.String("Office = :office AND begins_with(SK, :sk)"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":office": &types.AttributeValueMemberS{Value: officeName},
				":sk":     &types.AttributeValueMemberS{Value: "FLOORINFO#"},
			},
		})
		if err != nil {
			log.Println(err.Error())
			return errors.New("error fetching office")
		}

		if len(scanRes.Items) > 0 {
			if scanRes.Items[0]["OfficeId"] != nil {
				officeId = uuid.MustParse(scanRes.Items[0]["OfficeId"].(*types.AttributeValueMemberS).Value)
			}
		}
	}

	userId := uuid.New()

	item := map[string]types.AttributeValue{
		"PK":           &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", email)},
		"SK":           &types.AttributeValueMemberS{Value: fmt.Sprintf("PROFILE#%s", userId.String())},
		"Email":        &types.AttributeValueMemberS{Value: email},
		"Id":           &types.AttributeValueMemberS{Value: userId.String()},
		"Username":     &types.AttributeValueMemberS{Value: name},
		"PasswordHash": &types.AttributeValueMemberS{Value: password},
		"Role":         &types.AttributeValueMemberS{Value: role.String()},
		"IsActive":     &types.AttributeValueMemberBOOL{Value: true},
	}

	if officeId != uuid.Nil {
		item["Office"] = &types.AttributeValueMemberS{Value: officeName}
		item["OfficeId"] = &types.AttributeValueMemberS{Value: officeId.String()}
	}

	_, err := nosqlur.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Item:      item,
	})
	if err != nil {
		log.Println(err.Error())
		return errors.New("error creating user")
	}

	return nil
}

// helper for dynamo to user
func (nosqlur *NOSQLUserRepository) itemToUser(item map[string]types.AttributeValue) models.User {
	var user models.User

	user.UserID = uuid.MustParse(item["Id"].(*types.AttributeValueMemberS).Value)
	user.Name = item["Username"].(*types.AttributeValueMemberS).Value
	user.Email = item["Email"].(*types.AttributeValueMemberS).Value

	if item["PasswordHash"] != nil {
		user.Password = item["PasswordHash"].(*types.AttributeValueMemberS).Value
	}

	roleStr := item["Role"].(*types.AttributeValueMemberS).Value
	switch roleStr {
	case "Admin":
		user.Role = roles.Admin
	case "Customer":
		user.Role = roles.Customer
	default:
		user.Role = roles.Customer
	}

	user.IsActive = item["IsActive"].(*types.AttributeValueMemberBOOL).Value

	// Get office info if present
	if item["Office"] != nil {
		user.Office.OfficeName = item["Office"].(*types.AttributeValueMemberS).Value
	}
	if item["OfficeId"] != nil {
		user.OfficeID = uuid.MustParse(item["OfficeId"].(*types.AttributeValueMemberS).Value)
		user.Office.OfficeID = user.OfficeID
	}

	return user
}
