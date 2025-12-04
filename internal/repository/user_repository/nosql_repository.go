package userrepository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

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

	// query in reverse lookup to get uuid
	lookupQuery, err := nosqlur.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk AND SK = :sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "USER"},
			":sk": &types.AttributeValueMemberS{Value: email},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return user, errors.New("error fetching user")
	}

	if len(lookupQuery.Items) == 0 {
		log.Println("user not found in GetUserByEmail")
		return user, errors.New("user not found")
	}

	id := lookupQuery.Items[0]["UUID"].(*types.AttributeValueMemberS).Value

	userQuery, err := nosqlur.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk and begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", id)},
			":sk": &types.AttributeValueMemberS{Value: "PROFILE"},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return user, errors.New("error fetching user")
	}

	item := userQuery.Items[0]
	user = nosqlur.itemToUser(ctx, item)

	return user, nil
}

func (nosqlur *NOSQLUserRepository) GetUserById(ctx context.Context, id string) (models.User, error) {
	var user models.User

	userQuery, err := nosqlur.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk and begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", id)},
			":sk": &types.AttributeValueMemberS{Value: "PROFILE"},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return user, errors.New("error fetching user")
	}

	if len(userQuery.Items) == 0 {
		log.Println("user not found in GetUserById")
		return user, errors.New("user not found")
	}

	item := userQuery.Items[0]
	user = nosqlur.itemToUser(ctx, item)

	return user, nil
}

func (nosqlur *NOSQLUserRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	// 1. Get all UUIDs from reverse lookup
	lookupQuery, err := nosqlur.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "USER"},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("error fetching user list")
	}

	var userKeys []map[string]types.AttributeValue
	for _, item := range lookupQuery.Items {
		if val, ok := item["UUID"]; ok {
			uid := val.(*types.AttributeValueMemberS).Value
			userKeys = append(userKeys, map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", uid)},
				"SK": &types.AttributeValueMemberS{Value: "PROFILE"},
			})
		}
	}

	// 2. BatchGetItem to fetch user profiles
	var userItems []map[string]types.AttributeValue
	batchSize := 100

	for i := 0; i < len(userKeys); i += batchSize {
		end := i + batchSize
		if end > len(userKeys) {
			end = len(userKeys)
		}
		batchKeys := userKeys[i:end]

		input := &dynamodb.BatchGetItemInput{
			RequestItems: map[string]types.KeysAndAttributes{
				config.DynamoDBTable: {
					Keys: batchKeys,
				},
			},
		}

		for {
			out, err := nosqlur.client.BatchGetItem(ctx, input)
			if err != nil {
				log.Println("BatchGetItem error:", err)
				return nil, err
			}
			userItems = append(userItems, out.Responses[config.DynamoDBTable]...)

			if len(out.UnprocessedKeys) > 0 {
				input.RequestItems = out.UnprocessedKeys
			} else {
				break
			}
		}
	}

	// 3. Process items in parallel (to keep Office fetch fast)
	var users []models.User
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, item := range userItems {
		wg.Add(1)
		go func(itm map[string]types.AttributeValue) {
			defer wg.Done()
			user := nosqlur.itemToUser(ctx, itm)
			// filter inactive users
			if user.IsActive {
				mu.Lock()
				users = append(users, user)
				mu.Unlock()
			}
		}(item)
	}

	wg.Wait()
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
		"#role": "Role",
	}

	// update office if provided
	if user.OfficeID != uuid.Nil {
		updateExpression += ", OfficeId = :officeId"
		expressionValues[":officeId"] = &types.AttributeValueMemberS{Value: user.OfficeID.String()}
	}

	// update password if provided
	if user.Password != "" {
		updateExpression += ", PasswordHash = :password"
		expressionValues[":password"] = &types.AttributeValueMemberS{Value: user.Password}
	}

	_, err := nosqlur.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(config.DynamoDBTable),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", user.UserID.String())},
			"SK": &types.AttributeValueMemberS{Value: "PROFILE"},
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

