package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"go_26_10/domain"
)

func TestLocationHandler(t *testing.T) {
	router := setupTestRouter()

	// Register user to obtain valid JWT token
	regPayload := domain.RegisterRequest{
		Name:     "Location User",
		Email:    "locuser@example.com",
		Password: "password123",
	}
	regBody, _ := json.Marshal(regPayload)
	regReq, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(regBody))
	regReq.Header.Set("Content-Type", "application/json")
	regW := httptest.NewRecorder()
	router.ServeHTTP(regW, regReq)

	var regResp struct {
		Data domain.AuthResponse `json:"data"`
	}
	_ = json.Unmarshal(regW.Body.Bytes(), &regResp)
	token := regResp.Data.Token

	var targetLocationID uuid.UUID

	t.Run("Search Location - Jakarta", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/locations/search?q=Jakarta", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var resp struct {
			Success bool              `json:"success"`
			Data    []domain.Location `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(resp.Data) == 0 {
			t.Fatalf("expected at least 1 search result for Jakarta")
		}

		if resp.Data[0].City != "Jakarta" {
			t.Errorf("expected city Jakarta, got %s", resp.Data[0].City)
		}

		targetLocationID = resp.Data[0].ID
	})

	t.Run("Add Saved Location", func(t *testing.T) {
		payload := domain.SaveLocationRequest{
			LocationID: targetLocationID,
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/locations/saved", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got %d. Body: %s", w.Code, w.Body.String())
		}

		var resp struct {
			Success bool                     `json:"success"`
			Data    domain.UserSavedLocation `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if resp.Data.LocationID != targetLocationID {
			t.Errorf("expected LocationID %s, got %s", targetLocationID, resp.Data.LocationID)
		}
	})

	t.Run("Prevent Duplicate Saved Location", func(t *testing.T) {
		payload := domain.SaveLocationRequest{
			LocationID: targetLocationID,
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/locations/saved", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400 for duplicate saved location, got %d", w.Code)
		}
	})

	t.Run("Get All Saved Locations", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/locations/saved", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var resp struct {
			Success bool                               `json:"success"`
			Data    []domain.SavedLocationItemResponse `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(resp.Data) != 1 {
			t.Fatalf("expected 1 saved location, got %d", len(resp.Data))
		}

		if resp.Data[0].Location.City != "Jakarta" {
			t.Errorf("expected saved city Jakarta, got %s", resp.Data[0].Location.City)
		}
		if resp.Data[0].Location.Country != "Indonesia" {
			t.Errorf("expected country Indonesia, got %s", resp.Data[0].Location.Country)
		}
		if !resp.Data[0].IsDefault {
			t.Errorf("expected first saved location to be marked default")
		}
	})

	t.Run("Search Non-Indonesian Location Returns Empty", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/locations/search?q=Tokyo", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var resp struct {
			Success bool              `json:"success"`
			Data    []domain.Location `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(resp.Data) != 0 {
			t.Errorf("expected 0 results for foreign city Tokyo, got %d", len(resp.Data))
		}
	})
}
