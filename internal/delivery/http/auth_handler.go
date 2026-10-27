package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go_26_10/domain"
)

type AuthHandler struct {
	userUsecase domain.UserUsecase
}

func NewAuthHandler(userUsecase domain.UserUsecase) *AuthHandler {
	return &AuthHandler{
		userUsecase: userUsecase,
	}
}

// Register handles user registration (POST /api/v1/auth/register)
func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.WebResponse{
			Success: false,
			Message: "Validasi gagal: " + err.Error(),
		})
		return
	}

	res, err := h.userUsecase.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.WebResponse{
			Success: false,
			Message: "Registrasi gagal: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, domain.WebResponse{
		Success: true,
		Message: "Registrasi pengguna berhasil",
		Data:    res,
	})
}

// Login handles user authentication (POST /api/v1/auth/login)
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.WebResponse{
			Success: false,
			Message: "Validasi gagal: " + err.Error(),
		})
		return
	}

	res, err := h.userUsecase.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.WebResponse{
			Success: false,
			Message: "Login gagal: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.WebResponse{
		Success: true,
		Message: "Login berhasil",
		Data:    res,
	})
}
