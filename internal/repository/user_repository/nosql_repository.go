package userrepository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/Kaushik1766/ParkingManagement/internal/config"
	"github.com/Kaushik1766/ParkingManagement/internal/constants"
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
			":pk": &types.AttributeValueMemberS{Value: constants.PKUser},
			":sk": &types.AttributeValueMemberS{Value: email},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return user, errors.New(constants.ErrFetchingUser)
	}

	if len(lookupQuery.Items) == 0 {
		log.Println("user not found in GetUserByEmail")
		return user, errors.New(constants.ErrUserNotFound)
	}

	id := lookupQuery.Items[0]["UUID"].(*types.AttributeValueMemberS).Value

	userQuery, err := nosqlur.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk and begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PrefixUser, id)},
			":sk": &types.AttributeValueMemberS{Value: constants.PrefixProfile},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return user, errors.New(constants.ErrFetchingUser)
	}

	item := userQuery.Items[0]
	var ok bool
	user, ok = nosqlur.itemToUser(ctx, item)
	if !ok {
		log.Println("invalid user data in GetUserByEmail")
		return models.User{}, errors.New(constants.ErrFetchingUser)
	}

	return user, nil
}

func (nosqlur *NOSQLUserRepository) GetUserById(ctx context.Context, id string) (models.User, error) {
	var user models.User

	userQuery, err := nosqlur.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk and begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PrefixUser, id)},
			":sk": &types.AttributeValueMemberS{Value: constants.PrefixProfile},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return user, errors.New(constants.ErrFetchingUser)
	}

	if len(userQuery.Items) == 0 {
		log.Println("user not found in GetUserById")
		return user, errors.New(constants.ErrUserNotFound)
	}

	item := userQuery.Items[0]
	var ok bool
	user, ok = nosqlur.itemToUser(ctx, item)
	if !ok {
		log.Println("invalid user data in GetUserById")
		return models.User{}, errors.New(constants.ErrFetchingUser)
	}

	return user, nil
}

func (nosqlur *NOSQLUserRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	// 1. Get all UUIDs from reverse lookup
	lookupQuery, err := nosqlur.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: constants.PKUser},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New(constants.ErrFetchingUserList)
	}

	var userKeys []map[string]types.AttributeValue
	for _, item := range lookupQuery.Items {
		if val, ok := item["UUID"]; ok {
			uid := val.(*types.AttributeValueMemberS).Value
			userKeys = append(userKeys, map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PrefixUser, uid)},
				"SK": &types.AttributeValueMemberS{Value: constants.PrefixProfile},
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
	var validCount, skippedCount int

	for _, item := range userItems {
		wg.Add(1)
		go func(itm map[string]types.AttributeValue) {
			defer wg.Done()
			user, ok := nosqlur.itemToUser(ctx, itm)
			// filter inactive users or records with invalid data
			if ok && user.IsActive {
				mu.Lock()
				users = append(users, user)
				validCount++
				mu.Unlock()
			} else {
				mu.Lock()
				skippedCount++
				mu.Unlock()
			}
		}(item)
	}

	wg.Wait()
	log.Printf("userrepository: fetched %d user profiles, kept %d active/valid, skipped %d", len(userItems), validCount, skippedCount)
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
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PrefixUser, user.UserID.String())},
			"SK": &types.AttributeValueMemberS{Value: constants.PrefixProfile},
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionValues,
		ExpressionAttributeNames:  expressionNames,
	})
	if err != nil {
		log.Println(err.Error())
		return errors.New(constants.ErrSavingUser)
	}

	return nil
}

