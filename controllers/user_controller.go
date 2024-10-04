package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
	"net/http"
	"serviceNest/interfaces"
	"serviceNest/logger"
	"serviceNest/model"
	"serviceNest/response"
	"serviceNest/util"
)

var CheckPassword = util.CheckPasswordHash
var HashPassword = util.HashPassword
var validate *validator.Validate
var GenerateJWT = util.GenerateJWT
var ValidatePassword = util.ValidatePassword

func init() {
	// initialize new validator
	validate = validator.New()
}

type UserController struct {
	userService interfaces.UserService
}

func NewUserController(userService interfaces.UserService) *UserController {
	return &UserController{userService: userService}
}

// LoginUser handles POST /login
func (u *UserController) LoginUser(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	var err error
	if err = json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid input", 1001)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	// validate userInput
	err = validate.Struct(userInput)
	if err != nil {
		logger.Error("Validation error", nil)
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", 1001)
		return
	}

	// Check if user exists and get the user details
	var user *model.User
	user, err = u.userService.CheckUserExists(userInput.Email)
	if err != nil {
		logger.Error("Invalid email or password", map[string]interface{}{"email": userInput.Email})
		response.ErrorResponse(w, http.StatusUnauthorized, "Invalid email or password", 1005)
		//http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Verify the password
	if !CheckPassword(userInput.Password, user.Password) {
		logger.Error("Invalid password", map[string]interface{}{"email": userInput.Email})
		response.ErrorResponse(w, http.StatusUnauthorized, "Invalid password", 1005)
		//http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}
	if user.IsActive == false {
		logger.Error("User is deactivated by admin", map[string]interface{}{"email": userInput.Email})
		response.ErrorResponse(w, http.StatusUnauthorized, "user Deactivated by admin", 1007)
		return
	}

	// Generate JWT token
	var tokenString string
	tokenString, err = GenerateJWT(user.ID, user.Role)
	if err != nil {
		logger.Error("Error generating token", map[string]interface{}{"email": userInput.Email})
		response.ErrorResponse(w, http.StatusInternalServerError, "Error generating token", 1006)
		//http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Return the token as JSON
	logger.Info("token generated", map[string]interface{}{"email": userInput.Email})
	response.SuccessResponse(w, map[string]interface{}{"token": tokenString}, "Token generate successfully", http.StatusCreated)
	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
func (u *UserController) SignupUser(w http.ResponseWriter, r *http.Request) {
	var newUser struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
		Role     string `json:"role" validate:"required"`
		Address  string `json:"address" validate:"required"`
		Contact  string `json:"contact" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		logger.Error("Invalid input", nil)
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid input", 1001)
		return
	}
	err := validate.Struct(newUser)
	if err != nil {
		logger.Error("Validation error", nil)
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", 1001)
		return
	}
	// Check if user already exists
	exists, err := u.userService.CheckUserExists(newUser.Email)

	if exists != nil && err == nil {
		logger.Error("User already exists", map[string]interface{}{"email": newUser.Email})
		response.ErrorResponse(w, http.StatusConflict, "User already exists", 1009)
		//http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Hash the password before saving
	hashedPassword, err := HashPassword(newUser.Password)
	if err != nil {
		logger.Error("Error hashing password", map[string]interface{}{"email": newUser.Email})
		response.ErrorResponse(w, http.StatusInternalServerError, "Error hashing password", 1006)
		//http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	newUser.Password = hashedPassword

	// Save the user
	user := &model.User{
		Name:     newUser.Name,
		Email:    newUser.Email,
		Password: newUser.Password,
		Role:     newUser.Role,
		Address:  newUser.Address,
		Contact:  newUser.Contact,
	}
	err = u.userService.CreateUser(user)
	if err != nil {
		logger.Error("Error creating user", map[string]interface{}{"email": newUser.Email})
		response.ErrorResponse(w, http.StatusInternalServerError, "Error creating user", 1006)
		//http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	logger.Info("User created sucessfully", map[string]interface{}{"email": newUser.Email})
	response.SuccessResponse(w, nil, "User created successfully", http.StatusCreated)
	//w.WriteHeader(http.StatusCreated)
	//
	//json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

// ViewProfileByIDHandler handles GET /users/{id}
func (u *UserController) ViewProfileByIDHandler(w http.ResponseWriter, r *http.Request) {

	userID := r.Context().Value("userID").(string)

	// Call the UserService to get the user profile
	user, err := u.userService.ViewProfileByID(userID)
	if err != nil {
		logger.Error(err.Error(), nil)
		response.ErrorResponse(w, http.StatusNotFound, "error viewing user", 1008)
		//http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	logger.Info("user profile fetched", map[string]interface{}{"email": user.Email})
	response.SuccessResponse(w, user, "User profile fetched successfully", http.StatusOK)

}

// UpdateUserHandler handles PUT /users/{id}
func (u *UserController) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	// Parse incoming JSON data
	var updateData struct {
		Email    *string `json:"email"`
		Password *string `json:"password"`
		Address  *string `json:"address"`
		Contact  *string `json:"contact"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", 1001)
		//http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	err := ValidatePassword(*updateData.Password)
	if err != nil {
		logger.Error("Error validating password", map[string]interface{}{"email": *updateData.Email})
		response.ErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("%v", err), 1001)
		//http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}
	hashedPassword, err := HashPassword(*updateData.Password)
	if err != nil {
		logger.Error("Error hashing password", map[string]interface{}{"email": *updateData.Email})
		response.ErrorResponse(w, http.StatusInternalServerError, "Error hashing password", 1006)
		//http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Call the UserService to update the user profile
	err = u.userService.UpdateUser(userID, updateData.Email, &hashedPassword, updateData.Address, updateData.Contact)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, "Error updating user", 1006)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Info("user updated sucessfully", map[string]interface{}{"email": *updateData.Email})
	response.SuccessResponse(w, nil, "User updated successfully", http.StatusOK)

}
