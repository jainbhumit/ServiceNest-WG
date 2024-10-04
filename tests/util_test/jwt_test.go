package util_test

import (
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"serviceNest/config"
	"serviceNest/util"
	"strings"
	"testing"
	"time"
)

func TestGenerateJWT(t *testing.T) {
	userID := "testUserID"
	role := "Householder"

	// Generate JWT
	tokenString, err := util.GenerateJWT(userID, role)

	// Assert no error occurred
	assert.NoError(t, err, "expected no error when generating JWT")

	// Assert token is not empty
	assert.NotEmpty(t, tokenString, "expected generated token not to be empty")

	// Assert the token contains the user ID and role
	assert.True(t, strings.Contains(tokenString, "."), "expected a valid JWT with three parts")

	// Verify the token contains the correct claims
	token, err := util.VerifyJWT(tokenString)
	assert.NoError(t, err, "expected no error when verifying JWT")

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		assert.Equal(t, userID, claims["user_id"], "expected the user ID in the claims to match")
		assert.Equal(t, role, claims["role"], "expected the role in the claims to match")
	} else {
		t.Error("failed to extract claims from token")
	}
}

func TestVerifyJWT_InvalidToken(t *testing.T) {
	// Create an invalid token string (just random text)
	invalidToken := "this.is.an.invalid.token"

	// Try to verify the invalid token
	_, err := util.VerifyJWT(invalidToken)

	// Assert an error is returned for an invalid token
	assert.Error(t, err, "expected an error when verifying an invalid token")
}

func TestVerifyJWT_ExpiredToken(t *testing.T) {
	// Create a token that is already expired
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = "expiredUserID"
	claims["role"] = "Admin"
	claims["exp"] = time.Now().Add(-time.Hour * 1).Unix() // Expired 1 hour ago

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredTokenString, _ := token.SignedString([]byte(config.SECRET))

	// Verify the expired token
	_, err := util.VerifyJWT(expiredTokenString)

	// Assert an error is returned for an expired token
	assert.Error(t, err, "expected an error when verifying an expired token")
}

func TestVerifyJWT_ValidToken(t *testing.T) {
	// Create a valid token
	userID := "validUserID"
	role := "Admin"
	tokenString, err := util.GenerateJWT(userID, role)

	assert.NoError(t, err, "expected no error when generating a valid token")

	// Verify the token
	token, err := util.VerifyJWT(tokenString)

	// Assert no error is returned for a valid token
	assert.NoError(t, err, "expected no error when verifying a valid token")

	// Assert the token is valid and contains the correct claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		assert.Equal(t, userID, claims["user_id"], "expected the user ID in the claims to match")
		assert.Equal(t, role, claims["role"], "expected the role in the claims to match")
	} else {
		t.Error("failed to extract claims from token")
	}
}
