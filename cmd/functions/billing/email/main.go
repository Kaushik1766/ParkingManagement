package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
	emailservice "github.com/Kaushik1766/ParkingManagement/internal/service/email_service"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var emailSvc *emailservice.EmailService

func init() {
	emailSvc = emailservice.NewEmailService()
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) (events.SQSEventResponse, error) {
	var batchItemFailures []events.SQSBatchItemFailure

	for _, message := range sqsEvent.Records {
		var emailMessage models.SQSEmailMessage
		if err := json.Unmarshal([]byte(message.Body), &emailMessage); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			batchItemFailures = append(batchItemFailures, events.SQSBatchItemFailure{
				ItemIdentifier: message.MessageId,
			})
			continue
		}

		if err := emailSvc.SendEmail(ctx, emailMessage); err != nil {
			log.Printf("Error sending email: %v", err)
			batchItemFailures = append(batchItemFailures, events.SQSBatchItemFailure{
				ItemIdentifier: message.MessageId,
			})
			continue
		}

		log.Printf("Successfully sent email to: %s", emailMessage.To)
	}

	return events.SQSEventResponse{
		BatchItemFailures: batchItemFailures,
	}, nil
}

func main() {
	lambda.Start(handler)
}
