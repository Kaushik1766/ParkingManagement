package emailservice

import (
	"context"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
)

type EmailManager interface {
	SendEmail(ctx context.Context, emailMessage models.SQSEmailMessage) error
}
