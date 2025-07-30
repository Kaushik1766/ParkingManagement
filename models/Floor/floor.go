package floor

import (
	"fmt"

	"github.com/google/uuid"
)

type Floor struct {
	BuildingId  uuid.UUID
	FloorNumber int
}

func (f Floor) GetID() string {
	return fmt.Sprintf("%v%v", f.BuildingId, f.FloorNumber)
}
