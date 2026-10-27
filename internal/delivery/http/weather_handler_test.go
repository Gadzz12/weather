package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go_26_10/domain"
	deliveryHttp "go_26_10/internal/delivery/http"
	"go_26_10/internal/repository"
	"go_26_10/internal/usecase"
)

func setupTestRouter() http.Handler {
	weatherRepo := repository.NewWeatherRepository("", "https://api.bmkg.go.id")
	locationRepo := repository.NewLocationRepository(nil)
	userRepo := repository.NewUserRepository(nil)

	weatherUsecase := usecase.NewWeatherUsecase(weatherRepo, locationRepo, userRepo)
	userUsecase := usecase.NewUserUsecase(userRepo, "secret")
	locationUsecase := usecase.NewLocationUsecase(locationRepo)

	authHandler := deliveryHttp.NewAuthHandler(userUsecase)
	userHandler := deliveryHttp.NewUserHandler(userUsecase)
	locationHandler := deliveryHttp.NewLocationHandler(locationUsecase)
	weatherHandler := deliveryHttp.NewWeatherHandler(weatherUsecase)

	return deliveryHttp.SetupRouter(deliveryHttp.RouterConfig{
		AuthHandler:     authHandler,
		UserHandler:     userHandler,
		LocationHandler: locationHandler,
		WeatherHandler:  weatherHandler,
		JWTSecret:       "secret",
	})
}

