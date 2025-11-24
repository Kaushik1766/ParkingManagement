package mocks

import (
	context "context"
	reflect "reflect"

	models "github.com/Kaushik1766/ParkingManagement/internal/models"
	gomock "go.uber.org/mock/gomock"
)

// MockBillStorage is a mock of BillStorage interface.
type MockBillStorage struct {
	ctrl     *gomock.Controller
	recorder *MockBillStorageMockRecorder
}

// MockBillStorageMockRecorder is the mock recorder for MockBillStorage.
type MockBillStorageMockRecorder struct {
	mock *MockBillStorage
}

// NewMockBillStorage creates a new mock instance.
func NewMockBillStorage(ctrl *gomock.Controller) *MockBillStorage {
	mock := &MockBillStorage{ctrl: ctrl}
	mock.recorder = &MockBillStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBillStorage) EXPECT() *MockBillStorageMockRecorder {
	return m.recorder
}

// SaveBill mocks base method.
func (m *MockBillStorage) SaveBill(ctx context.Context, bill models.BillDTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveBill", ctx, bill)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveBill indicates an expected call of SaveBill.
func (mr *MockBillStorageMockRecorder) SaveBill(ctx, bill interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveBill", reflect.TypeOf((*MockBillStorage)(nil).SaveBill), ctx, bill)
}

// GetBill mocks base method.
func (m *MockBillStorage) GetBill(ctx context.Context, userId string, month, year int) (models.BillDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBill", ctx, userId, month, year)
	ret0, _ := ret[0].(models.BillDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBill indicates an expected call of GetBill.
func (mr *MockBillStorageMockRecorder) GetBill(ctx, userId, month, year interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBill", reflect.TypeOf((*MockBillStorage)(nil).GetBill), ctx, userId, month, year)
}
