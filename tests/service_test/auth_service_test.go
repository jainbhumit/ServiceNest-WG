package service_test

//func TestSignUp(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockUserRepo := mocks.NewMockUserRepository(ctrl)
//
//	// Set expectations for the mock
//	mockUserRepo.EXPECT().
//		GetUserByEmail("test@example.com").
//		Return(nil, nil)
//
//	mockUserRepo.EXPECT().
//		SaveUser(gomock.Any()).
//		Return(nil)
//
//	// Mocking user inputs (could be automated using a helper function)
//	// Ensure the function SignUp uses the mock correctly
//	user, err := service.SignUp(mockUserRepo)
//	assert.NoError(t, err)
//	assert.NotNil(t, user)
//	assert.Equal(t, "test@example.com", user.Email)
//}

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
//	mockUser := &model.User{
//		Email:    "test@example.com",
//		Password: string(hashedPassword),
//	}
//
//	// Set expectations for the mock
//	mockUserRepo.EXPECT().
//		GetUserByEmail("test@example.com").
//		Return(mockUser, nil)
//
//	// Mocking user inputs (could be automated using a helper function)
//	// Ensure the function Login uses the mock correctly
//	user, err := service.Login(mockUserRepo)
//	assert.NoError(t, err)
//	assert.NotNil(t, user)
//	assert.Equal(t, "test@example.com", user.Email)
//}
