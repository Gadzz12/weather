package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"go_26_10/internal/repository"
	"go_26_10/internal/usecase"
)

func TestWeatherUsecase_HourlyChart(t *testing.T) {
	weatherRepo := repository.NewWeatherRepository("", "")
	locationRepo := repository.NewLocationRepository(nil)
	userRepo := repository.NewUserRepository(nil)
	weatherUc := usecase.NewWeatherUsecase(weatherRepo, locationRepo, userRepo)

	ctx := context.Background()

	t.Run("Hourly chart in Celsius", func(t *testing.T) {
		chart, err := weatherUc.GetHourlyChart(ctx, nil, nil, "celsius")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if chart.Unit != "Celsius" {
			t.Errorf("expected unit Celsius, got %s", chart.Unit)
		}

		if len(chart.Labels) != 8 {
			t.Errorf("expected 8 data points, got %d", len(chart.Labels))
		}

		if chart.MinTemp != 10 {
			t.Errorf("expected min temp 10, got %f", chart.MinTemp)
		}

		if chart.MaxTemp != 22 {
			t.Errorf("expected max temp 22, got %f", chart.MaxTemp)
		}

		if chart.PeakTime != "03 pm" {
			t.Errorf("expected peak time '03 pm', got '%s'", chart.PeakTime)
		}

		if chart.LowestTime != "06 am" {
			t.Errorf("expected lowest time '06 am', got '%s'", chart.LowestTime)
		}
	})

	t.Run("Hourly chart in Fahrenheit", func(t *testing.T) {
		chart, err := weatherUc.GetHourlyChart(ctx, nil, nil, "fahrenheit")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if chart.Unit != "Fahrenheit" {
			t.Errorf("expected unit Fahrenheit, got %s", chart.Unit)
		}

		// 10 C = 50 F
		if chart.MinTemp != 50 {
			t.Errorf("expected min temp 50 F, got %f", chart.MinTemp)
		}

		// 22 C = 71.6 F
		if chart.MaxTemp != 71.6 {
			t.Errorf("expected max temp 71.6 F, got %f", chart.MaxTemp)
		}
	})
}

func TestWeatherUsecase_ChanceOfRain(t *testing.T) {
	weatherRepo := repository.NewWeatherRepository("", "")
	locationRepo := repository.NewLocationRepository(nil)
	userRepo := repository.NewUserRepository(nil)
	weatherUc := usecase.NewWeatherUsecase(weatherRepo, locationRepo, userRepo)

	ctx := context.Background()

	rainAnalysis, err := weatherUc.GetChanceOfRainAnalysis(ctx, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rainAnalysis.PeakPercentage != 90 {
		t.Errorf("expected peak percentage 90, got %d", rainAnalysis.PeakPercentage)
	}

	if rainAnalysis.PeakTime != "06 pm" {
		t.Errorf("expected peak time '06 pm', got '%s'", rainAnalysis.PeakTime)
	}

	if rainAnalysis.PeakIntensity != "Heavy Rain" {
		t.Errorf("expected peak intensity 'Heavy Rain', got '%s'", rainAnalysis.PeakIntensity)
	}

	if rainAnalysis.RiskLevel != "Very High" {
		t.Errorf("expected risk level 'Very High', got '%s'", rainAnalysis.RiskLevel)
	}

	if len(rainAnalysis.TimeRanges) == 0 {
		t.Errorf("expected daily time ranges breakdown")
	}

	// Verify location is populated
	expectedID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	if rainAnalysis.Location.ID != expectedID {
		t.Errorf("expected default location id %s, got %s", expectedID, rainAnalysis.Location.ID)
	}
}

func TestWeatherUsecase_CompareCities(t *testing.T) {
	weatherRepo := repository.NewWeatherRepository("", "")
	locationRepo := repository.NewLocationRepository(nil)
	userRepo := repository.NewUserRepository(nil)
	weatherUc := usecase.NewWeatherUsecase(weatherRepo, locationRepo, userRepo)

	ctx := context.Background()

	t.Run("Default comparison between Jakarta and Bandung", func(t *testing.T) {
		res, err := weatherUc.CompareCities(ctx, nil, "", "", nil, nil, "celsius")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if res.City1.Location.City != "Jakarta" {
			t.Errorf("expected City1 to be Jakarta, got %s", res.City1.Location.City)
		}
		if res.City2.Location.City != "Bandung" {
			t.Errorf("expected City2 to be Bandung, got %s", res.City2.Location.City)
		}
		if len(res.Comparison.HourlyComparison) == 0 {
			t.Errorf("expected hourly comparison items")
		}
		if res.Comparison.TempDiffFormatted == "" {
			t.Errorf("expected non-empty temp diff formatted description")
		}
	})

	t.Run("Explicit city query Jakarta and Surabaya", func(t *testing.T) {
		res, err := weatherUc.CompareCities(ctx, nil, "Jakarta", "Surabaya", nil, nil, "celsius")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if res.City1.Location.City != "Jakarta" {
			t.Errorf("expected City1 to be Jakarta, got %s", res.City1.Location.City)
		}
		if res.City2.Location.City != "Surabaya" {
			t.Errorf("expected City2 to be Surabaya, got %s", res.City2.Location.City)
		}
		if len(res.Comparison.HourlyComparison) == 0 {
			t.Errorf("expected hourly comparison items")
		}
	})
}

func TestWeatherUsecase_Calendar(t *testing.T) {
	weatherRepo := repository.NewWeatherRepository("", "")
	locationRepo := repository.NewLocationRepository(nil)
	userRepo := repository.NewUserRepository(nil)
	weatherUc := usecase.NewWeatherUsecase(weatherRepo, locationRepo, userRepo)

	ctx := context.Background()

	t.Run("Calendar for default month and year in Celsius", func(t *testing.T) {
		cal, err := weatherUc.GetCalendar(ctx, nil, "Jakarta", 2026, 10, "celsius")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cal.City != "Jakarta" {
			t.Errorf("expected city Jakarta, got %s", cal.City)
		}
		if cal.Year != 2026 {
			t.Errorf("expected year 2026, got %d", cal.Year)
		}
		if cal.Month != 10 {
			t.Errorf("expected month 10, got %d", cal.Month)
		}
		if len(cal.Days) != 31 {
			t.Errorf("expected 31 days, got %d", len(cal.Days))
		}

		// First day checks
		d1 := cal.Days[0]
		if d1.Day != 1 {
			t.Errorf("expected day 1, got %d", d1.Day)
		}
		if d1.Weekday != "Thursday" {
			t.Errorf("expected Thursday, got %s", d1.Weekday)
		}
	})

	t.Run("Calendar in Fahrenheit", func(t *testing.T) {
		cal, err := weatherUc.GetCalendar(ctx, nil, "Bandung", 2026, 10, "fahrenheit")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(cal.Days) != 31 {
			t.Errorf("expected 31 days, got %d", len(cal.Days))
		}
		// Temp high in Fahrenheit should be converted (> 40)
		if cal.Days[0].TempHigh < 40 {
			t.Errorf("expected Fahrenheit temp > 40, got %f", cal.Days[0].TempHigh)
		}
	})
}
