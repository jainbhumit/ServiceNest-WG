// Code generated by MockGen. DO NOT EDIT.
// Source: C:\Users\bjain\serviceNest\interfaces\service_provider_repository_interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	model "serviceNest/model"

	gomock "github.com/golang/mock/gomock"
)

// MockServiceProviderRepository is a mock of ServiceProviderRepository interface.
type MockServiceProviderRepository struct {
	ctrl     *gomock.Controller
	recorder *MockServiceProviderRepositoryMockRecorder
}

// MockServiceProviderRepositoryMockRecorder is the mock recorder for MockServiceProviderRepository.
type MockServiceProviderRepositoryMockRecorder struct {
	mock *MockServiceProviderRepository
}

// NewMockServiceProviderRepository creates a new mock instance.
func NewMockServiceProviderRepository(ctrl *gomock.Controller) *MockServiceProviderRepository {
	mock := &MockServiceProviderRepository{ctrl: ctrl}
	mock.recorder = &MockServiceProviderRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceProviderRepository) EXPECT() *MockServiceProviderRepositoryMockRecorder {
	return m.recorder
}

// AddReview mocks base method.
func (m *MockServiceProviderRepository) AddReview(providerID, householderID, review string, rating float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddReview", providerID, householderID, review, rating)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddReview indicates an expected call of AddReview.
func (mr *MockServiceProviderRepositoryMockRecorder) AddReview(providerID, householderID, review, rating interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddReview", reflect.TypeOf((*MockServiceProviderRepository)(nil).AddReview), providerID, householderID, review, rating)
}

// GetProviderByID mocks base method.
func (m *MockServiceProviderRepository) GetProviderByID(providerID string) (*model.ServiceProvider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProviderByID", providerID)
	ret0, _ := ret[0].(*model.ServiceProvider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProviderByID indicates an expected call of GetProviderByID.
func (mr *MockServiceProviderRepositoryMockRecorder) GetProviderByID(providerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProviderByID", reflect.TypeOf((*MockServiceProviderRepository)(nil).GetProviderByID), providerID)
}

// GetProviderByServiceID mocks base method.
func (m *MockServiceProviderRepository) GetProviderByServiceID(serviceID string) (*model.ServiceProvider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProviderByServiceID", serviceID)
	ret0, _ := ret[0].(*model.ServiceProvider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProviderByServiceID indicates an expected call of GetProviderByServiceID.
func (mr *MockServiceProviderRepositoryMockRecorder) GetProviderByServiceID(serviceID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProviderByServiceID", reflect.TypeOf((*MockServiceProviderRepository)(nil).GetProviderByServiceID), serviceID)
}

// GetProvidersByServiceType mocks base method.
func (m *MockServiceProviderRepository) GetProvidersByServiceType(serviceType string) ([]model.ServiceProvider, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProvidersByServiceType", serviceType)
	ret0, _ := ret[0].([]model.ServiceProvider)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProvidersByServiceType indicates an expected call of GetProvidersByServiceType.
func (mr *MockServiceProviderRepositoryMockRecorder) GetProvidersByServiceType(serviceType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProvidersByServiceType", reflect.TypeOf((*MockServiceProviderRepository)(nil).GetProvidersByServiceType), serviceType)
}

// SaveServiceProvider mocks base method.
func (m *MockServiceProviderRepository) SaveServiceProvider(provider model.ServiceProvider) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveServiceProvider", provider)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveServiceProvider indicates an expected call of SaveServiceProvider.
func (mr *MockServiceProviderRepositoryMockRecorder) SaveServiceProvider(provider interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveServiceProvider", reflect.TypeOf((*MockServiceProviderRepository)(nil).SaveServiceProvider), provider)
}

// UpdateServiceProvider mocks base method.
func (m *MockServiceProviderRepository) UpdateServiceProvider(provider *model.ServiceProvider) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateServiceProvider", provider)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateServiceProvider indicates an expected call of UpdateServiceProvider.
func (mr *MockServiceProviderRepositoryMockRecorder) UpdateServiceProvider(provider interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateServiceProvider", reflect.TypeOf((*MockServiceProviderRepository)(nil).UpdateServiceProvider), provider)
}
