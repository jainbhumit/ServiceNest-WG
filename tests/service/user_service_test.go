package service

//
//import (
//	"serviceNest/model"
//	"serviceNest/repository"
//	"serviceNest/service"
//	"testing"
//
//	"github.com/stretchr/testify/assert"
//)
//
//func TestUserService_UpdateUser(t *testing.T) {
//	userRepo := repository.NewUserRepository("test_users.json")
//	userService := service.NewUserService(userRepo)
//
//	user := model.User{
//		ID:       "1",
//		Email:    "test@example.com",
//		Password: "password",
//		Contact:  "1234567890",
//		Address:  "123 Street",
//	}
//
//	err := userService.UpdateUser(user.ID, &user.Email, &user.Password, &user.Address, &user.Contact)
//	assert.NoError(t, err, "")
//
//	updatedUser, err := userRepo.GetUserByID(user.ID)
//	assert.NoError(t, err, "")
//	assert.Equal(t, user.Email, updatedUser.Email, "")
//}
//
//func TestUserService_ViewProfileByID(t *testing.T) {
//	userRepo := repository.NewUserRepository("test_users.json")
//	userService := service.NewUserService(userRepo)
//
//	user := model.User{
//		ID:    "1",
//		Email: "test@example.com",
//	}
//
//	viewedUser, err := userService.ViewProfileByID(user.ID)
//	assert.NoError(t, err, "Expected no error when viewing profile")
//	assert.Equal(t, user.Email, viewedUser.Email, "Expected email to match")
//}
