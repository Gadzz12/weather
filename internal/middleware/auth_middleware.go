package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
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

// allowedOrigins membaca daftar origin yang diizinkan dari variabel lingkungan ALLOWED_ORIGINS
// (pisahkan dengan koma). Jika tidak diatur atau bernilai "*", semua origin diizinkan agar memudahkan hosting frontend dan backend.
func allowedOrigins() []string {
	env := strings.TrimSpace(os.Getenv("ALLOWED_ORIGINS"))
	if env == "" || env == "*" {
		return []string{"*"}
	}
	var result []string
	for _, o := range strings.Split(env, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			result = append(result, o)
		}
	}
	if len(result) == 0 {
		return []string{"*"}
	}
	return result
}

// isOriginAllowed memeriksa apakah origin terdapat dalam daftar yang diizinkan atau wildcard (*).
func isOriginAllowed(origin string, allowed []string) bool {
	for _, o := range allowed {
		if o == "*" || o == origin {
			return true
		}
	}
	return false
}

// CORSMiddleware handles cross-origin resource sharing requests using github.com/gin-contrib/cors
func CORSMiddleware() gin.HandlerFunc {
	env := strings.TrimSpace(os.Getenv("ALLOWED_ORIGINS"))
	allowed := allowedOrigins()

	corsConfig := cors.Config{
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"Accept",
			"Cache-Control",
			"X-Requested-With",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}

	if env == "" || env == "*" || isOriginAllowed("*", allowed) {
		corsConfig.AllowOriginFunc = func(origin string) bool {
			return true
		}
	} else {
		corsConfig.AllowOrigins = allowed
	}

	return cors.New(corsConfig)
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
