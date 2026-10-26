package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go_26_10/domain"
)

func TestAuthAndUserSettings(t *testing.T) {
	router := setupTestRouter()

	var token string
	email := "testuser@example.com"
	password := "password123"

	t.Run("Register New User", func(t *testing.T) {
		payload := domain.RegisterRequest{
			Name:     "Test User",
			Email:    email,
			Password: password,
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got %d. Body: %s", w.Code, w.Body.String())
		}

		var resp struct {
			Success bool                `json:"success"`
			Data    domain.AuthResponse `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if !resp.Success || resp.Data.Token == "" {
			t.Fatalf("expected valid token in auth response")
		}

		if resp.Data.User.TemperatureUnit != "Celsius" {
			t.Errorf("expected default unit Celsius, got %s", resp.Data.User.TemperatureUnit)
		}

		token = resp.Data.Token
	})

	t.Run("Login User", func(t *testing.T) {
		payload := domain.LoginRequest{
			Email:    email,
			Password: password,
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var resp struct {
			Success bool                `json:"success"`
			Data    domain.AuthResponse `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if !resp.Success || resp.Data.Token == "" {
			t.Fatalf("expected valid token in login response")
		}
		token = resp.Data.Token
	})

	t.Run("Update Settings - Unauthorized without Token", func(t *testing.T) {
		payload := domain.UpdateSettingsRequest{
			TemperatureUnit: "Fahrenheit",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPut, "/api/v1/users/settings", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected status 401, got %d", w.Code)
		}
	})

	t.Run("Update Settings - Success with Bearer Token", func(t *testing.T) {
		payload := domain.UpdateSettingsRequest{
			TemperatureUnit: "Fahrenheit",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPut, "/api/v1/users/settings", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var resp struct {
			Success bool        `json:"success"`
			Data    domain.User `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.Data.TemperatureUnit != "Fahrenheit" {
			t.Errorf("expected updated unit Fahrenheit, got %s", resp.Data.TemperatureUnit)
		}
	})

	t.Run("Get Profile with Bearer Token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/profile", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var resp struct {
			Success bool        `json:"success"`
			Data    domain.User `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.Data.Email != email {
			t.Errorf("expected email %s, got %s", email, resp.Data.Email)
		}
		if resp.Data.TemperatureUnit != "Fahrenheit" {
			t.Errorf("expected temperature_unit Fahrenheit, got %s", resp.Data.TemperatureUnit)
		}
	})
}
