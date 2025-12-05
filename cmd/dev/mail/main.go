package main

import (
	"context"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
	emailservice "github.com/Kaushik1766/ParkingManagement/internal/service/email_service"
)

func main() {
	emailService := emailservice.NewEmailService()

	emailService.SendEmail(context.Background(), models.SQSEmailMessage{
		To:     "noobitanobi176@gmail.com",
		Header: "hello world",
		Body:   "hello from go",
	})
}
