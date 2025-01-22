package middlewares_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"serviceNest/middlewares"
	"serviceNest/util"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "userID123", r.Context().Value("userID"))
		assert.Equal(t, "Householder", r.Context().Value("role"))
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Valid Token", func(t *testing.T) {
		tokenString, _ := util.GenerateJWT("userID123", "Householder") // Implement this in your util
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)

		rr := httptest.NewRecorder()
		middleware := middlewares.AuthMiddleware(handler)

		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Missing Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rr := httptest.NewRecorder()
		middleware := middlewares.AuthMiddleware(handler)

		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer invalidToken")
		rr := httptest.NewRecorder()
		middleware := middlewares.AuthMiddleware(handler)

		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}

func TestHouseHolderAuthMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Valid Role", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		ctx := context.WithValue(req.Context(), "role", "Householder")
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		middleware := middlewares.HouseHolderAuthMiddleware(handler)

		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Invalid Role", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		ctx := context.WithValue(req.Context(), "role", "ServiceProvider")
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		middleware := middlewares.HouseHolderAuthMiddleware(handler)

		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}

func TestServiceProviderAuthMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Valid Role", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		ctx := context.WithValue(req.Context(), "role", "ServiceProvider")
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		middleware := middlewares.ServiceProviderAuthMiddleware(handler)

		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Invalid Role", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		ctx := context.WithValue(req.Context(), "role", "Householder")
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		middleware := middlewares.ServiceProviderAuthMiddleware(handler)

		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}

func TestAdminAuthMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Valid Role", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		ctx := context.WithValue(req.Context(), "role", "Admin")
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		middleware := middlewares.AdminAuthMiddleware(handler)

		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Invalid Role", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		ctx := context.WithValue(req.Context(), "role", "User")
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		middleware := middlewares.AdminAuthMiddleware(handler)

		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}
