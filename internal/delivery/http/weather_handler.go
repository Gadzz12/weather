package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go_26_10/domain"
)

type WeatherHandler struct {
	weatherUsecase domain.WeatherUsecase
}

func NewWeatherHandler(weatherUsecase domain.WeatherUsecase) *WeatherHandler {
	return &WeatherHandler{
		weatherUsecase: weatherUsecase,
	}
}

// GetDashboard handles weather overview data (GET /api/v1/weather/dashboard)
func (h *WeatherHandler) GetDashboard(c *gin.Context) {
	var locationIDPtr *uuid.UUID
	if locStr := c.Query("location_id"); locStr != "" {
		if parsedID, err := uuid.Parse(locStr); err == nil {
			locationIDPtr = &parsedID
		}
	}

	var userIDPtr *uuid.UUID
	if val, exists := c.Get("userID"); exists {
		if uid, ok := val.(uuid.UUID); ok {
			userIDPtr = &uid
		}
	}

	data, err := h.weatherUsecase.GetDashboardData(c.Request.Context(), userIDPtr, locationIDPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.WebResponse{
			Success: false,
			Message: "Gagal memuat data cuaca: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.WebResponse{
		Success: true,
		Data:    data,
	})
}

// GetHourlyChart handles hourly temperature chart and time-series data (GET /api/v1/weather/hourly)
func (h *WeatherHandler) GetHourlyChart(c *gin.Context) {
	var locationIDPtr *uuid.UUID
	if locStr := c.Query("location_id"); locStr != "" {
		if parsedID, err := uuid.Parse(locStr); err == nil {
			locationIDPtr = &parsedID
		}
	}

	var userIDPtr *uuid.UUID
	if val, exists := c.Get("userID"); exists {
		if uid, ok := val.(uuid.UUID); ok {
			userIDPtr = &uid
		}
	}

	unit := c.Query("unit")

	data, err := h.weatherUsecase.GetHourlyChart(c.Request.Context(), userIDPtr, locationIDPtr, unit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.WebResponse{
			Success: false,
			Message: "Gagal memuat grafik suhu per jam: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.WebResponse{
		Success: true,
		Data:    data,
	})
}

// GetChanceOfRain handles rain probability analysis across time ranges (GET /api/v1/weather/chance-of-rain)
func (h *WeatherHandler) GetChanceOfRain(c *gin.Context) {
	var locationIDPtr *uuid.UUID
	if locStr := c.Query("location_id"); locStr != "" {
		if parsedID, err := uuid.Parse(locStr); err == nil {
			locationIDPtr = &parsedID
		}
	}

	var userIDPtr *uuid.UUID
	if val, exists := c.Get("userID"); exists {
		if uid, ok := val.(uuid.UUID); ok {
			userIDPtr = &uid
		}
	}

	data, err := h.weatherUsecase.GetChanceOfRainAnalysis(c.Request.Context(), userIDPtr, locationIDPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.WebResponse{
			Success: false,
			Message: "Gagal memuat analisis peluang hujan: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.WebResponse{
		Success: true,
		Data:    data,
	})
}

// GetForecast handles multi-day weather predictions with high/low temperatures (GET /api/v1/weather/forecast)
func (h *WeatherHandler) GetForecast(c *gin.Context) {
	var locationIDPtr *uuid.UUID
	if locStr := c.Query("location_id"); locStr != "" {
		if parsedID, err := uuid.Parse(locStr); err == nil {
			locationIDPtr = &parsedID
		}
	}

	var userIDPtr *uuid.UUID
	if val, exists := c.Get("userID"); exists {
		if uid, ok := val.(uuid.UUID); ok {
			userIDPtr = &uid
		}
	}

	days := 3
	if daysStr := c.Query("days"); daysStr != "" {
		if parsedDays, err := strconv.Atoi(daysStr); err == nil && parsedDays > 0 {
			days = parsedDays
		}
	}

	unit := c.Query("unit")

	data, err := h.weatherUsecase.GetForecast(c.Request.Context(), userIDPtr, locationIDPtr, days, unit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.WebResponse{
			Success: false,
			Message: "Gagal memuat prakiraan cuaca: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.WebResponse{
		Success: true,
		Data:    data,
	})
}

// CompareCities handles two-city weather comparison (GET /api/v1/weather/compare)
func (h *WeatherHandler) CompareCities(c *gin.Context) {
	city1 := c.Query("city1")
	city2 := c.Query("city2")

	var locID1Ptr *uuid.UUID
	if locStr := c.Query("location_id1"); locStr != "" {
		if parsedID, err := uuid.Parse(locStr); err == nil {
			locID1Ptr = &parsedID
		}
	}

	var locID2Ptr *uuid.UUID
	if locStr := c.Query("location_id2"); locStr != "" {
		if parsedID, err := uuid.Parse(locStr); err == nil {
			locID2Ptr = &parsedID
		}
	}

	var userIDPtr *uuid.UUID
	if val, exists := c.Get("userID"); exists {
		if uid, ok := val.(uuid.UUID); ok {
			userIDPtr = &uid
		}
	}

	unit := c.Query("unit")

	data, err := h.weatherUsecase.CompareCities(c.Request.Context(), userIDPtr, city1, city2, locID1Ptr, locID2Ptr, unit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.WebResponse{
			Success: false,
			Message: "Gagal membandingkan cuaca dua kota: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.WebResponse{
		Success: true,
		Data:    data,
	})
}

// GetCalendar handles monthly weather calendar data (GET /api/v1/weather/calendar)
func (h *WeatherHandler) GetCalendar(c *gin.Context) {
	city := c.Query("city")
	year := 0
	if yStr := c.Query("year"); yStr != "" {
		if y, err := strconv.Atoi(yStr); err == nil {
			year = y
		}
	}
	month := 0
	if mStr := c.Query("month"); mStr != "" {
		if m, err := strconv.Atoi(mStr); err == nil {
			month = m
		}
	}

	var userIDPtr *uuid.UUID
	if val, exists := c.Get("userID"); exists {
		if uid, ok := val.(uuid.UUID); ok {
			userIDPtr = &uid
		}
	}

	unit := c.Query("unit")

	data, err := h.weatherUsecase.GetCalendar(c.Request.Context(), userIDPtr, city, year, month, unit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.WebResponse{
			Success: false,
			Message: "Gagal memuat kalender cuaca: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.WebResponse{
		Success: true,
		Data:    data,
	})
}
