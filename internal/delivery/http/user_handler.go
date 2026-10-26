package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go_26_10/domain"
)

type UserHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(userUsecase domain.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

// UpdateSettings handles updating user preferences (PUT /api/v1/users/settings)
func (h *UserHandler) UpdateSettings(c *gin.Context) {
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

	var req domain.UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.WebResponse{
			Success: false,
			Message: "Validasi gagal: " + err.Error(),
		})
		return
	}

	updatedUser, err := h.userUsecase.UpdateSettings(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.WebResponse{
			Success: false,
			Message: "Gagal memperbarui pengaturan: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.WebResponse{
		Success: true,
		Message: "Pengaturan preferensi suhu berhasil diperbarui",
		Data:    updatedUser,
	})
}

// GetProfile handles fetching user profile (GET /api/v1/users/profile)
func (h *UserHandler) GetProfile(c *gin.Context) {
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

	user, err := h.userUsecase.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, domain.WebResponse{
			Success: false,
			Message: "Profil pengguna tidak ditemukan: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.WebResponse{
		Success: true,
		Data:    user,
	})
}
