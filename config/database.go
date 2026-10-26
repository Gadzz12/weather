package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds all application configuration values
type Config struct {
	AppPort            string
	AppEnv             string
	DatabaseURL        string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	DBSSLMode          string
	WeatherAPIKey      string
	WeatherAPIBaseURL  string
	JWTSecret          string
	JWTExpirationHours string
}

func cleanEnv(val string) string {
	val = strings.TrimSpace(val)
	val = strings.TrimPrefix(val, "'")
	val = strings.TrimSuffix(val, "'")
	val = strings.TrimPrefix(val, "\"")
	val = strings.TrimSuffix(val, "\"")
	return strings.TrimSpace(val)
}

// LoadConfig reads configuration from environment variables or .env file
func LoadConfig() *Config {
	loaded := false
	for _, envPath := range []string{".env", "../.env", "../../.env"} {
		if err := godotenv.Load(envPath); err == nil {
			loaded = true
			break
		}
	}
	if !loaded {
		log.Println("[INFO] No .env file found or error loading, using system environment variables")
	}

	rawBaseURL := cleanEnv(getEnv("WEATHER_API_BASE_URL", "https://api.weatherapi.com/v1"))
	// Normalisasi base URL jika pengguna memasukkan endpoint lengkap seperti /current.json?...
	if idx := strings.Index(rawBaseURL, "/v1"); idx != -1 {
		rawBaseURL = rawBaseURL[:idx+3]
	}

	return &Config{
		AppPort:            cleanEnv(getEnv("PORT", "")),
		AppEnv:             cleanEnv(getEnv("APP_ENV", "development")),
		DatabaseURL:        cleanEnv(getEnv("DATABASE_URL", "")),
		DBHost:             cleanEnv(getEnv("DB_HOST", "localhost")),
		DBPort:             cleanEnv(getEnv("DB_PORT", "5432")),
		DBUser:             cleanEnv(getEnv("DB_USER", "postgres")),
		DBPassword:         cleanEnv(getEnv("DB_PASSWORD", "secret")),
		DBName:             cleanEnv(getEnv("DB_NAME", "go_vibe_weather_db")),
		DBSSLMode:          cleanEnv(getEnv("DB_SSLMODE", "disable")),
		WeatherAPIKey:      cleanEnv(getEnv("WEATHER_API_KEY", "")),
		WeatherAPIBaseURL:  rawBaseURL,
		JWTSecret:          cleanEnv(getEnv("JWT_SECRET", "default_secret_key")),
		JWTExpirationHours: cleanEnv(getEnv("JWT_EXPIRATION_HOURS", "24")),
	}
}

// InitDatabase initializes PostgreSQL connection using GORM
func InitDatabase(cfg *Config) (*gorm.DB, error) {
	var dsn string

	if cfg.DatabaseURL != "" {
		dsn = cfg.DatabaseURL
		// Hilangkan parameter channel_binding yang terkadang tidak kompatibel dengan driver pgx tertentu
		dsn = strings.Replace(dsn, "&channel_binding=require", "", -1)
		dsn = strings.Replace(dsn, "?channel_binding=require", "?", -1)
	} else {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
			cfg.DBHost,
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBName,
			cfg.DBPort,
			cfg.DBSSLMode,
		)
	}

	logMode := logger.Warn
	if strings.EqualFold(cfg.AppEnv, "development") {
		logMode = logger.Info
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("gagal terhubung ke database PostgreSQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil instance koneksi DB: %w", err)
	}

	// Uji konektivitas database
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database unreachable (ping failed): %w", err)
	}

	return db, nil
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}
