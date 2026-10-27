package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"go_26_10/domain"
)

// JWTClaims defines JWT claims payload
type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// CORSMiddleware menangani permintaan cross-origin secara dinamis dan aman untuk domain manapun
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Mengambil origin browser (misal http://localhost:5500, http://127.0.0.1:5500, https://...)
		if origin != "" && origin != "null" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		// Mengizinkan header yang diminta oleh preflight browser secara dinamis
		reqHeaders := c.Request.Header.Get("Access-Control-Request-Headers")
		if reqHeaders != "" {
			c.Writer.Header().Set("Access-Control-Allow-Headers", reqHeaders)
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE, HEAD")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Authorization")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		// Preflight OPTIONS langsung dihentikan dengan status 204 No Content
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// ParseJWTToken parses and validates a JWT token string
func ParseJWTToken(tokenStr, secretKey string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// AuthMiddleware validates JWT Bearer token from Authorization header
func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, domain.WebResponse{
				Success: false,
				Message: "Authorization header is required",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, domain.WebResponse{
				Success: false,
				Message: "Authorization header must be Bearer token",
			})
			return
		}

		claims, err := ParseJWTToken(parts[1], secretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, domain.WebResponse{
				Success: false,
				Message: "Token autentikasi tidak valid atau sudah kedaluwarsa: " + err.Error(),
			})
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}

// OptionalAuthMiddleware extracts user info if Bearer token is provided, but doesn't block if absent
func OptionalAuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				if claims, err := ParseJWTToken(parts[1], secretKey); err == nil {
					c.Set("userID", claims.UserID)
					c.Set("email", claims.Email)
				}
			}
		}
		c.Next()
	}
}
