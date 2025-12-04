package officeservice

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/Kaushik1766/ParkingManagement/internal/models"
	officerepository "github.com/Kaushik1766/ParkingManagement/internal/repository/office_repository"
	"github.com/Kaushik1766/ParkingManagement/mocks"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestNewOfficeService(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOfficeRepo := mocks.NewMockOfficeStorage(ctrl)

	type args struct {
		officeRepo officerepository.OfficeStorage
	}
	tests := []struct {
		name string
		args args
		want *OfficeService
	}{
		{
			name: "valide repo",
			args: args{
				mockOfficeRepo,
			},
			want: &OfficeService{
				officeRepo: mockOfficeRepo,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOfficeService(tt.args.officeRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOfficeService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOfficeService_AddOffice(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOfficeRepo := mocks.NewMockOfficeStorage(ctrl)
	mockOfficeRepo.EXPECT().AddOffice(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	type fields struct {
		officeRepo officerepository.OfficeStorage
	}
	type args struct {
		ctx         context.Context
		officeName  string
		buildingId  string
		floorNumber int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid inputs",
			fields: fields{
				officeRepo: mockOfficeRepo,
			},
			args: args{
				ctx:         context.Background(),
				officeName:  "Test",
				buildingId:  "Test",
				floorNumber: 1,
			},
			wantErr: false,
		},
		{
			name: "invalid inputs",
			fields: fields{
				officeRepo: mockOfficeRepo,
			},
			args: args{
				ctx:         context.Background(),
				officeName:  "Test",
				buildingId:  "Test",
				floorNumber: -1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			officeServ := &OfficeService{
				officeRepo: tt.fields.officeRepo,
			}
			if err := officeServ.AddOffice(tt.args.ctx, tt.args.officeName, tt.args.buildingId, tt.args.floorNumber); (err != nil) != tt.wantErr {
				t.Errorf("AddOffice() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOfficeService_GetAllOfficeNames(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOfficeRepo := mocks.NewMockOfficeStorage(ctrl)
	mockOfficeRepo.EXPECT().GetAllOffices(gomock.Any()).Return([]models.Office{
		{
			OfficeID:    uuid.Nil,
			OfficeName:  "wg",
			BuildingID:  uuid.Nil,
			FloorNumber: 1,
		},
	}, nil)
	mockOfficeRepo.EXPECT().GetAllOffices(gomock.Any()).Return(nil, errors.New("repo error"))

	type fields struct {
		officeRepo officerepository.OfficeStorage
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "office names found",
			fields: fields{
				officeRepo: mockOfficeRepo,
			},
			args: args{
				ctx: context.Background(),
			},
			want:    []string{"wg"},
			wantErr: false,
		},
		{
			name: "no office name found",
			fields: fields{
				officeRepo: mockOfficeRepo,
			},
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			officeServ := &OfficeService{
				officeRepo: tt.fields.officeRepo,
			}
			got, err := officeServ.GetAllOfficeNames(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllOfficeNames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllOfficeNames() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOfficeService_ListOfficesByBuilding(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOfficeRepo := mocks.NewMockOfficeStorage(ctrl)
	mockOfficeRepo.EXPECT().GetOfficesByBuilding(gomock.Any(), gomock.Any()).Return([]models.Office{
		{
			OfficeID:    uuid.Nil,
			OfficeName:  "wg",
			BuildingID:  uuid.Nil,
			FloorNumber: 1,
		},
	}, nil)
	mockOfficeRepo.EXPECT().GetOfficesByBuilding(gomock.Any(), gomock.Any()).Return(nil, errors.New("repo error"))

	type fields struct {
		officeRepo officerepository.OfficeStorage
	}
	type args struct {
		ctx        context.Context
		buildingId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.OfficeDTO
		wantErr bool
	}{
		{
			name: "office names found",
			fields: fields{
				officeRepo: mockOfficeRepo,
			},
			args: args{
				ctx:        context.Background(),
				buildingId: "Test",
			},
			wantErr: false,
			want: []models.OfficeDTO{
				{
					BuildingID:  uuid.Nil.String(),
					FloorNumber: 1,
					OfficeName:  "wg",
					OfficeID:    uuid.Nil.String(),
				},
			},
		},
		{
			name: "office names not found",
			fields: fields{
				officeRepo: mockOfficeRepo,
			},
			args: args{
				ctx:        context.Background(),
				buildingId: "Test",
			},
			wantErr: true,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			officeServ := &OfficeService{
				officeRepo: tt.fields.officeRepo,
			}
			got, err := officeServ.ListOfficesByBuilding(tt.args.ctx, tt.args.buildingId)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListOfficesByBuilding() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListOfficesByBuilding() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOfficeService_RemoveOffice(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOfficeRepo := mocks.NewMockOfficeStorage(ctrl)
	mockOfficeRepo.EXPECT().DeleteOffice(gomock.Any(), gomock.Any()).Return(nil)

	type fields struct {
		officeRepo officerepository.OfficeStorage
	}
	type args struct {
		ctx      context.Context
		officeId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid officeId",
			fields: fields{
				officeRepo: mockOfficeRepo,
			},
			args: args{
				ctx:      context.Background(),
				officeId: "Test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			officeServ := &OfficeService{
				officeRepo: tt.fields.officeRepo,
			}
			if err := officeServ.RemoveOffice(tt.args.ctx, tt.args.officeId); (err != nil) != tt.wantErr {
				t.Errorf("RemoveOffice() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
