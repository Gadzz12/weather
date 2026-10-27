package domain

import (
	"context"

	"github.com/google/uuid"
)

// CurrentWeather represents real-time weather information
type CurrentWeather struct {
	Temperature   float64 `json:"temperature"`
	Condition     string  `json:"condition"`
	ConditionIcon string  `json:"condition_icon"`
	DateTime      string  `json:"date_time"`
}

// WeatherHighlights represents key climate metrics
type WeatherHighlights struct {
	Precipitation string `json:"precipitation"`
	Humidity      string `json:"humidity"`
	WindSpeed     string `json:"wind_speed"`
	Sunrise       string `json:"sunrise"`
	Sunset        string `json:"sunset"`
}

// HourlyTemperature represents hourly temperature projection for chart and dashboard
type HourlyTemperature struct {
	Time      string  `json:"time"`
	Hour      int     `json:"hour,omitempty"`
	Temp      float64 `json:"temp"`
	TempC     float64 `json:"temp_c,omitempty"`
	TempF     float64 `json:"temp_f,omitempty"`
	FeelsLike float64 `json:"feels_like,omitempty"`
	Humidity  int     `json:"humidity,omitempty"`
	Condition string  `json:"condition,omitempty"`
	Icon      string  `json:"icon,omitempty"`
}

// HourlyChartData represents complete time-series and analytical structure for hourly temperature chart
type HourlyChartData struct {
	Location   WeatherLocationInfo `json:"location"`
	Unit       string              `json:"unit"`
	MinTemp    float64             `json:"min_temp"`
	MaxTemp    float64             `json:"max_temp"`
	AvgTemp    float64             `json:"avg_temp"`
	PeakTime   string              `json:"peak_time"`
	LowestTime string              `json:"lowest_time"`
	Trend      string              `json:"trend"`
	Labels     []string            `json:"labels"`
	DataPoints []float64           `json:"data_points"`
	Series     []HourlyTemperature `json:"series"`
}

// ChanceOfRain represents hourly rain likelihood and intensity
type ChanceOfRain struct {
	Time       string  `json:"time"`
	Hour       int     `json:"hour,omitempty"`
	Intensity  string  `json:"intensity"`
	Percentage int     `json:"percentage"`
	PrecipMM   float64 `json:"precip_mm,omitempty"`
}

// RainTimeRange represents aggregated rain probability across a time slot of the day
type RainTimeRange struct {
	Period        string `json:"period"`
	TimeSlot      string `json:"time_slot"`
	AverageChance int    `json:"average_chance"`
	MaxChance     int    `json:"max_chance"`
	Intensity     string `json:"intensity"`
	Description   string `json:"description"`
}

// RainAnalysis represents comprehensive chance of rain analytics
type RainAnalysis struct {
	Location       WeatherLocationInfo `json:"location"`
	Date           string              `json:"date"`
	DailyAverage   int                 `json:"daily_average"`
	PeakTime       string              `json:"peak_time"`
	PeakPercentage int                 `json:"peak_percentage"`
	PeakIntensity  string              `json:"peak_intensity"`
	RiskLevel      string              `json:"risk_level"`
	Advisory       string              `json:"advisory"`
	TimeRanges     []RainTimeRange     `json:"time_ranges"`
	HourlyForecast []ChanceOfRain      `json:"hourly_forecast"`
}

// DayForecast represents a daily weather summary with min/max temperatures
type DayForecast struct {
	Day             string  `json:"day"`
	Date            string  `json:"date,omitempty"`
	Condition       string  `json:"condition"`
	TempHigh        float64 `json:"temp_high"`
	TempLow         float64 `json:"temp_low"`
	TempHighDisplay string  `json:"temp_high_display,omitempty"` // e.g. "↑ 26°C" or "↑ 78.8°F"
	TempLowDisplay  string  `json:"temp_low_display,omitempty"`  // e.g. "↓ 11°C" or "↓ 51.8°F"
	Icon            string  `json:"icon"`
	Summary         string  `json:"summary,omitempty"`
}

// ForecastData represents multi-day weather prediction response
type ForecastData struct {
	Location WeatherLocationInfo `json:"location"`
	Unit     string              `json:"unit"`
	Days     int                 `json:"days"`
	Forecast []DayForecast       `json:"forecast"`
}

// WeatherLocationInfo represents simplified location info for dashboard header
type WeatherLocationInfo struct {
	ID       uuid.UUID `json:"id"`
	City     string    `json:"city"`
	State    *string   `json:"state,omitempty"`
	Country  string    `json:"country"`
	ImageURL *string   `json:"image_url,omitempty"`
}