func TestWeatherDashboardHandler(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/weather/dashboard", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var response struct {
		Success bool                        `json:"success"`
		Data    domain.WeatherDashboardData `json:"data"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !response.Success {
		t.Errorf("expected success to be true")
	}

	// Verify Location
	if response.Data.Location.City == "" {
		t.Errorf("expected location city to be populated")
	}

	// Verify Current Weather
	if response.Data.CurrentWeather.Condition == "" {
		t.Errorf("expected current weather condition to be populated")
	}

	// Verify Highlights
	if response.Data.Highlights.Humidity == "" {
		t.Errorf("expected highlights humidity to be populated")
	}
	if response.Data.Highlights.Sunrise == "" {
		t.Errorf("expected highlights sunrise to be populated")
	}
	if response.Data.Highlights.Sunset == "" {
		t.Errorf("expected highlights sunset to be populated")
	}
	if response.Data.Highlights.WindSpeed == "" {
		t.Errorf("expected highlights wind_speed to be populated")
	}

	// Verify Hourly Temperature time-series
	if len(response.Data.HourlyTemperature) == 0 {
		t.Errorf("expected hourly temperature list to be populated")
	}

	// Verify Chance of Rain
	if len(response.Data.ChanceOfRain) == 0 {
		t.Errorf("expected chance of rain list to be populated")
	}

	// Verify 3-Days Forecast
	if len(response.Data.ThreeDaysForecast) != 3 {
		t.Errorf("expected 3 days forecast, got %d", len(response.Data.ThreeDaysForecast))
	}
	for _, f := range response.Data.ThreeDaysForecast {
		if !strings.HasPrefix(f.TempHighDisplay, "↑") || !strings.HasPrefix(f.TempLowDisplay, "↓") {
			t.Errorf("expected max/min indicators ↑ and ↓, got high: %s, low: %s", f.TempHighDisplay, f.TempLowDisplay)
		}
	}

	// Verify HourlyChart & RainAnalysis inside dashboard
	if response.Data.HourlyChart == nil {
		t.Errorf("expected hourly_chart to be populated in dashboard")
	}
	if response.Data.RainAnalysis == nil {
		t.Errorf("expected rain_analysis to be populated in dashboard")
	}
}

func TestHourlyTemperatureChartHandler(t *testing.T) {
	router := setupTestRouter()

	t.Run("Default Celsius", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/weather/hourly", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var response struct {
			Success bool                   `json:"success"`
			Data    domain.HourlyChartData `json:"data"`
		}

		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if !response.Success {
			t.Errorf("expected success to be true")
		}

		if response.Data.Unit != "Celsius" {
			t.Errorf("expected unit Celsius, got %s", response.Data.Unit)
		}

		if len(response.Data.Labels) == 0 || len(response.Data.DataPoints) == 0 {
			t.Errorf("expected chart labels and data_points to have elements")
		}

		if len(response.Data.Labels) != len(response.Data.DataPoints) {
			t.Errorf("expected labels count to match data_points count")
		}

		if response.Data.PeakTime == "" {
			t.Errorf("expected peak_time to be populated")
		}

		if response.Data.LowestTime == "" {
			t.Errorf("expected lowest_time to be populated")
		}

		if response.Data.Trend == "" {
			t.Errorf("expected trend description to be populated")
		}
	})

	t.Run("Fahrenheit Unit Conversion", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/weather/hourly?unit=fahrenheit", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var response struct {
			Success bool                   `json:"success"`
			Data    domain.HourlyChartData `json:"data"`
		}

		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if response.Data.Unit != "Fahrenheit" {
			t.Errorf("expected unit Fahrenheit, got %s", response.Data.Unit)
		}

		// When temp is 10 C -> 50 F
		if response.Data.MinTemp < 32 {
			t.Errorf("expected min temp in Fahrenheit to be >= 32, got %f", response.Data.MinTemp)
		}
	})
}

func TestChanceOfRainAnalysisHandler(t *testing.T) {
	router := setupTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/weather/chance-of-rain", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var response struct {
		Success bool                `json:"success"`
		Data    domain.RainAnalysis `json:"data"`
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !response.Success {
		t.Errorf("expected success to be true")
	}

	if response.Data.PeakPercentage <= 0 {
		t.Errorf("expected positive peak percentage, got %d", response.Data.PeakPercentage)
	}

	if response.Data.PeakTime == "" {
		t.Errorf("expected peak_time to be non-empty")
	}

	if response.Data.PeakIntensity == "" {
		t.Errorf("expected peak_intensity to be non-empty")
	}

	if response.Data.RiskLevel == "" {
		t.Errorf("expected risk_level to be populated")
	}

	if response.Data.Advisory == "" {
		t.Errorf("expected advisory recommendation to be non-empty")
	}

	if len(response.Data.TimeRanges) == 0 {
		t.Errorf("expected time_ranges period analysis to have elements")
	}

	for _, tr := range response.Data.TimeRanges {
		if tr.Period == "" || tr.TimeSlot == "" || tr.Description == "" {
			t.Errorf("expected complete time range breakdown, got %+v", tr)
		}
	}
}

func TestMultiDayForecastHandler(t *testing.T) {
	router := setupTestRouter()

	t.Run("Default 3 Days Forecast in Celsius", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/weather/forecast", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var response struct {
			Success bool                `json:"success"`
			Data    domain.ForecastData `json:"data"`
		}

		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if !response.Success {
			t.Errorf("expected success to be true")
		}

		if response.Data.Unit != "Celsius" {
			t.Errorf("expected unit Celsius, got %s", response.Data.Unit)
		}

		if len(response.Data.Forecast) != 3 {
			t.Errorf("expected 3 forecast days, got %d", len(response.Data.Forecast))
		}

		for _, item := range response.Data.Forecast {
			if !strings.HasPrefix(item.TempHighDisplay, "↑") {
				t.Errorf("expected TempHighDisplay to start with ↑, got %s", item.TempHighDisplay)
			}
			if !strings.HasPrefix(item.TempLowDisplay, "↓") {
				t.Errorf("expected TempLowDisplay to start with ↓, got %s", item.TempLowDisplay)
			}
			if item.TempHigh < item.TempLow {
				t.Errorf("expected TempHigh (%f) to be >= TempLow (%f)", item.TempHigh, item.TempLow)
			}
		}
	})

	t.Run("Forecast in Fahrenheit", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/weather/forecast?unit=fahrenheit", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var response struct {
			Success bool                `json:"success"`
			Data    domain.ForecastData `json:"data"`
		}

		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if response.Data.Unit != "Fahrenheit" {
			t.Errorf("expected unit Fahrenheit, got %s", response.Data.Unit)
		}

		for _, item := range response.Data.Forecast {
			if !strings.Contains(item.TempHighDisplay, "°F") {
				t.Errorf("expected °F in TempHighDisplay, got %s", item.TempHighDisplay)
			}
		}
	})
}

func TestWeatherCompareHandler(t *testing.T) {
	router := setupTestRouter()

	t.Run("Compare default cities (Jakarta vs Bandung)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/weather/compare", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var response struct {
			Success bool                     `json:"success"`
			Data    domain.WeatherComparison `json:"data"`
		}

		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if !response.Success {
			t.Errorf("expected success to be true")
		}
		if response.Data.City1.Location.City != "Jakarta" {
			t.Errorf("expected City1 Jakarta, got %s", response.Data.City1.Location.City)
		}
		if response.Data.City2.Location.City != "Bandung" {
			t.Errorf("expected City2 Bandung, got %s", response.Data.City2.Location.City)
		}
		if len(response.Data.Comparison.HourlyComparison) == 0 {
			t.Errorf("expected hourly comparison to be populated")
		}
	})

	t.Run("Compare custom cities by name (Jakarta vs Surabaya)", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/weather/compare?city1=Jakarta&city2=Surabaya", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var response struct {
			Success bool                     `json:"success"`
			Data    domain.WeatherComparison `json:"data"`
		}

		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if response.Data.City1.Location.City != "Jakarta" {
			t.Errorf("expected City1 Jakarta, got %s", response.Data.City1.Location.City)
		}
		if response.Data.City2.Location.City != "Surabaya" {
			t.Errorf("expected City2 Surabaya, got %s", response.Data.City2.Location.City)
		}
		if len(response.Data.Comparison.HourlyComparison) == 0 {
			t.Errorf("expected hourly comparison to be populated")
		}
	})
}

func TestWeatherCalendarHandler(t *testing.T) {
	router := setupTestRouter()

	t.Run("Get calendar data default city", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/weather/calendar?city=Jakarta&year=2026&month=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var response struct {
			Success bool                       `json:"success"`
			Data    domain.WeatherCalendarData `json:"data"`
		}

		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if !response.Success {
			t.Errorf("expected success to be true")
		}

		if response.Data.City != "Jakarta" {
			t.Errorf("expected City 'Jakarta', got '%s'", response.Data.City)
		}

		if response.Data.Year != 2026 {
			t.Errorf("expected Year 2026, got %d", response.Data.Year)
		}

		if response.Data.Month != 10 {
			t.Errorf("expected Month 10, got %d", response.Data.Month)
		}

		if len(response.Data.Days) != 31 {
			t.Errorf("expected 31 days in October, got %d", len(response.Data.Days))
		}

		// Verify first day has proper fields
		firstDay := response.Data.Days[0]
		if firstDay.Day != 1 {
			t.Errorf("expected day 1, got %d", firstDay.Day)
		}
		if firstDay.Condition == "" {
			t.Errorf("expected condition to be non-empty")
		}
	})

	t.Run("Get calendar data with Fahrenheit", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/weather/calendar?city=Jakarta&year=2026&month=10&unit=fahrenheit", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var response struct {
			Success bool                       `json:"success"`
			Data    domain.WeatherCalendarData `json:"data"`
		}

		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if !response.Success {
			t.Errorf("expected success to be true")
		}

		if len(response.Data.Days) != 31 {
			t.Errorf("expected 31 days, got %d", len(response.Data.Days))
		}
		// In Fahrenheit, temp_high should be > 60
		if response.Data.Days[0].TempHigh < 60 {
			t.Errorf("expected Fahrenheit temperature > 60, got %f", response.Data.Days[0].TempHigh)
		}
	})
}
