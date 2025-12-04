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

	msg := []byte("To: " + emailMessage.To + "\r\n" +
		"Subject: " + emailMessage.Header + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		emailMessage.Body)

	err := smtp.SendMail(smtpHost+":"+strconv.Itoa(smtpPort), auth, smtpUsername, []string{emailMessage.To}, msg)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
