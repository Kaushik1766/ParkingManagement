package building

import "github.com/google/uuid"

type Building struct {
	BuildingId   uuid.UUID
	BuildingName string
}

func (b Building) GetID() string {
	return b.BuildingId.String()
}
