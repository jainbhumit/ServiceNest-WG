package tests

import (
	"serviceNest/model"
	"time"
)

var MockServices = []model.Service{
	{
		ID:          "service1",
		Name:        "Household",
		Description: "House cleaning service",
		Price:       50.0,
		ProviderID:  "provider1",
		Category:    "Household",
	},
	{
		ID:          "service2",
		Name:        "Plumbing",
		Description: "Fix plumbing issues",
		Price:       75.0,
		ProviderID:  "provider2",
		Category:    "Maintenance",
	},
}

var MockServiceProvider = []model.ServiceProvider{
	{
		User: model.User{
			ID:      "provider1",
			Name:    "provider",
			Address: "test address",
			Contact: "9999999999",
		},
		ServicesOffered: []model.Service{{
			ID:          "service1",
			Name:        "Household",
			Description: "House cleaning service",
			Price:       50.0,
			ProviderID:  "provider1",
			Category:    "Household",
		}},
		Rating:       0,
		Reviews:      []*model.Review{},
		Availability: true,
		IsActive:     true,
	},
}
var MockHouseUser = []model.User{
	{
		ID:      "householder1",
		Name:    "householder",
		Address: "test address",
		Contact: "0000000000",
	},
}

var MockHouseholder = []model.Householder{
	{
		User: model.User{
			ID:      "householder1",
			Name:    "householder",
			Address: "test address",
			Contact: "0000000000",
		},
	},
}

var MockServiceRequest = []model.ServiceRequest{
	{
		ID:              "request1",
		HouseholderID:   &MockHouseholder[0].ID,
		ServiceID:       "Household",
		RequestedTime:   time.Now(),
		Status:          "pending",
		ProviderDetails: nil,
	},
}
var MockServiceRequestWithAccept = []model.ServiceRequest{
	{
		ID:                "request1",
		HouseholderID:     &MockHouseholder[0].ID,
		ServiceID:         "Household",
		ServiceProviderID: "provider1",
		RequestedTime:     time.Now(),
		Status:            "Accepted",
		ProviderDetails: &model.ServiceProviderDetails{
			Name:    "provider",
			Contact: "9999999999",
			Address: "test address",
			Rating:  0,
			Reviews: []*model.Review{},
		},
	},
}
