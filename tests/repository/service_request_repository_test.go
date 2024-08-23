package repository

import (
	"serviceNest/model"
	"serviceNest/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceRequestRepository_SaveServiceRequest(t *testing.T) {
	repo := repository.NewServiceRequestRepository("test_service_requests.json")

	request := model.ServiceRequest{
		ID:            "request1",
		ServiceID:     "service1",
		HouseholderID: "householder1",
		Status:        "pending",
	}

	err := repo.SaveServiceRequest(request)
	assert.NoError(t, err, "Expected no error when saving service request")

	savedRequest, err := repo.GetServiceRequestByID(request.ID)
	assert.NoError(t, err, "Expected no error when fetching service request")
	assert.Equal(t, request.Status, savedRequest.Status, "Expected service request status to match")
}

func TestServiceRequestRepository_UpdateServiceRequest(t *testing.T) {
	repo := repository.NewServiceRequestRepository("test_service_requests.json")

	request := model.ServiceRequest{
		ID:     "request1",
		Status: "accepted",
	}

	err := repo.UpdateServiceRequest(request)
	assert.NoError(t, err, "Expected no error when updating service request")

	updatedRequest, err := repo.GetServiceRequestByID(request.ID)
	assert.NoError(t, err, "Expected no error when fetching service request")
	assert.Equal(t, request.Status, updatedRequest.Status, "Expected service request status to be updated")
}
