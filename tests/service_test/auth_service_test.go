//go:build !test
// +build !test

package service_test

//
//import (
//	"bytes"
//	"github.com/golang/mock/gomock"
//	"github.com/stretchr/testify/assert"
//	"golang.org/x/crypto/bcrypt"
//	"io"
//	"serviceNest/interfaces/mocks"
//	"serviceNest/model"
//	"serviceNest/service"
//	"testing"
//)
//
//func TestSignUp(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockUserRepo := mocks.NewMockUserRepository(ctrl)
//
//	// Simulate user inputs
//	input := "John Doe\njohn@example.com\nPassword123!\n1\n123 Main St\n1234567890\n"
//	service.SetInputReader(io.NopCloser(bytes.NewReader([]byte(input))))
//
//	// Set expectations for the mock
//	mockUserRepo.EXPECT().
//		GetUserByEmail("john@example.com").
//		Return(nil, nil)
//
//	mockUserRepo.EXPECT().
//		SaveUser(gomock.Any()).
//		Return(nil)
//
//	// Execute SignUp
//	user, err := service.SignUp(mockUserRepo)
//	assert.NoError(t, err)
//	assert.NotNil(t, user)
//	assert.Equal(t, "john@example.com", user.Email)
//	assert.Equal(t, "Householder", user.Role)
//}
//
//func TestLogin(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockUserRepo := mocks.NewMockUserRepository(ctrl)
//
//	// Create a hashed password for the mock user
//	password := "password123"
//	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
//
////	mockUser := &model.User{
////		Email:    "john@example.com",
////		Password: string(hashedPassword),
////	}
////
////	// Set expectations for the mock
////	mockUserRepo.EXPECT().
////		GetUserByEmail("john@example.com").
////		Return(mockUser, nil)
////
////	// Simulate user inputs
////	input := "john@example.com\npassword123\n"
////	service.SetInputReader(io.NopCloser(bytes.NewReader([]byte(input))))
////
////	// Execute Login
////	user, err := service.Login(mockUserRepo)
////	assert.NoError(t, err)
////	assert.NotNil(t, user)
////	assert.Equal(t, "john@example.com", u
//
//// TestSignUp_Success tests the successful scenario for SignUp
//func TestSignUp_Success(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockUserRepo := mocks.NewMockUserRepository(ctrl)
//
//	input := "John Doe\njohn.doe@example.com\nPassword123!\n1\n123 Main St\n1234567890\n"
//	r := strings.NewReader(input)
//	service.SetInputReader(r)
//
//	mockUserRepo.EXPECT().GetUserByEmail("john.doe@example.com").Return(nil, nil)
//	mockUserRepo.EXPECT().SaveUser(gomock.Any()).Return(nil)
//
//	user, err := service.SignUp(mockUserRepo)
//
//	assert.NoError(t, err)
//	assert.Equal(t, "John Doe", user.Name)
//	assert.Equal(t, "john.doe@example.com", user.Email)
//	assert.Equal(t, "Householder", user.Role)
//	assert.Equal(t, "123 Main St", user.Address)
//	assert.Equal(t, "1234567890", user.Contact)
//
//	// Check if the password is hashed correctly
//	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("Password123!"))
//	assert.NoError(t, err, "password should be hashed correctly")
//}
//
//// TestSignUp_EmailAlreadyExists tests the scenario where the email is already registered
//func TestSignUp_EmailAlreadyExists(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockUserRepo := mocks.NewMockUserRepository(ctrl)
//
//	input := "John Doe\njohn.doe@example.com\n"
//	r := strings.NewReader(input)
//	service.SetInputReader(r)
//
//	existingUser := &model.User{Email: "john.doe@example.com"}
//	mockUserRepo.EXPECT().GetUserByEmail("john.doe@example.com").Return(existingUser, nil)
//
//	_, err := service.SignUp(mockUserRepo)
//	assert.Error(t, err)
//	assert.Equal(t, "email already registered. Please use a different email address.", err.Error())
//}
//
//// TestLogin_Success tests the successful scenario for Login
//func TestLogin_Success(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockUserRepo := mocks.NewMockUserRepository(ctrl)
//
//	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)
//	existingUser := &model.User{
//		Email:    "john.doe@example.com",
//		Password: string(hashedPassword),
//	}
//	mockUserRepo.EXPECT().GetUserByEmail("john.doe@example.com").Return(existingUser, nil)
//
//	input := "john.doe@example.com\nPassword123!\n"
//	r := strings.NewReader(input)
//	service.SetInputReader(r)
//
//	user, err := service.Login(mockUserRepo)
//	assert.NoError(t, err)
//	assert.Equal(t, existingUser, user)
//}
//
//// TestLogin_InvalidCredentials tests the scenario where the user enters invalid credentials
//func TestLogin_InvalidCredentials(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockUserRepo := mocks.NewMockUserRepository(ctrl)
//
//	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)
//	existingUser := &model.User{
//		Email:    "john.doe@example.com",
//		Password: string(hashedPassword),
//	}
//	mockUserRepo.EXPECT().GetUserByEmail("john.doe@example.com").Return(existingUser, nil)
//
//	input := "john.doe@example.com\nWrongPassword!\n"
//	r := strings.NewReader(input)
//	service.SetInputReader(r)
//
//	_, err := service.Login(mockUserRepo)
//	assert.Error(t, err)
//	assert.Equal(t, "invalid credentials", err.Error())
//}
//
//// TestLogin_UserNotFound tests the scenario where the user is not found
//func TestLogin_UserNotFound(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockUserRepo := mocks.NewMockUserRepository(ctrl)
//
//	mockUserRepo.EXPECT().GetUserByEmail("unknown@example.com").Return(nil, errors.New("user not found"))
//
//	input := "unknown@example.com\nPassword123!\n"
//	r := strings.NewReader(input)
//	service.SetInputReader(r)
//
//	_, err := service.Login(mockUserRepo)
//	assert.Error(t, err)
//	assert.Equal(t, "user not found", err.Error())
//}
