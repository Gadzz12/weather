package http

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"go_26_10/domain"
	"go_26_10/internal/middleware"
)

type RouterConfig struct {
	AuthHandler     *AuthHandler
	UserHandler     *UserHandler
	LocationHandler *LocationHandler
	WeatherHandler  *WeatherHandler
	JWTSecret       string
}

// SetupRouter initializes Gin engine and maps all RESTful API routes
func SetupRouter(cfg RouterConfig) *gin.Engine {
	r := gin.Default()

	// Global Middlewares (CORS, Logger, Recovery)
	r.Use(middleware.CORSMiddleware())
	r.Use(gin.Recovery())

	// Deteksi lokasi folder web secara fleksibel (root workspace, folder go_vibe, maupun subpackage test)
	candidates := []string{
		"./web",
		"./go_vibe/web",
		"../web",
		"../../web",
		"../../../web",
	}
	webDir := ""
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			if _, err := os.Stat(filepath.Join(c, "index.html")); err == nil {
				webDir = c
				break
			}
		}
	}

	hasWebFiles := (webDir != "")
	indexPath := ""
	if hasWebFiles {
		indexPath = filepath.Join(webDir, "index.html")
		r.Static("/web", webDir)
		r.StaticFile("/index.html", indexPath)
		r.StaticFile("/style.css", filepath.Join(webDir, "style.css"))
		r.StaticFile("/app.js", filepath.Join(webDir, "app.js"))
	}

	apiEndpointsDoc := gin.H{
		"health":         "GET /health",
		"dashboard":      "GET /api/v1/weather/dashboard",
		"hourly_chart":   "GET /api/v1/weather/hourly",
		"chance_of_rain": "GET /api/v1/weather/chance-of-rain",
		"forecast":       "GET /api/v1/weather/forecast",
		"calendar":       "GET /api/v1/weather/calendar?city={city}&year={year}&month={month}",
		"compare":        "GET /api/v1/weather/compare?city1={city}&city2={city}",
		"auth":           []string{"POST /api/v1/auth/register", "POST /api/v1/auth/login", "POST /login", "POST /register"},
		"users":          []string{"PUT /api/v1/users/settings", "GET /api/v1/users/profile"},
		"locations":      []string{"GET /api/v1/locations/search?q={city}", "POST /api/v1/locations/saved", "GET /api/v1/locations/saved"},
	}

	// Root endpoint: Menyajikan Web UI jika diakses melalui browser, atau API endpoints jika diminta JSON
	r.GET("/", func(c *gin.Context) {
		accept := c.GetHeader("Accept")
		if strings.Contains(accept, "application/json") || !hasWebFiles {
			c.JSON(http.StatusOK, gin.H{
				"success":   true,
				"message":   "Welcome to go_vibe Weather & Analytics API",
				"endpoints": apiEndpointsDoc,
			})
			return
		}

		c.File(indexPath)
	})

	// Dedicated API documentation endpoint
	r.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success":   true,
			"message":   "Welcome to go_vibe Weather & Analytics API",
			"endpoints": apiEndpointsDoc,
		})
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, domain.WebResponse{
			Success: true,
			Message: "go_vibe weather & analytics service is healthy",
		})
	})

	// Alias rute auth langsung di root (untuk kompatibilitas frontend fetch /login & /register)
	r.POST("/login", cfg.AuthHandler.Login)
	r.POST("/register", cfg.AuthHandler.Register)

	apiV1 := r.Group("/api/v1")
	{
		// A. Authentication (Public)
		authGroup := apiV1.Group("/auth")
		{
			authGroup.POST("/register", cfg.AuthHandler.Register)
			authGroup.POST("/login", cfg.AuthHandler.Login)
		}

		// Protected routes with JWT Auth
		authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)
		optionalAuthMiddleware := middleware.OptionalAuthMiddleware(cfg.JWTSecret)

		// Settings / Users (Protected)
		userGroup := apiV1.Group("/users")
		userGroup.Use(authMiddleware)
		{
			userGroup.PUT("/settings", cfg.UserHandler.UpdateSettings)
			userGroup.GET("/profile", cfg.UserHandler.GetProfile)
		}

		// B. Weather & Dashboard Analytics (Authenticated / Public)
		weatherGroup := apiV1.Group("/weather")
		weatherGroup.Use(optionalAuthMiddleware)
		{
			weatherGroup.GET("/dashboard", cfg.WeatherHandler.GetDashboard)
			weatherGroup.GET("/hourly", cfg.WeatherHandler.GetHourlyChart)
			weatherGroup.GET("/chance-of-rain", cfg.WeatherHandler.GetChanceOfRain)
			weatherGroup.GET("/forecast", cfg.WeatherHandler.GetForecast)
			weatherGroup.GET("/calendar", cfg.WeatherHandler.GetCalendar)
			weatherGroup.GET("/compare", cfg.WeatherHandler.CompareCities)
		}

		// C. Location Management
		locationGroup := apiV1.Group("/locations")
		{
			locationGroup.GET("/search", optionalAuthMiddleware, cfg.LocationHandler.Search)
			locationGroup.POST("/saved", authMiddleware, cfg.LocationHandler.AddSaved)
			locationGroup.GET("/saved", authMiddleware, cfg.LocationHandler.GetSaved)
		}
	}

	return r
}
