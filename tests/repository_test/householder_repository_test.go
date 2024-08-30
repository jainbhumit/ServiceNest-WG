package repository_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"serviceNest/model"
	"serviceNest/repository"
	"serviceNest/tests/mocks"
	"testing"
)

func TestSaveHouseholder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	repo := &repository.HouseholderRepository{
		Collection: mockCollection,
	}

	householder := model.Householder{
		User: model.User{
			ID:       "householder1",
			Name:     "John Doe",
			Email:    "johndoe@example.com",
			Password: "password123",
		},
	}

	mockCollection.EXPECT().
		InsertOne(gomock.Any(), householder).
		Return(&mongo.InsertOneResult{}, nil)

	err := repo.SaveHouseholder(householder)
	assert.NoError(t, err)
}

func TestSaveHouseholder_InsertError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	repo := &repository.HouseholderRepository{
		Collection: mockCollection,
	}

	householder := model.Householder{
		User: model.User{
			ID:       "householder1",
			Name:     "John Doe",
			Email:    "johndoe@example.com",
			Password: "password123",
		},
	}

	mockCollection.EXPECT().
		InsertOne(gomock.Any(), householder).
		Return(nil, errors.New("insert error"))

	err := repo.SaveHouseholder(householder)
	assert.Error(t, err)
	assert.EqualError(t, err, "insert error")
}

func TestGetHouseholderByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := mocks.NewMockMongoCollection(ctrl)
	repo := &repository.HouseholderRepository{
		Collection: mockCollection,
	}

	householder := &model.Householder{
		User: model.User{
			ID:    "householder1",
			Name:  "John Doe",
			Email: "johndoe@example.com",
		},
	}

	mockResult := mongo.NewSingleResultFromDocument(householder, nil, nil)

	mockCollection.EXPECT().
		FindOne(gomock.Any(), gomock.Eq(bson.M{"id": "householder1"})).
		Return(mockResult)

	result, err := repo.GetHouseholderByID("householder1")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, householder, result)
}

//func TestGetHouseholderByID_NotFound(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockCollection := mocks.NewMockMongoCollection(ctrl)
//	repo := &repository.HouseholderRepository{
//		Collection: mockCollection,
//	}
//
//	mockCollection.EXPECT().
//		FindOne(gomock.Any(), gomock.Eq(bson.M{"id": "householder1"})).
//		Return(NewSingleResultFromError(mongo.ErrNoDocuments))
//
//	result, err := repo.GetHouseholderByID("householder1")
//	assert.Error(t, err)
//	assert.Nil(t, result)
//	assert.EqualError(t, err, "householder not found")
//}
