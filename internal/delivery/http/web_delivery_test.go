package http_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWebFrontendServing(t *testing.T) {
	router := setupTestRouter()

	t.Run("Serve index.html at root GET /", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Accept", "text/html,application/xhtml+xml")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		body := w.Body.String()
		if !strings.Contains(body, "Weather Forecast & Analytics Dashboard") {
			t.Errorf("expected HTML title in response body, got: %s", body[:min(200, len(body))])
		}
	})

	t.Run("Serve style.css at GET /style.css", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/style.css", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		if !strings.Contains(w.Body.String(), ".dashboard-wrapper") {
			t.Errorf("expected CSS content in response body")
		}
	})

	t.Run("Serve app.js at GET /app.js", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/app.js", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		if !strings.Contains(w.Body.String(), "[go_vibe]") {
			t.Errorf("expected JavaScript content in response body")
		}
	})

	t.Run("Root endpoint with Accept application/json returns API endpoints JSON", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Accept", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		if !strings.Contains(w.Body.String(), "Weather") || !strings.Contains(w.Body.String(), "endpoints") {
			t.Errorf("expected JSON API overview in response body, got: %s", w.Body.String())
		}
	})

	t.Run("Verify Calendar feature components in web files", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Accept", "text/html")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		body := w.Body.String()
		if !strings.Contains(body, "calendarGrid") {
			t.Errorf("expected calendarGrid in index.html")
		}
		if !strings.Contains(body, "calendarNavBtn") {
			t.Errorf("expected calendarNavBtn in index.html")
		}
		if !strings.Contains(body, "Kalender Cuaca") {
			t.Errorf("expected Kalender Cuaca in index.html")
		}

		// Verify CSS has calendar styles
		reqCss, _ := http.NewRequest(http.MethodGet, "/style.css", nil)
		wCss := httptest.NewRecorder()
		router.ServeHTTP(wCss, reqCss)
		if !strings.Contains(wCss.Body.String(), ".calendar-grid") {
			t.Errorf("expected .calendar-grid in style.css")
		}

		// Verify JS has calendar logic
		reqJs, _ := http.NewRequest(http.MethodGet, "/app.js", nil)
		wJs := httptest.NewRecorder()
		router.ServeHTTP(wJs, reqJs)
		if !strings.Contains(wJs.Body.String(), "setupCalendarFeature") {
			t.Errorf("expected setupCalendarFeature in app.js")
		}
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
