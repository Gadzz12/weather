package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go_26_10/domain"
)

type LocationHandler struct {
	locationUsecase domain.LocationUsecase
}

func NewLocationHandler(locationUsecase domain.LocationUsecase) *LocationHandler {
	return &LocationHandler{
		locationUsecase: locationUsecase,
	}
}

// Search handles searching locations by query string (GET /api/v1/locations/search?q={city_name})
func (h *LocationHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, domain.WebResponse{
			Success: false,
			Message: "Parameter query 'q' wajib disertakan",
		})
		return
	}

	locations, err := h.locationUsecase.SearchLocations(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.WebResponse{
			Success: false,
			Message: "Gagal mencari lokasi: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.WebResponse{
		Success: true,
		Data:    locations,
	})
}

// AddSaved handles bookmarking a location (POST /api/v1/locations/saved)
func (h *LocationHandler) AddSaved(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, domain.WebResponse{
			Success: false,
			Message: "Pengguna tidak terautentikasi",
		})
		return
	}
	userID, ok := val.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, domain.WebResponse{
			Success: false,
			Message: "Format token pengguna tidak valid",
		})
		return
	}

	var req domain.SaveLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.WebResponse{
			Success: false,
			Message: "Validasi gagal: " + err.Error(),
		})
		return
	}

	saved, err := h.locationUsecase.SaveLocation(c.Request.Context(), userID, req.LocationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.WebResponse{
			Success: false,
			Message: "Gagal menyimpan lokasi: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, domain.WebResponse{
		Success: true,
		Message: "Lokasi berhasil disimpan ke daftar favorit",
		Data:    saved,
	})
}

// GetSaved handles retrieving user's saved locations (GET /api/v1/locations/saved)
func (h *LocationHandler) GetSaved(c *gin.Context) {
	val, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, domain.WebResponse{
			Success: false,
			Message: "Pengguna tidak terautentikasi",
		})
		return
	}
	userID, ok := val.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, domain.WebResponse{
			Success: false,
			Message: "Format token pengguna tidak valid",
		})
		return
	}

	savedList, err := h.locationUsecase.GetSavedLocations(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.WebResponse{
			Success: false,
			Message: "Gagal memuat daftar lokasi tersimpan: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.WebResponse{
		Success: true,
		Data:    savedList,
	})
}