// DONE: get office id instead of name
func (nosqlur *NOSQLUserRepository) CreateUser(ctx context.Context, name, email, password, officeId string, role roles.Role) error {
	if officeId == "" {
		log.Println("no office id provided")
		return errors.New("office id is required")
	}

	// check if office exists
	officeQuery, err := nosqlur.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk AND SK = :sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "OFFICE"},
			":sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("DETAILS#%s", officeId)},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return errors.New("error fetching office")
	}
	if len(officeQuery.Items) == 0 {
		log.Printf("no office with id %s found", officeId)
		return errors.New("specified office not found")
	}

	// check if user exists
	userQuery, err := nosqlur.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk AND SK = :sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "USER"},
			":sk": &types.AttributeValueMemberS{Value: email},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return errors.New("error checking existing user")
	}
	if len(userQuery.Items) > 0 {
		log.Printf("user with email %s already exists", email)
		return errors.New("user with given email already exists")
	}

	userID := uuid.New().String()

	item := map[string]types.AttributeValue{
		"PK":           &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
		"SK":           &types.AttributeValueMemberS{Value: "PROFILE"},
		"Email":        &types.AttributeValueMemberS{Value: email},
		"Id":           &types.AttributeValueMemberS{Value: userID},
		"Username":     &types.AttributeValueMemberS{Value: name},
		"PasswordHash": &types.AttributeValueMemberS{Value: password},
		"Role":         &types.AttributeValueMemberS{Value: role.String()},
		"IsActive":     &types.AttributeValueMemberBOOL{Value: true},
		"OfficeId":     &types.AttributeValueMemberS{Value: officeId},
	}

	// transaction for putting value to rev lookup and user data
	_, err = nosqlur.client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					TableName: aws.String(config.DynamoDBTable),
					Item: map[string]types.AttributeValue{
						"PK":   &types.AttributeValueMemberS{Value: "USER"},
						"SK":   &types.AttributeValueMemberS{Value: email},
						"UUID": &types.AttributeValueMemberS{Value: userID},
					},
				},
			},
			{
				Put: &types.Put{
					TableName: aws.String(config.DynamoDBTable),
					Item:      item,
				},
			},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return errors.New("error creating user")
	}

	return nil
}

// helper for dynamo to user
func (nosqlur *NOSQLUserRepository) itemToUser(ctx context.Context, item map[string]types.AttributeValue) models.User {
	var user models.User

	user.UserID = uuid.MustParse(item["Id"].(*types.AttributeValueMemberS).Value)
	user.Name = item["Username"].(*types.AttributeValueMemberS).Value
	user.Email = item["Email"].(*types.AttributeValueMemberS).Value

	if item["PasswordHash"] != nil {
		user.Password = item["PasswordHash"].(*types.AttributeValueMemberS).Value
	}

	roleStr := item["Role"].(*types.AttributeValueMemberS).Value
	switch roleStr {
	case roles.Admin.String():
		user.Role = roles.Admin
	case roles.Customer.String():
		user.Role = roles.Customer
	default:
		user.Role = roles.Customer
	}

	user.IsActive = item["IsActive"].(*types.AttributeValueMemberBOOL).Value

	// if item["Office"] != nil {
	// 	user.Office.OfficeName = item["Office"].(*types.AttributeValueMemberS).Value
	// }
	if item["OfficeId"] != nil {
		user.OfficeID = uuid.MustParse(item["OfficeId"].(*types.AttributeValueMemberS).Value)
		user.Office.OfficeID = user.OfficeID

		res, err := nosqlur.client.Query(ctx, &dynamodb.QueryInput{
			TableName:              aws.String(config.DynamoDBTable),
			KeyConditionExpression: aws.String("PK = :pk AND SK = :sk"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk": &types.AttributeValueMemberS{Value: "OFFICE"},
				":sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("DETAILS#%s", user.OfficeID.String())},
			},
			AttributesToGet: []string{
				"OfficeName",
			},
		})
		if err == nil && len(res.Items) > 0 {
			user.Office.OfficeName = res.Items[0]["OfficeName"].(*types.AttributeValueMemberS).Value
		}
	}

	return user
}
