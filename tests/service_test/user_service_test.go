package service_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"serviceNest/model"
	"serviceNest/service"
	"serviceNest/tests/mocks"
	"testing"
)

func TestViewProfileByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	//defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	userService := service.NewUserService(mockUserRepo)

	userID := "12345"
	user := &model.User{ID: userID, Email: "test@example.com"}

	// Define test cases
	tests := []struct {
		name            string
		mockGetUserByID func(*mocks.MockUserRepository)
		expectedUser    *model.User
		expectedError   error
	}{
		{
			name: "Success",
			mockGetUserByID: func(m *mocks.MockUserRepository) {
				m.EXPECT().GetUserByID(userID).Return(user, nil)
			},
			expectedUser:  user,
			expectedError: nil,
		},
		{
			name: "User Not Found",
			mockGetUserByID: func(m *mocks.MockUserRepository) {
				m.EXPECT().GetUserByID(userID).Return(nil, nil)
			},
			expectedUser:  nil,
			expectedError: errors.New("user not found"),
		},
		{
			name: "Error Getting User",
			mockGetUserByID: func(m *mocks.MockUserRepository) {
				m.EXPECT().GetUserByID(userID).Return(nil, errors.New("database error"))
			},
			expectedUser:  nil,
			expectedError: errors.New("could not find user: database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockGetUserByID(mockUserRepo)
			user, err := userService.ViewProfileByID(userID)
			assert.Equal(t, tt.expectedUser, user)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	userService := service.NewUserService(mockUserRepo)

	userID := "12345"
	existingUser := &model.User{ID: userID, Email: "old@example.com"}

	// Define test cases
	tests := []struct {
		name               string
		mockGetUserByID    func(*mocks.MockUserRepository)
		mockGetUserByEmail func(*mocks.MockUserRepository)
		mockUpdateUser     func(*mocks.MockUserRepository)
		newEmail           *string
		newPassword        *string
		newAddress         *string
		newPhone           *string
		expectedError      error
	}{
		{
			name: "Success",
			mockGetUserByID: func(m *mocks.MockUserRepository) {
				m.EXPECT().GetUserByID(userID).Return(existingUser, nil)
			},
			mockGetUserByEmail: func(m *mocks.MockUserRepository) {
				m.EXPECT().GetUserByEmail("new@example.com").Return(nil, errors.New("user not found"))
			},
			mockUpdateUser: func(m *mocks.MockUserRepository) {
				m.EXPECT().UpdateUser(existingUser).Return(nil)
			},
			newEmail:      stringPtr("new@example.com"),
			newPassword:   stringPtr("NewPassword@123"),
			newAddress:    stringPtr("New Address"),
			newPhone:      stringPtr("1234567890"),
			expectedError: nil,
		},
		{
			name: "Email Already In Use",
			mockGetUserByID: func(m *mocks.MockUserRepository) {
				m.EXPECT().GetUserByID(userID).Return(existingUser, nil)
			},
			mockGetUserByEmail: func(m *mocks.MockUserRepository) {
				m.EXPECT().GetUserByEmail("new@example.com").Return(&model.User{ID: "67890"}, nil)
			},
			mockUpdateUser: func(m *mocks.MockUserRepository) {
				// No expectations on UpdateUser here
			},
			newEmail:      stringPtr("new@example.com"),
			expectedError: errors.New("email already in use by another user"),
		},
		{
			name: "Error Updating User",
			mockGetUserByID: func(m *mocks.MockUserRepository) {
				m.EXPECT().GetUserByID(userID).Return(existingUser, nil)
			},
			mockGetUserByEmail: func(m *mocks.MockUserRepository) {
				m.EXPECT().GetUserByEmail("new@example.com").Return(nil, errors.New("user not found"))
			},
			mockUpdateUser: func(m *mocks.MockUserRepository) {
				m.EXPECT().UpdateUser(existingUser).Return(errors.New("update error"))
			},
			newEmail:      stringPtr("new@example.com"),
			expectedError: errors.New("could not update user: update error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockGetUserByID(mockUserRepo)
			tt.mockGetUserByEmail(mockUserRepo)
			tt.mockUpdateUser(mockUserRepo)
			err := userService.UpdateUser(userID, tt.newEmail, tt.newPassword, tt.newAddress, tt.newPhone)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

// Helper function to get pointer of string
func stringPtr(s string) *string {
	return &s
}
