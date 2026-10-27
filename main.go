package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"go_26_10/config"
	"go_26_10/domain"
	httpDelivery "go_26_10/internal/delivery/http"
	"go_26_10/internal/repository"
	"go_26_10/internal/usecase"
)

func main() {
	// 1. Inisialisasi Konfigurasi Environment
	cfg := config.LoadConfig()
	log.Printf("[INFO] Memulai layanan go_vibe dalam mode: %s\n", cfg.AppEnv)

	// 2. Inisialisasi Koneksi Database PostgreSQL via GORM
	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Printf("[WARN] Koneksi database belum tersedia: %v\n", err)
		log.Println("[WARN] Layanan berjalan dalam mode fallback/in-memory. Endpoint tetap siap diuji.")
	} else {
		log.Println("[INFO] Berhasil terhubung ke database PostgreSQL")
		_ = db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error
		_ = db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`).Error

		// Auto Migrate tabel jika belum dibuat oleh migration
		if err := db.AutoMigrate(&domain.User{}, &domain.Location{}, &domain.UserSavedLocation{}); err != nil {
			log.Printf("[WARN] AutoMigrate warning: %v\n", err)
		} else {
			log.Println("[INFO] Skema database berhasil dimigrasikan secara otomatis")
		}
	}

	// 3. Inisialisasi Repository Layer (Database & External API)
	userRepo := repository.NewUserRepository(db)
	locationRepo := repository.NewLocationRepository(db)
	weatherRepo := repository.NewWeatherRepository(cfg.WeatherAPIKey, cfg.WeatherAPIBaseURL)

	// 4. Inisialisasi Usecase Layer (Business Logic & Analytics)
	userUsecase := usecase.NewUserUsecase(userRepo, cfg.JWTSecret)
	locationUsecase := usecase.NewLocationUsecase(locationRepo)
	weatherUsecase := usecase.NewWeatherUsecase(weatherRepo, locationRepo, userRepo)

	// 5. Inisialisasi Delivery / HTTP Handler Layer
	authHandler := httpDelivery.NewAuthHandler(userUsecase)
	userHandler := httpDelivery.NewUserHandler(userUsecase)
	locationHandler := httpDelivery.NewLocationHandler(locationUsecase)
	weatherHandler := httpDelivery.NewWeatherHandler(weatherUsecase)

	// 6. Setup Router & Routing Endpoints
	router := httpDelivery.SetupRouter(httpDelivery.RouterConfig{
		AuthHandler:     authHandler,
		UserHandler:     userHandler,
		LocationHandler: locationHandler,
		WeatherHandler:  weatherHandler,
		JWTSecret:       cfg.JWTSecret,
	})

	// 7. Konfigurasi HTTP Server dengan Graceful Shutdown
	port := strings.TrimPrefix(strings.TrimSpace(cfg.AppPort), ":")
	if port == "" {
		port = "8080"
	}
	serverAddr := fmt.Sprintf(":%s", port)
	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Jalankan server di background goroutine
	go func() {
		log.Printf("[INFO] Server HTTP aktif dan mendengarkan pada http://localhost%s\n", serverAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[FATAL] Server HTTP berhenti karena error: %v\n", err)
		}
	}()

	// Menunggu sinyal interupsi OS (Ctrl+C atau SIGTERM) untuk graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("[INFO] Menerima sinyal mematikan server...")

	// Waktu batas shutdown 5 detik
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[FATAL] Server dipaksa berhenti: %v\n", err)
	}

	log.Println("[INFO] Server go_vibe berhasil dihentikan dengan aman.")
}
