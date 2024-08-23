package service

import (
	"serviceNest/model"
	"serviceNest/repository"
	"serviceNest/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHouseholderService_AddReview(t *testing.T) {
	householderRepo := repository.NewHouseholderRepository("test_user.json")
	//reviewRepo := repository.NewReviewRepository("test_reviews.json")
	serviceRepo := repository.NewServiceRepository("test_services.json")
	providerRepo := repository.NewServiceProviderRepository("test_providers.json")
	serviceRequestRepo := repository.NewServiceRequestRepository("test_service_requests.json")
	householderService := service.NewHouseholderService(householderRepo, providerRepo, serviceRepo, serviceRequestRepo)

	review := model.Review{
		ID:            "1",
		ServiceID:     "service1",
		HouseholderID: "householder1",
		Rating:        4.5,
		Comments:      "Great service!",
		ReviewDate:    time.Now(),
	}

	err := householderService.AddReview(review.HouseholderID, review.ServiceID, review.Comments, review.Rating)
	assert.NoError(t, err, "Expected no error when adding review")

	//savedReview, err := reviewRepo.GetReviewByID(review.ID)
	//assert.NoError(t, err, "Expected no error when fetching review")
	//assert.Equal(t, review.Comments, savedReview.Comments, "Expected review comments to match")
}

//func TestHouseholderService_ViewProfileByID(t *testing.T) {
//	userRepo := repository.NewUserRepository("test_users.json")
//	householderService := NewHouseholderService(userRepo)
//
//	profile, err := householderService.ViewProfileByID("householder1")
//	assert.NoError(t, err, "Expected no error when viewing profile")
//	assert.Equal(t, "householder1", profile.ID, "Expected profile ID to match")
//}
