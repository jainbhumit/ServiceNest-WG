package controllers

import (
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"serviceNest/controllers"
	"serviceNest/model"
	"serviceNest/tests/mocks"
	"testing"
)

func TestAdminController_ViewAllService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockAdminService(ctrl)

	// Create mock response data
	mockServices := []model.Service{
		{ID: "1", Name: "Service 1", Description: "Description 1"},
		{ID: "2", Name: "Service 2", Description: "Description 2"},
	}

	// Set up expectations for the mock service
	mockService.EXPECT().GetAllService().Return(mockServices, nil)

	// Create AdminController with mock service
	adminController := controllers.NewAdminController(mockService)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/admin/services", nil)
	assert.NoError(t, err)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler function
	handler := http.HandlerFunc(adminController.ViewAllService)
	handler.ServeHTTP(rr, req)

	// Assert that the status code is 200 OK
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert the body contains the expected services
	expected := `{
        "data": [
            {
                "id": "1",
                "name": "Service 1",
                "description": "Description 1",
                "category": "",
                "price": 0,
                "provider_id": ""
            },
            {
                "id": "2",
                "name": "Service 2",
                "description": "Description 2",
                "category": "",
                "price": 0,
                "provider_id": ""
            }
        ],
        "message": "All available services",
        "status": "Success"
    }`

	// Assert the body contains the expected services
	assert.JSONEq(t, expected, rr.Body.String())
}
func TestAdminController_ViewAllService_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockAdminService(ctrl)

	// Set up expectations for the mock service
	mockService.EXPECT().GetAllService().Return(nil, errors.New("error getting services"))

	// Create AdminController with mock service
	adminController := controllers.NewAdminController(mockService)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/admin/services", nil)
	assert.NoError(t, err)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler function
	handler := http.HandlerFunc(adminController.ViewAllService)
	handler.ServeHTTP(rr, req)

	// Assert that the status code is 200 OK
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Assert the body contains the expected services
	expected := `{
        "message": "error fetching services",
        "status": "Fail"
    }`

	// Assert the body contains the expected services
	assert.JSONEq(t, expected, rr.Body.String())
}

func TestAdminController_DeleteService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockAdminService(ctrl)

	// Set up expectations for the mock service
	mockService.EXPECT().DeleteService("1").Return(nil)

	// Create AdminController with mock service
	adminController := controllers.NewAdminController(mockService)

	// Create a new HTTP request with service ID
	req, err := http.NewRequest("DELETE", "/admin/service/1", nil)
	assert.NoError(t, err)

	// Set up mux router to handle route variables
	req = mux.SetURLVars(req, map[string]string{"serviceID": "1"})

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler function
	handler := http.HandlerFunc(adminController.DeleteService)
	handler.ServeHTTP(rr, req)

	// Assert that the status code is 200 OK
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert the body contains success message
	expected := `{"message":"Service deleted successfully","status":"Success"}`
	assert.JSONEq(t, expected, rr.Body.String())
}

func TestAdminController_ViewReports(t *testing.T) {
	// Initialize mock admin service
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockAdminService(ctrl)

	// Create mock response data
	mockReports := []model.ServiceRequest{
		{ID: "1", ServiceName: "Service 1", Status: "Pending"},
		{ID: "2", ServiceName: "Service 2", Status: "Cancelled"},
	}

	// Set up expectations for the mock service
	mockService.EXPECT().ViewReports().Return(mockReports, nil)

	// Create AdminController with mock service
	adminController := controllers.NewAdminController(mockService)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/admin/reports", nil)
	assert.NoError(t, err)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler function
	handler := http.HandlerFunc(adminController.ViewReports)
	handler.ServeHTTP(rr, req)

	// Assert that the status code is 200 OK
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert the body contains the expected reports
	// Parse the actual response body
	var actualResponse map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)

	// Focus on asserting important fields
	data := actualResponse["data"].([]interface{})

	// Assert that the first report matches the expected values
	firstReport := data[0].(map[string]interface{})
	assert.Equal(t, "1", firstReport["id"])
	assert.Equal(t, "Service 1", firstReport["service_name"])
	assert.Equal(t, "Pending", firstReport["status"])

	// Assert that the second report matches the expected values
	secondReport := data[1].(map[string]interface{})
	assert.Equal(t, "2", secondReport["id"])
	assert.Equal(t, "Service 2", secondReport["service_name"])
	assert.Equal(t, "Cancelled", secondReport["status"])
}

func TestAdminController_DeactivateUserAccount(t *testing.T) {
	// Initialize mock admin service
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockAdminService(ctrl)

	// Set up expectations for the mock service
	mockService.EXPECT().DeactivateAccount("1").Return(nil)

	// Create AdminController with mock service
	adminController := controllers.NewAdminController(mockService)

	// Create a new HTTP request with provider ID
	req, err := http.NewRequest("PUT", "/admin/provider/1/deactivate", nil)
	assert.NoError(t, err)

	// Set up mux router to handle route variables
	req = mux.SetURLVars(req, map[string]string{"providerID": "1"})

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler function
	handler := http.HandlerFunc(adminController.DeactivateUserAccount)
	handler.ServeHTTP(rr, req)

	// Assert that the status code is 200 OK
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert the body contains success message
	expected := `{"message":"Account deactivated successfully","status":"Success"}`
	assert.JSONEq(t, expected, rr.Body.String())
}

func TestAdminController_ViewReports_Error(t *testing.T) {
	// Initialize mock admin service
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockAdminService(ctrl)

	// Set up expectations for the mock service
	mockService.EXPECT().ViewReports().Return(nil, errors.New("Error fetching requests"))

	// Create AdminController with mock service
	adminController := controllers.NewAdminController(mockService)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/admin/reports", nil)
	assert.NoError(t, err)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler function
	handler := http.HandlerFunc(adminController.ViewReports)
	handler.ServeHTTP(rr, req)

	// Assert that the status code is 200 OK
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

}
func TestAdminController_DeactivateUserAccount_Error(t *testing.T) {
	// Initialize mock admin service
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockAdminService(ctrl)

	// Set up expectations for the mock service
	mockService.EXPECT().DeactivateAccount("1").Return(errors.New("Error deactive user"))

	// Create AdminController with mock service
	adminController := controllers.NewAdminController(mockService)

	// Create a new HTTP request with provider ID
	req, err := http.NewRequest("PUT", "/admin/provider/1/deactivate", nil)
	assert.NoError(t, err)

	// Set up mux router to handle route variables
	req = mux.SetURLVars(req, map[string]string{"providerID": "1"})

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler function
	handler := http.HandlerFunc(adminController.DeactivateUserAccount)
	handler.ServeHTTP(rr, req)

	// Assert that the status code is 200 OK
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Assert the body contains success message
	expected := `{"message":"Error deactivating account","status":"Fail"}`
	assert.JSONEq(t, expected, rr.Body.String())
}