// DONE: get office id instead of name
func (nosqlur *NOSQLUserRepository) CreateUser(ctx context.Context, name, email, password, officeId string, role roles.Role) error {
	if officeId == "" {
		log.Println("no office id provided")
		return errors.New(constants.ErrOfficeIDRequired)
	}

	// check if office exists
	officeQuery, err := nosqlur.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk AND SK = :sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: constants.PKOffice},
			":sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PrefixDetails, officeId)},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return errors.New(constants.ErrFetchingOffice)
	}
	if len(officeQuery.Items) == 0 {
		log.Printf("no office with id %s found", officeId)
		return errors.New(constants.ErrSpecifiedOfficeNotFound)
	}

	// check if user exists
	userQuery, err := nosqlur.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(config.DynamoDBTable),
		KeyConditionExpression: aws.String("PK = :pk AND SK = :sk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: constants.PKUser},
			":sk": &types.AttributeValueMemberS{Value: email},
		},
	})
	if err != nil {
		log.Println(err.Error())
		return errors.New(constants.ErrCheckingExistingUser)
	}
	if len(userQuery.Items) > 0 {
		log.Printf("user with email %s already exists", email)
		return errors.New(constants.ErrUserAlreadyExists)
	}

	userID := uuid.New().String()

	item := map[string]types.AttributeValue{
		"PK":           &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PrefixUser, userID)},
		"SK":           &types.AttributeValueMemberS{Value: constants.PrefixProfile},
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
						"PK":   &types.AttributeValueMemberS{Value: constants.PKUser},
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
		return errors.New(constants.ErrCreatingUser)
	}

	return nil
}

// helper for dynamo to user; returns false when data is malformed (e.g., bad UUIDs)
func (nosqlur *NOSQLUserRepository) itemToUser(ctx context.Context, item map[string]types.AttributeValue) (models.User, bool) {
	var user models.User
	valid := true

	if idAttr, ok := item["Id"].(*types.AttributeValueMemberS); ok && idAttr != nil {
		if userUUID, err := uuid.Parse(idAttr.Value); err == nil {
			user.UserID = userUUID
		} else {
			log.Printf("invalid user Id uuid: %v", err)
			valid = false
		}
	} else {
		log.Printf("missing user Id attribute in item")
		valid = false
	}

	if usernameAttr, ok := item["Username"].(*types.AttributeValueMemberS); ok && usernameAttr != nil {
		user.Name = usernameAttr.Value
	}

	if emailAttr, ok := item["Email"].(*types.AttributeValueMemberS); ok && emailAttr != nil {
		user.Email = emailAttr.Value
	}

	if pwAttr, ok := item["PasswordHash"].(*types.AttributeValueMemberS); ok && pwAttr != nil {
		user.Password = pwAttr.Value
	}

	if roleAttr, ok := item["Role"].(*types.AttributeValueMemberS); ok && roleAttr != nil {
		roleStr := roleAttr.Value
		switch roleStr {
		case roles.Admin.String():
			user.Role = roles.Admin
		case roles.Customer.String():
			user.Role = roles.Customer
		default:
			user.Role = roles.Customer
		}
	}

	// handle missing IsActive gracefully for legacy rows; default to true
	if isActiveAttr, ok := item["IsActive"].(*types.AttributeValueMemberBOOL); ok && isActiveAttr != nil {
		user.IsActive = isActiveAttr.Value
	} else {
		user.IsActive = true
	}

	if officeAttr, ok := item["OfficeId"].(*types.AttributeValueMemberS); ok && officeAttr != nil {
		if officeUUID, err := uuid.Parse(officeAttr.Value); err == nil {
			user.OfficeID = officeUUID
			user.Office.OfficeID = user.OfficeID

			res, err := nosqlur.client.Query(ctx, &dynamodb.QueryInput{
				TableName:              aws.String(config.DynamoDBTable),
				KeyConditionExpression: aws.String("PK = :pk AND SK = :sk"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":pk": &types.AttributeValueMemberS{Value: constants.PKOffice},
					":sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", constants.PrefixDetails, user.OfficeID.String())},
				},
				AttributesToGet: []string{
					"OfficeName",
				},
			})
			if err == nil && len(res.Items) > 0 {
				user.Office.OfficeName = res.Items[0]["OfficeName"].(*types.AttributeValueMemberS).Value
			}
		} else {
			log.Printf("invalid OfficeId uuid for user %s: %v", user.UserID.String(), err)
			valid = false
		}
	}

	return user, valid
}