// WeatherDashboardData represents the aggregate dashboard payload
type WeatherDashboardData struct {
	Location          WeatherLocationInfo `json:"location"`
	CurrentWeather    CurrentWeather      `json:"current_weather"`
	Highlights        WeatherHighlights   `json:"highlights"`
	HourlyTemperature []HourlyTemperature `json:"hourly_temperature"`
	ChanceOfRain      []ChanceOfRain      `json:"chance_of_rain"`
	ThreeDaysForecast []DayForecast       `json:"three_days_forecast"`
	HourlyChart       *HourlyChartData    `json:"hourly_chart,omitempty"`
	RainAnalysis      *RainAnalysis       `json:"rain_analysis,omitempty"`
}

// HourlyComparisonItem represents comparison of temperature at a specific hour
type HourlyComparisonItem struct {
	Time  string  `json:"time"`
	Temp1 float64 `json:"temp1"`
	Temp2 float64 `json:"temp2"`
	Diff  float64 `json:"diff"`
}

// ComparisonSummary contains analytical comparison between two cities
type ComparisonSummary struct {
	TempDifference     float64                `json:"temp_difference"`
	TempDiffFormatted  string                 `json:"temp_diff_formatted"`
	WarmerCity         string                 `json:"warmer_city"`
	ColderCity         string                 `json:"colder_city"`
	HumidityComparison string                 `json:"humidity_comparison"`
	RainierCity        string                 `json:"rainier_city"`
	WindierCity        string                 `json:"windier_city"`
	ConditionSummary   string                 `json:"condition_summary"`
	Recommendation     string                 `json:"recommendation"`
	HourlyComparison   []HourlyComparisonItem `json:"hourly_comparison"`
}

// WeatherComparison represents the complete comparison data of two cities
type WeatherComparison struct {
	City1      WeatherDashboardData `json:"city1"`
	City2      WeatherDashboardData `json:"city2"`
	Comparison ComparisonSummary    `json:"comparison"`
}

// CalendarDay represents weather forecast and metrics for a specific calendar date
type CalendarDay struct {
	Date         string  `json:"date"`
	Day          int     `json:"day"`
	Weekday      string  `json:"weekday"`
	Condition    string  `json:"condition"`
	TempHigh     float64 `json:"temp_high"`
	TempLow      float64 `json:"temp_low"`
	ChanceOfRain int     `json:"chance_of_rain"`
	Icon         string  `json:"icon"`
	Advisory     string  `json:"advisory"`
	IsToday      bool    `json:"is_today"`
}

// WeatherCalendarData represents the monthly calendar weather overview
type WeatherCalendarData struct {
	City      string        `json:"city"`
	Year      int           `json:"year"`
	Month     int           `json:"month"`
	MonthName string        `json:"month_name"`
	Days      []CalendarDay `json:"days"`
}

// WeatherRepository defines external Weather API provider integration contract
type WeatherRepository interface {
	GetRealtimeWeather(ctx context.Context, lat, lon float64) (*CurrentWeather, error)
	GetHighlights(ctx context.Context, lat, lon float64) (*WeatherHighlights, error)
	GetHourlyForecast(ctx context.Context, lat, lon float64) ([]HourlyTemperature, []ChanceOfRain, error)
	GetMultiDayForecast(ctx context.Context, lat, lon float64, days int) ([]DayForecast, error)
	GetCalendarForecast(ctx context.Context, city string, year, month int) (*WeatherCalendarData, error)
}

// WeatherUsecase defines the business logic contract for weather analytics
type WeatherUsecase interface {
	GetDashboardData(ctx context.Context, userID *uuid.UUID, locationID *uuid.UUID) (*WeatherDashboardData, error)
	GetHourlyChart(ctx context.Context, userID *uuid.UUID, locationID *uuid.UUID, unit string) (*HourlyChartData, error)
	GetChanceOfRainAnalysis(ctx context.Context, userID *uuid.UUID, locationID *uuid.UUID) (*RainAnalysis, error)
	GetForecast(ctx context.Context, userID *uuid.UUID, locationID *uuid.UUID, days int, unit string) (*ForecastData, error)
	CompareCities(ctx context.Context, userID *uuid.UUID, query1, query2 string, locID1, locID2 *uuid.UUID, unit string) (*WeatherComparison, error)
	GetCalendar(ctx context.Context, userID *uuid.UUID, city string, year, month int, unit string) (*WeatherCalendarData, error)
}
