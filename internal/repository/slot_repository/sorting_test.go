package slotrepository

import (
	"testing"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
)

func TestSortSlots(t *testing.T) {
	tests := []struct {
		name  string
		slots []models.Slot
		want  []int
	}{
		{
			name: "unsorted slots",
			slots: []models.Slot{
				{SlotNumber: 3},
				{SlotNumber: 1},
				{SlotNumber: 2},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "already sorted",
			slots: []models.Slot{
				{SlotNumber: 1},
				{SlotNumber: 2},
				{SlotNumber: 3},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "empty",
			slots: []models.Slot{},
			want: []int{},
		},
		{
			name: "single element",
			slots: []models.Slot{
				{SlotNumber: 1},
			},
			want: []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortSlots(tt.slots)
			if len(tt.slots) != len(tt.want) {
				t.Errorf("got len %d, want %d", len(tt.slots), len(tt.want))
			}
			for i, slot := range tt.slots {
				if slot.SlotNumber != tt.want[i] {
					t.Errorf("index %d: got %d, want %d", i, slot.SlotNumber, tt.want[i])
				}
			}
		})
	}
}
