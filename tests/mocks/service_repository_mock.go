// Code generated by MockGen. DO NOT EDIT.
// Source: C:\Users\bjain\serviceNest\interfaces\service_repository_interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	model "serviceNest/model"

	gomock "github.com/golang/mock/gomock"
)

// MockServiceRepository is a mock of ServiceRepository interface.
type MockServiceRepository struct {
	ctrl     *gomock.Controller
	recorder *MockServiceRepositoryMockRecorder
}

// MockServiceRepositoryMockRecorder is the mock recorder for MockServiceRepository.
type MockServiceRepositoryMockRecorder struct {
	mock *MockServiceRepository
}

// NewMockServiceRepository creates a new mock instance.
func NewMockServiceRepository(ctrl *gomock.Controller) *MockServiceRepository {
	mock := &MockServiceRepository{ctrl: ctrl}
	mock.recorder = &MockServiceRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceRepository) EXPECT() *MockServiceRepositoryMockRecorder {
	return m.recorder
}

// GetAllServices mocks base method.
func (m *MockServiceRepository) GetAllServices() ([]model.Service, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllServices")
	ret0, _ := ret[0].([]model.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllServices indicates an expected call of GetAllServices.
func (mr *MockServiceRepositoryMockRecorder) GetAllServices() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllServices", reflect.TypeOf((*MockServiceRepository)(nil).GetAllServices))
}

// GetServiceByID mocks base method.
func (m *MockServiceRepository) GetServiceByID(serviceID string) (*model.Service, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetServiceByID", serviceID)
	ret0, _ := ret[0].(*model.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetServiceByID indicates an expected call of GetServiceByID.
func (mr *MockServiceRepositoryMockRecorder) GetServiceByID(serviceID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServiceByID", reflect.TypeOf((*MockServiceRepository)(nil).GetServiceByID), serviceID)
}

// GetServiceByName mocks base method.
func (m *MockServiceRepository) GetServiceByName(serviceName string) (*model.Service, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetServiceByName", serviceName)
	ret0, _ := ret[0].(*model.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetServiceByName indicates an expected call of GetServiceByName.
func (mr *MockServiceRepositoryMockRecorder) GetServiceByName(serviceName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServiceByName", reflect.TypeOf((*MockServiceRepository)(nil).GetServiceByName), serviceName)
}

// GetServiceByProviderID mocks base method.
func (m *MockServiceRepository) GetServiceByProviderID(providerID string) ([]model.Service, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetServiceByProviderID", providerID)
	ret0, _ := ret[0].([]model.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetServiceByProviderID indicates an expected call of GetServiceByProviderID.
func (mr *MockServiceRepositoryMockRecorder) GetServiceByProviderID(providerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServiceByProviderID", reflect.TypeOf((*MockServiceRepository)(nil).GetServiceByProviderID), providerID)
}

// RemoveService mocks base method.
func (m *MockServiceRepository) RemoveService(serviceID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveService", serviceID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveService indicates an expected call of RemoveService.
func (mr *MockServiceRepositoryMockRecorder) RemoveService(serviceID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveService", reflect.TypeOf((*MockServiceRepository)(nil).RemoveService), serviceID)
}

// RemoveServiceByProviderID mocks base method.
func (m *MockServiceRepository) RemoveServiceByProviderID(providerID, serviceID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveServiceByProviderID", providerID, serviceID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveServiceByProviderID indicates an expected call of RemoveServiceByProviderID.
func (mr *MockServiceRepositoryMockRecorder) RemoveServiceByProviderID(providerID, serviceID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveServiceByProviderID", reflect.TypeOf((*MockServiceRepository)(nil).RemoveServiceByProviderID), providerID, serviceID)
}

// SaveAllServices mocks base method.
func (m *MockServiceRepository) SaveAllServices(services []model.Service) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveAllServices", services)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveAllServices indicates an expected call of SaveAllServices.
func (mr *MockServiceRepositoryMockRecorder) SaveAllServices(services interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveAllServices", reflect.TypeOf((*MockServiceRepository)(nil).SaveAllServices), services)
}

// SaveService mocks base method.
func (m *MockServiceRepository) SaveService(service model.Service) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveService", service)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveService indicates an expected call of SaveService.
func (mr *MockServiceRepositoryMockRecorder) SaveService(service interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveService", reflect.TypeOf((*MockServiceRepository)(nil).SaveService), service)
}

// UpdateService mocks base method.
func (m *MockServiceRepository) UpdateService(providerID string, updatedService model.Service) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateService", providerID, updatedService)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateService indicates an expected call of UpdateService.
func (mr *MockServiceRepositoryMockRecorder) UpdateService(providerID, updatedService interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateService", reflect.TypeOf((*MockServiceRepository)(nil).UpdateService), providerID, updatedService)
}
