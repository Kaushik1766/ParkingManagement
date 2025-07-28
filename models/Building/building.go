package building

import "github.com/google/uuid"

type Building struct {
	BuildingId   uuid.UUID
	BuildingName string
}
