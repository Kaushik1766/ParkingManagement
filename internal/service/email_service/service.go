package emailservice

import (
	"context"
	"log"
	"net/smtp"
	"os"
	"strconv"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
)

type EmailService struct {
}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (emailService *EmailService) SendEmail(ctx context.Context, emailMessage models.SQSEmailMessage) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	err := smtp.SendMail(smtpHost+":"+strconv.Itoa(smtpPort), auth, smtpUsername, []string{emailMessage.To}, []byte("To: "+emailMessage.To+"\r\n"+"Subject: "+emailMessage.Header+"\r\n"+"\r\n"+emailMessage.Body))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
