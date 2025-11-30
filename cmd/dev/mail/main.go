package main

import (
	"context"

	emailservice "github.com/Kaushik1766/ParkingManagement/internal/service/email_service"
)

func main() {
	emailService := emailservice.NewEmailService()

	emailService.SendEmail(context.Background(), "noobitanobi176@gmail.com", "hello", "hello world")
}
