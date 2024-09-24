package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"serviceNest/controllers"
	"serviceNest/model"
	"serviceNest/tests/mocks"
	"testing"
)

func TestLoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	userController := controllers.NewUserController(mockUserService)

	t.Run("Successfully login", func(t *testing.T) {
		userInput := map[string]string{
			"email":    "user@example.com",
			"password": "password123",
		}
		jsonInput, _ := json.Marshal(userInput)

		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonInput))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		mockUser := &model.User{
			Email:    "user@example.com",
			Password: "$2a$10$VwW8.kUj8lB.gPiQVHGX5uOEPJFM6PgfVHLZ4g/cOhUHTBZZ6XmdS", // hashed "password123"
			Role:     "user",
			ID:       "1",
		}

		mockUserService.EXPECT().CheckUserExists("user@example.com").Return(mockUser, nil)
		mockUserService.EXPECT().CheckUserExists(gomock.Any()).Return(nil, errors.New("user not found")).AnyTimes()

		original := controllers.CheckPassword
		defer func() {
			controllers.CheckPassword = original
		}()
		controllers.CheckPassword = func(password, hash string) bool {
			return true
		}

		userController.LoginUser(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		expectedResponse := `{"status":"Success","message":"Token generate successfully","data":`
		if !bytes.Contains(rr.Body.Bytes(), []byte(expectedResponse)) {
			t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
		}
	})

	t.Run("Invalid json body", func(t *testing.T) {
		userInput := map[string]string{
			"user_email": "user@example.com",
			"password":   "password123",
		}
		jsonInput, _ := json.Marshal(userInput)

		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonInput))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		userController.LoginUser(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

	})
	t.Run("Invalid Request body", func(t *testing.T) {
		userInput := map[string]string{
			"email":    "user@example.com",
			"password": "",
		}
		jsonInput, _ := json.Marshal(userInput)

		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonInput))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		userController.LoginUser(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

	})

	t.Run("User not exist", func(t *testing.T) {
		userInput := map[string]string{
			"email":    "user@example.com",
			"password": "password123",
		}
		jsonInput, _ := json.Marshal(userInput)

		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonInput))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		mockUserService.EXPECT().CheckUserExists("user@example.com").Return(nil, errors.New("user not found")).AnyTimes()

		userController.LoginUser(rr, req)

		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
		}

		//expectedResponse := `{"status":"Success","message":"Token generate successfully","data":`
		//if !bytes.Contains(rr.Body.Bytes(), []byte(expectedResponse)) {
		//	t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
		//}
	})

	t.Run("Error generating token", func(t *testing.T) {
		userInput := map[string]string{
			"email":    "user@example.com",
			"password": "password123",
		}
		jsonInput, _ := json.Marshal(userInput)

		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonInput))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		mockUser := &model.User{
			Email:    "user@example.com",
			Password: "$2a$10$VwW8.kUj8lB.gPiQVHGX5uOEPJFM6PgfVHLZ4g/cOhUHTBZZ6XmdS", // hashed "password123"
			Role:     "user",
			ID:       "1",
		}

		mockUserService.EXPECT().CheckUserExists("user@example.com").Return(mockUser, nil)

		originalCheckPassword := controllers.CheckPassword
		defer func() {
			controllers.CheckPassword = originalCheckPassword
		}()
		controllers.CheckPassword = func(password, hash string) bool {
			return true
		}

		// Mock JWT generation to return an error
		originalGenerateJWT := controllers.GenerateJWT
		defer func() {
			controllers.GenerateJWT = originalGenerateJWT
		}()
		controllers.GenerateJWT = func(id, role string) (string, error) {
			return "", errors.New("error generating token")
		}

		// Call the login handler
		userController.LoginUser(rr, req)

		// Check that the status code is 500 Internal Server Error
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}

		// Check that the body contains the correct error message
		expectedResponse := `{"status":"Fail","message":"Error generating token"}`
		if !bytes.Contains(rr.Body.Bytes(), []byte(expectedResponse)) {
			t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
		}
	})

}
func TestLoginUser_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	userController := controllers.NewUserController(mockUserService)

	userInput := map[string]string{
		"email":    "user@example.com",
		"password": "wrongpassword",
	}
	jsonInput, _ := json.Marshal(userInput)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonInput))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	mockUser := &model.User{
		Email:    "user@example.com",
		Password: "$2a$10$VwW8.kUj8lB.gPiQVHGX5uOEPJFM6PgfVHLZ4g/cOhUHTBZZ6XmdS", // hashed "password123"
		Role:     "user",
		ID:       "1",
	}

	mockUserService.EXPECT().CheckUserExists("user@example.com").Return(mockUser, nil)

	original := controllers.CheckPassword
	defer func() {
		controllers.CheckPassword = original
	}()
	controllers.CheckPassword = func(password, hash string) bool {
		return false
	}

	userController.LoginUser(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	expectedResponse := `{"status":"Fail","message":"Invalid password"}`
	assert.JSONEq(t, expectedResponse, rr.Body.String())
}
func TestSignupUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	userController := controllers.NewUserController(mockUserService)

	t.Run("Successfully user registered", func(t *testing.T) {
		newUser := map[string]string{
			"name":     "John Doe",
			"email":    "john@example.com",
			"password": "password123",
			"role":     "user",
			"address":  "123 Main St",
			"contact":  "1234567890",
		}
		jsonInput, _ := json.Marshal(newUser)

		req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonInput))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		// Simulate no existing user
		mockUserService.EXPECT().CheckUserExists("john@example.com").Return(nil, errors.New("user not found"))

		// Simulate successful user creation
		mockUserService.EXPECT().CreateUser(gomock.Any()).Return(nil)

		userController.SignupUser(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		expectedResponse := `{"status":"Success","message":"User created successfully"}`
		assert.JSONEq(t, expectedResponse, rr.Body.String())

	})
	t.Run("User already exists", func(t *testing.T) {
		newUser := map[string]string{
			"name":     "John Doe",
			"email":    "john@example.com",
			"password": "password123",
			"role":     "user",
			"address":  "123 Main St",
			"contact":  "1234567890",
		}
		jsonInput, _ := json.Marshal(newUser)

		req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonInput))
		assert.NoError(t, err)

		user := &model.User{
			Name:     "user1",
			Password: "password",
		}
		rr := httptest.NewRecorder()

		// Simulate no existing user
		mockUserService.EXPECT().CheckUserExists("john@example.com").Return(user, nil)

		userController.SignupUser(rr, req)

		assert.Equal(t, http.StatusConflict, rr.Code)
		expectedResponse := `{"status":"Fail","message":"User already exists"}`
		assert.JSONEq(t, expectedResponse, rr.Body.String())
	})
	t.Run("Invalid request body", func(t *testing.T) {
		newUser := map[string]string{
			"username": "John Doe",
			"email":    "john@example.com",
			"password": "password123",
			"role":     "user",
			"address":  "123 Main St",
			"contact":  "1234567890",
		}
		jsonInput, _ := json.Marshal(newUser)

		req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonInput))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		// Simulate no existing user

		userController.SignupUser(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		expectedResponse := `{"status":"Fail","message":"Invalid request body"}`
		assert.JSONEq(t, expectedResponse, rr.Body.String())

	})
	t.Run("Error hashing password", func(t *testing.T) {
		newUser := map[string]string{
			"name":     "John Doe",
			"email":    "john@example.com",
			"password": "password123",
			"role":     "user",
			"address":  "123 Main St",
			"contact":  "1234567890",
		}
		jsonInput, _ := json.Marshal(newUser)

		req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonInput))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		// Simulate no existing user
		mockUserService.EXPECT().CheckUserExists("john@example.com").Return(nil, errors.New("user not found"))

		original := controllers.HashPassword
		defer func() {
			controllers.HashPassword = original
		}()
		controllers.HashPassword = func(password string) (string, error) {
			return "", errors.New("error hashing password")
		}

		userController.SignupUser(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		expectedResponse := `{"status":"Fail","message":"Error hashing password"}`
		assert.JSONEq(t, expectedResponse, rr.Body.String())

	})
	t.Run("Successfully user registered", func(t *testing.T) {
		newUser := map[string]string{
			"name":     "John Doe",
			"email":    "john@example.com",
			"password": "password123",
			"role":     "user",
			"address":  "123 Main St",
			"contact":  "1234567890",
		}
		jsonInput, _ := json.Marshal(newUser)

		req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonInput))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		// Simulate no existing user
		mockUserService.EXPECT().CheckUserExists("john@example.com").Return(nil, errors.New("user not found"))

		// Simulate successful user creation
		mockUserService.EXPECT().CreateUser(gomock.Any()).Return(errors.New("error creating user"))

		userController.SignupUser(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		expectedResponse := `{"status":"Fail","message":"Error creating user"}`
		assert.JSONEq(t, expectedResponse, rr.Body.String())

	})

}
func TestViewProfileByIDHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	userController := controllers.NewUserController(mockUserService)

	t.Run("Successfully view profile", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/users/profile", nil)
		assert.NoError(t, err)

		userID := "1"
		ctx := context.WithValue(req.Context(), "userID", userID)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		mockUser := &model.User{
			ID:    "1",
			Email: "user@example.com",
		}

		mockUserService.EXPECT().ViewProfileByID(userID).Return(mockUser, nil)

		userController.ViewProfileByIDHandler(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		expectedResponse := `{"status":"Success","message":"User profile fetched successfully","data":`
		assert.Contains(t, rr.Body.String(), expectedResponse)
	})
	t.Run("Error viewing profile", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/users/profile", nil)
		assert.NoError(t, err)

		userID := "1"
		ctx := context.WithValue(req.Context(), "userID", userID)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		mockUserService.EXPECT().ViewProfileByID(userID).Return(nil, errors.New("error viewing profile"))

		userController.ViewProfileByIDHandler(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		expectedResponse := `{"status":"Fail","message":"error viewing user"`
		assert.Contains(t, rr.Body.String(), expectedResponse)
	})

}

func TestUpdateUserHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mocks.NewMockUserService(ctrl)
	userController := controllers.NewUserController(mockUserService)

	t.Run("Valid input with no error", func(t *testing.T) {
		updateData := map[string]string{
			"email":    "newemail@example.com",
			"password": "Newpassword@123",
			"address":  "456 New Address",
			"phone":    "0987654321",
		}
		jsonInput, _ := json.Marshal(updateData)

		req, err := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(jsonInput))
		assert.NoError(t, err)

		userID := "1"
		ctx := context.WithValue(req.Context(), "userID", userID)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		// Mock the HashPassword function
		original := controllers.HashPassword
		defer func() { controllers.HashPassword = original }()
		controllers.HashPassword = func(password string) (string, error) {
			return "hashedpassword123", nil
		}

		// Mock UpdateUser service call
		mockUserService.EXPECT().UpdateUser(userID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		// Call the handler
		userController.UpdateUserHandler(rr, req)

		// Assert response
		assert.Equal(t, http.StatusOK, rr.Code)
		expectedResponse := `{"status":"Success","message":"User updated successfully"}`
		assert.JSONEq(t, expectedResponse, rr.Body.String())
	})

	t.Run("Invalid JSON input", func(t *testing.T) {
		invalidJSON := `{"email": "invalidemail.com", "password": "123"}`
		req, err := http.NewRequest("PUT", "/users/1", bytes.NewBuffer([]byte(invalidJSON)))
		assert.NoError(t, err)

		ctx := context.WithValue(req.Context(), "userID", "1")
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		// Call the handler with invalid JSON
		userController.UpdateUserHandler(rr, req)

		// Assert response
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		expectedResponse := `{"status":"Fail","message":"password must be at least 8 characters long, contain at least one uppercase letter, one lowercase letter, one number, and one special symbol"}`
		assert.JSONEq(t, expectedResponse, rr.Body.String())
	})

	t.Run("Password validation failure", func(t *testing.T) {
		updateData := map[string]string{
			"email":    "newemail@example.com",
			"password": "123", // Weak password
			"address":  "456 New Address",
			"phone":    "0987654321",
		}
		jsonInput, _ := json.Marshal(updateData)

		req, err := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(jsonInput))
		assert.NoError(t, err)

		ctx := context.WithValue(req.Context(), "userID", "1")
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		// Mock ValidatePassword to return a failure
		original := controllers.ValidatePassword
		defer func() { controllers.ValidatePassword = original }()
		controllers.ValidatePassword = func(password string) error {
			return fmt.Errorf("password is too weak")
		}

		// Call the handler
		userController.UpdateUserHandler(rr, req)

		// Assert response
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		expectedResponse := `{"status":"Fail","message":"password is too weak"}`
		assert.JSONEq(t, expectedResponse, rr.Body.String())
	})
	t.Run("Error hashing password", func(t *testing.T) {
		updateData := map[string]string{
			"email":    "newemail@example.com",
			"password": "Newpassword@123",
			"address":  "456 New Address",
			"phone":    "0987654321",
		}
		jsonInput, _ := json.Marshal(updateData)

		req, err := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(jsonInput))
		assert.NoError(t, err)

		ctx := context.WithValue(req.Context(), "userID", "1")
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		// Mock HashPassword to return an error
		original := controllers.HashPassword
		defer func() { controllers.HashPassword = original }()
		controllers.HashPassword = func(password string) (string, error) {
			return "", fmt.Errorf("hashing failed")
		}

		// Call the handler
		userController.UpdateUserHandler(rr, req)

		// Assert response
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		expectedResponse := `{"status":"Fail","message":"Error hashing password"}`
		assert.JSONEq(t, expectedResponse, rr.Body.String())
	})
	t.Run("Error during UpdateUser service call", func(t *testing.T) {
		updateData := map[string]string{
			"email":    "newemail@example.com",
			"password": "Newpassword@123",
			"address":  "456 New Address",
			"phone":    "0987654321",
		}
		jsonInput, _ := json.Marshal(updateData)

		req, err := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(jsonInput))
		assert.NoError(t, err)

		ctx := context.WithValue(req.Context(), "userID", "1")
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		// Mock HashPassword
		original := controllers.HashPassword
		defer func() { controllers.HashPassword = original }()
		controllers.HashPassword = func(password string) (string, error) {
			return "hashedpassword123", nil
		}

		// Mock UpdateUser to return an error
		mockUserService.EXPECT().UpdateUser("1", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("db update failed")).Times(1)

		// Call the handler
		userController.UpdateUserHandler(rr, req)

		// Assert response
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		expectedResponse := `{"status":"Fail","message":"Error updating user"}`
		assert.JSONEq(t, expectedResponse, rr.Body.String())
	})
}
