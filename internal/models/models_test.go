package models

import (
	"testing"
	"time"

	vehicletypes "github.com/Kaushik1766/ParkingManagement/internal/models/enums/vehicle_types"
	"github.com/google/uuid"
)

func TestBillDTO_String(t *testing.T) {
	type fields struct {
		ParkingHistory []ParkingHistoryDTO
		TotalAmount    float64
		BillDate       string
		UserId         string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "bill with single parking history",
			fields: fields{
				ParkingHistory: []ParkingHistoryDTO{
					{
						TicketId:     "ticket123",
						NumberPlate:  "ka01hk1234",
						BuildingId:   "bldg-001",
						FLoorNumber:  1,
						SlotNumber:   10,
						StartTime:    time.Date(2025, 8, 28, 9, 0, 0, 0, time.UTC),
						EndTime:      time.Date(2025, 8, 28, 18, 0, 0, 0, time.UTC),
						VechicleType: vehicletypes.FourWheeler,
					},
				},
				TotalAmount: 150.50,
				BillDate:    "2025-08-28",
				UserId:      "kaushik@a.com",
			},
			want: "\n\nParkingHistory:\nTicketId: ticket123\nNumberPlate: ka01hk1234\nBuildingId: bldg-001\nFloorNumber: 1\nSlotNumber: 10\nStartTime: 2025-08-28 09:00:00 +0000 UTC\nEndTime: 2025-08-28 18:00:00 +0000 UTC\n\n\n TotalAmount: 150.50\n BillDate: 2025-08-28\n UserId: kaushik@a.com\n\n",
		},
		{
			name: "admin bill with multiple parking entries",
			fields: fields{
				ParkingHistory: []ParkingHistoryDTO{
					{
						TicketId:     "ticket456",
						NumberPlate:  "ka02ad5678",
						BuildingId:   "bldg-002",
						FLoorNumber:  2,
						SlotNumber:   15,
						StartTime:    time.Date(2025, 8, 28, 8, 30, 0, 0, time.UTC),
						EndTime:      time.Date(2025, 8, 28, 17, 30, 0, 0, time.UTC),
						VechicleType: vehicletypes.TwoWheeler,
					},
					{
						TicketId:     "ticket789",
						NumberPlate:  "ka03ad9012",
						BuildingId:   "bldg-003",
						FLoorNumber:  3,
						SlotNumber:   20,
						StartTime:    time.Date(2025, 8, 28, 10, 0, 0, 0, time.UTC),
						EndTime:      time.Date(2025, 8, 28, 19, 0, 0, 0, time.UTC),
						VechicleType: vehicletypes.FourWheeler,
					},
				},
				TotalAmount: 275.25,
				BillDate:    "2025-08-28",
				UserId:      "admin@a.com",
			},
			want: "\n\nParkingHistory:\nTicketId: ticket456\nNumberPlate: ka02ad5678\nBuildingId: bldg-002\nFloorNumber: 2\nSlotNumber: 15\nStartTime: 2025-08-28 08:30:00 +0000 UTC\nEndTime: 2025-08-28 17:30:00 +0000 UTC\n\nTicketId: ticket789\nNumberPlate: ka03ad9012\nBuildingId: bldg-003\nFloorNumber: 3\nSlotNumber: 20\nStartTime: 2025-08-28 10:00:00 +0000 UTC\nEndTime: 2025-08-28 19:00:00 +0000 UTC\n\n\n TotalAmount: 275.25\n BillDate: 2025-08-28\n UserId: admin@a.com\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bdto := &BillDTO{
				ParkingHistory: tt.fields.ParkingHistory,
				TotalAmount:    tt.fields.TotalAmount,
				BillDate:       tt.fields.BillDate,
				UserId:         tt.fields.UserId,
			}
			if got := bdto.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParkingHistoryDTO_String(t *testing.T) {
	type fields struct {
		TicketId     string
		NumberPlate  string
		BuildingId   string
		FLoorNumber  int
		SlotNumber   int
		StartTime    time.Time
		EndTime      time.Time
		VechicleType vehicletypes.VehicleType
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "kaushik vehicle parking history",
			fields: fields{
				TicketId:     "ticket-kaushik-001",
				NumberPlate:  "ka01hk1234",
				BuildingId:   "building-001",
				FLoorNumber:  1,
				SlotNumber:   5,
				StartTime:    time.Date(2025, 8, 28, 9, 0, 0, 0, time.UTC),
				EndTime:      time.Date(2025, 8, 28, 18, 0, 0, 0, time.UTC),
				VechicleType: vehicletypes.FourWheeler,
			},
			want: "TicketId: ticket-kaushik-001\nNumberPlate: ka01hk1234\nBuildingId: building-001\nFloorNumber: 1\nSlotNumber: 5\nStartTime: 2025-08-28 09:00:00 +0000 UTC\nEndTime: 2025-08-28 18:00:00 +0000 UTC",
		},
		{
			name: "admin vehicle parking history",
			fields: fields{
				TicketId:     "ticket-admin-002",
				NumberPlate:  "ka02ad5678",
				BuildingId:   "building-002",
				FLoorNumber:  2,
				SlotNumber:   10,
				StartTime:    time.Date(2025, 8, 28, 8, 30, 0, 0, time.UTC),
				EndTime:      time.Date(2025, 8, 28, 17, 30, 0, 0, time.UTC),
				VechicleType: vehicletypes.TwoWheeler,
			},
			want: "TicketId: ticket-admin-002\nNumberPlate: ka02ad5678\nBuildingId: building-002\nFloorNumber: 2\nSlotNumber: 10\nStartTime: 2025-08-28 08:30:00 +0000 UTC\nEndTime: 2025-08-28 17:30:00 +0000 UTC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phdto := &ParkingHistoryDTO{
				TicketId:     tt.fields.TicketId,
				NumberPlate:  tt.fields.NumberPlate,
				BuildingId:   tt.fields.BuildingId,
				FLoorNumber:  tt.fields.FLoorNumber,
				SlotNumber:   tt.fields.SlotNumber,
				StartTime:    tt.fields.StartTime,
				EndTime:      tt.fields.EndTime,
				VechicleType: tt.fields.VechicleType,
			}
			if got := phdto.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVehicleDTO_String(t *testing.T) {
	type fields struct {
		NumberPlate  string
		VehicleType  string
		AssignedSlot Slot
	}

	buildingID := uuid.New()

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "kaushik vehicle dto",
			fields: fields{
				NumberPlate: "ka01hk1234",
				VehicleType: "FourWheeler",
				AssignedSlot: Slot{
					BuildingID:  buildingID,
					FloorNumber: 1,
					SlotNumber:  5,
					SlotType:    vehicletypes.FourWheeler,
				},
			},
			want: "ka01hk1234 (FourWheeler)",
		},
		{
			name: "admin vehicle dto",
			fields: fields{
				NumberPlate: "ka02ad5678",
				VehicleType: "TwoWheeler",
				AssignedSlot: Slot{
					BuildingID:  buildingID,
					FloorNumber: 2,
					SlotNumber:  10,
					SlotType:    vehicletypes.TwoWheeler,
				},
			},
			want: "ka02ad5678 (TwoWheeler)",
		},
		{
			name: "unassigned vehicle dto",
			fields: fields{
				NumberPlate: "ka03uk9876",
				VehicleType: "FourWheeler",
				AssignedSlot: Slot{
					BuildingID:  uuid.Nil,
					FloorNumber: 0,
					SlotNumber:  0,
					SlotType:    vehicletypes.FourWheeler,
				},
			},
			want: "ka03uk9876 (FourWheeler)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := VehicleDTO{
				NumberPlate:  tt.fields.NumberPlate,
				VehicleType:  tt.fields.VehicleType,
				AssignedSlot: tt.fields.AssignedSlot,
			}
			if got := v.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
