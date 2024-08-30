package repository_test

//import (
//	"serviceNest/model"
//	"serviceNest/repository_test"
//	"testing"
//
//	"github.com/stretchr/testify/assert"
//)
//
//func TestServiceRequestRepository_SaveServiceRequest(t *testing.T) {
//	repo := repository_test.NewServiceRequestRepository("test_service_requests.json")
//
//	request := model.ServiceRequest{
//		ID:            "request1",
//		ServiceID:     "service1",
//		HouseholderID: "householder1",
//		Status:        "pending",
//	}
//
//	err := repo.SaveServiceRequest(request)
//	assert.NoError(t, err, "Expected no error when saving service_test request")
//
//	savedRequest, err := repo.GetServiceRequestByID(request.ID)
//	assert.NoError(t, err, "Expected no error when fetching service_test request")
//	assert.Equal(t, request.Status, savedRequest.Status, "Expected service_test request status to match")
//}
//
//func TestServiceRequestRepository_UpdateServiceRequest(t *testing.T) {
//	repo := repository_test.NewServiceRequestRepository("test_service_requests.json")
//
//	request := model.ServiceRequest{
//		ID:     "request1",
//		Status: "accepted",
//	}
//
//	err := repo.UpdateServiceRequest(request)
//	assert.NoError(t, err, "Expected no error when updating service_test request")
//
//	updatedRequest, err := repo.GetServiceRequestByID(request.ID)
//	assert.NoError(t, err, "Expected no error when fetching service_test request")
//	assert.Equal(t, request.Status, updatedRequest.Status, "Expected service_test request status to be updated")
//}
