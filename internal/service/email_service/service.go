package emailservice

import (
	"context"
	"log"
	"net/smtp"
	"os"
	"strconv"
)

type EmailService struct {
}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (emailService *EmailService) SendEmail(ctx context.Context, to, header, body string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	err := smtp.SendMail(smtpHost+":"+strconv.Itoa(smtpPort), auth, smtpUsername, []string{to}, []byte("To: "+to+"\r\n"+"Subject: "+header+"\r\n"+"\r\n"+body))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
