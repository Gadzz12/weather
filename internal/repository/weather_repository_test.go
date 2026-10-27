package repository_test

import (
	"context"
	"testing"

	"go_26_10/internal/repository"
)

func TestWeatherRepository_BMKGIntegration(t *testing.T) {
	repo := repository.NewWeatherRepository("", "https://api.bmkg.go.id")
	ctx := context.Background()

	// Jakarta coordinates
	cur, err := repo.GetRealtimeWeather(ctx, -6.2088, 106.8456)
	if err != nil {
		t.Fatalf("unexpected error fetching realtime: %v", err)
	}
	if cur.Condition == "" || cur.Temperature == 0 {
		t.Errorf("expected valid realtime weather, got %+v", cur)
	}

	high, err := repo.GetHighlights(ctx, -6.2088, 106.8456)
	if err != nil {
		t.Fatalf("unexpected error fetching highlights: %v", err)
	}
	if high.Humidity == "" || high.WindSpeed == "" {
		t.Errorf("expected valid highlights, got %+v", high)
	}

	hourly, rain, err := repo.GetHourlyForecast(ctx, -6.2088, 106.8456)
	if err != nil {
		t.Fatalf("unexpected error fetching hourly: %v", err)
	}
	if len(hourly) == 0 || len(rain) == 0 {
		t.Errorf("expected hourly and rain forecasts, got hourly: %d, rain: %d", len(hourly), len(rain))
	}

	forecast, err := repo.GetMultiDayForecast(ctx, -6.2088, 106.8456, 3)
	if err != nil {
		t.Fatalf("unexpected error fetching multiday: %v", err)
	}
	if len(forecast) < 3 {
		t.Errorf("expected 3 days forecast, got %d", len(forecast))
	}
}

func TestWeatherRepository_FallbackMode(t *testing.T) {
	repo := repository.NewWeatherRepository("", "")
	ctx := context.Background()

	cur, err := repo.GetRealtimeWeather(ctx, -6.2088, 106.8456)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cur.Condition != "Cerah Berawan" {
		t.Errorf("expected fallback condition Cerah Berawan, got %s", cur.Condition)
	}

	hourly, rain, err := repo.GetHourlyForecast(ctx, -6.2088, 106.8456)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(hourly) != 8 || len(rain) != 7 {
		t.Errorf("expected 8 hourly and 7 rain fallback items, got %d and %d", len(hourly), len(rain))
	}
}
