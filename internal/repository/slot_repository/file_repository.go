package slotrepository

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/Kaushik1766/ParkingManagement/internal/models/slot"
)

type FileSlotRepository struct {
	*sync.Mutex
	slots []slot.Slot
}

func NewFileSlotRepository() *FileSlotRepository {
	data, err := os.ReadFile("slots.json")
	if err != nil {
		os.WriteFile("slots.json", []byte("[]"), 0666)
		data, err = json.Marshal([]slot.Slot{})
		if err != nil {
			fmt.Println("unable to marshal")
		}
	}

	var slotData []slot.Slot
	err = json.Unmarshal(data, &slotData)
	if err != nil {
		fmt.Println(err)
		panic("corrupted data")
	}
	return &FileSlotRepository{
		Mutex: &sync.Mutex{},
		slots: slotData,
	}
}
