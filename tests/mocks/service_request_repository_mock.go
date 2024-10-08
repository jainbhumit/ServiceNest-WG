// Code generated by MockGen. DO NOT EDIT.
// Source: C:\Users\bjain\serviceNest\interfaces\service_request_repository_interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	model "serviceNest/model"

	gomock "github.com/golang/mock/gomock"
)

// MockServiceRequestRepository is a mock of ServiceRequestRepository interface.
type MockServiceRequestRepository struct {
	ctrl     *gomock.Controller
	recorder *MockServiceRequestRepositoryMockRecorder
}

// MockServiceRequestRepositoryMockRecorder is the mock recorder for MockServiceRequestRepository.
type MockServiceRequestRepositoryMockRecorder struct {
	mock *MockServiceRequestRepository
}

// NewMockServiceRequestRepository creates a new mock instance.
func NewMockServiceRequestRepository(ctrl *gomock.Controller) *MockServiceRequestRepository {
	mock := &MockServiceRequestRepository{ctrl: ctrl}
	mock.recorder = &MockServiceRequestRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceRequestRepository) EXPECT() *MockServiceRequestRepositoryMockRecorder {
	return m.recorder
}

// GetAllServiceRequests mocks base method.
func (m *MockServiceRequestRepository) GetAllServiceRequests() ([]model.ServiceRequest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllServiceRequests")
	ret0, _ := ret[0].([]model.ServiceRequest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllServiceRequests indicates an expected call of GetAllServiceRequests.
func (mr *MockServiceRequestRepositoryMockRecorder) GetAllServiceRequests() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllServiceRequests", reflect.TypeOf((*MockServiceRequestRepository)(nil).GetAllServiceRequests))
}

// GetServiceProviderByRequestID mocks base method.
func (m *MockServiceRequestRepository) GetServiceProviderByRequestID(requestID, providerID string) (*model.ServiceRequest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetServiceProviderByRequestID", requestID, providerID)
	ret0, _ := ret[0].(*model.ServiceRequest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetServiceProviderByRequestID indicates an expected call of GetServiceProviderByRequestID.
func (mr *MockServiceRequestRepositoryMockRecorder) GetServiceProviderByRequestID(requestID, providerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServiceProviderByRequestID", reflect.TypeOf((*MockServiceRequestRepository)(nil).GetServiceProviderByRequestID), requestID, providerID)
}

// GetServiceRequestByID mocks base method.
func (m *MockServiceRequestRepository) GetServiceRequestByID(requestID string) (*model.ServiceRequest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetServiceRequestByID", requestID)
	ret0, _ := ret[0].(*model.ServiceRequest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetServiceRequestByID indicates an expected call of GetServiceRequestByID.
func (mr *MockServiceRequestRepositoryMockRecorder) GetServiceRequestByID(requestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServiceRequestByID", reflect.TypeOf((*MockServiceRequestRepository)(nil).GetServiceRequestByID), requestID)
}

// GetServiceRequestsByHouseholderID mocks base method.
func (m *MockServiceRequestRepository) GetServiceRequestsByHouseholderID(householderID string) ([]model.ServiceRequest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetServiceRequestsByHouseholderID", householderID)
	ret0, _ := ret[0].([]model.ServiceRequest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetServiceRequestsByHouseholderID indicates an expected call of GetServiceRequestsByHouseholderID.
func (mr *MockServiceRequestRepositoryMockRecorder) GetServiceRequestsByHouseholderID(householderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServiceRequestsByHouseholderID", reflect.TypeOf((*MockServiceRequestRepository)(nil).GetServiceRequestsByHouseholderID), householderID)
}

// GetServiceRequestsByProviderID mocks base method.
func (m *MockServiceRequestRepository) GetServiceRequestsByProviderID(providerID string) ([]model.ServiceRequest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetServiceRequestsByProviderID", providerID)
	ret0, _ := ret[0].([]model.ServiceRequest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetServiceRequestsByProviderID indicates an expected call of GetServiceRequestsByProviderID.
func (mr *MockServiceRequestRepositoryMockRecorder) GetServiceRequestsByProviderID(providerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServiceRequestsByProviderID", reflect.TypeOf((*MockServiceRequestRepository)(nil).GetServiceRequestsByProviderID), providerID)
}

// SaveAllServiceRequests mocks base method.
func (m *MockServiceRequestRepository) SaveAllServiceRequests(serviceRequests []model.ServiceRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveAllServiceRequests", serviceRequests)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveAllServiceRequests indicates an expected call of SaveAllServiceRequests.
func (mr *MockServiceRequestRepositoryMockRecorder) SaveAllServiceRequests(serviceRequests interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveAllServiceRequests", reflect.TypeOf((*MockServiceRequestRepository)(nil).SaveAllServiceRequests), serviceRequests)
}

// SaveServiceRequest mocks base method.
func (m *MockServiceRequestRepository) SaveServiceRequest(request model.ServiceRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveServiceRequest", request)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveServiceRequest indicates an expected call of SaveServiceRequest.
func (mr *MockServiceRequestRepositoryMockRecorder) SaveServiceRequest(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveServiceRequest", reflect.TypeOf((*MockServiceRequestRepository)(nil).SaveServiceRequest), request)
}

// UpdateServiceRequest mocks base method.
func (m *MockServiceRequestRepository) UpdateServiceRequest(updatedRequest *model.ServiceRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateServiceRequest", updatedRequest)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateServiceRequest indicates an expected call of UpdateServiceRequest.
func (mr *MockServiceRequestRepositoryMockRecorder) UpdateServiceRequest(updatedRequest interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateServiceRequest", reflect.TypeOf((*MockServiceRequestRepository)(nil).UpdateServiceRequest), updatedRequest)
}
