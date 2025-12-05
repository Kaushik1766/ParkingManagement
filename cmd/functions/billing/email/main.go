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

// var billingService billingserviceBillingMgr
var emailService emailservice.EmailManager

func init() {
	emailService = emailservice.NewEmailService()
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, event events.SQSEvent) (events.SQSEventResponse, error) {
	failedEvents := []events.SQSBatchItemFailure{}

	log.Printf("num of messages received: %d", len(event.Records))

	for _, r := range event.Records {
		log.Printf("Body: %s, MessageId: %s\n", r.Body, r.MessageId)

		var sqsMessage models.SQSEmailMessage

		err := json.Unmarshal([]byte(r.Body), &sqsMessage)
		if err != nil {
			log.Println(err)
			failedEvents = append(failedEvents, events.SQSBatchItemFailure{
				ItemIdentifier: r.MessageId,
			})
		}

		err = emailService.SendEmail(ctx, sqsMessage)
		if err != nil {
			log.Println(err)
			failedEvents = append(failedEvents, events.SQSBatchItemFailure{
				ItemIdentifier: r.MessageId,
			})
		}
	}

	return events.SQSEventResponse{
		BatchItemFailures: failedEvents,
	}, nil
}