package repository

import (
	"github.com/stretchr/testify/assert"
	"serviceNest/model"
	"serviceNest/repository"
	"testing"
)

func TestServiceRepository_SaveService(t *testing.T) {
	repo := repository.NewServiceRepository("test_services.json")

	service := model.Service{
		ID:          "1",
		Name:        "Plumbing",
		Description: "Fixing pipes",
	}

	err := repo.SaveService(service)
	assert.NoError(t, err, "Expected no error when saving service")

	savedService, err := repo.GetServiceByID(service.ID)
	assert.NoError(t, err, "Expected no error when fetching service")
	assert.Equal(t, service.Name, savedService.Name, "Expected service name to match")
}

func TestServiceRepository_RemoveService(t *testing.T) {
	repo := repository.NewServiceRepository("test_services.json")

	err := repo.RemoveService("1")
	assert.NoError(t, err, "Expected no error when removing service")

	_, err = repo.GetServiceByID("1")
	assert.Error(t, err, "Expected an error when fetching removed service")
}
