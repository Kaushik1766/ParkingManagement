package floor

import (
	"github.com/google/uuid"
)

type Floor struct {
	FloorNumber int
	BuildingId  uuid.UUID
}
