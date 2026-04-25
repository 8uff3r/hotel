package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr              string
	DBPath            string
	SessionCookie     string
	HotelCookie       string
	SessionTTL        time.Duration
	RequestTimeout    time.Duration
	ShutdownTimeout   time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	SeedHotelCodeName string
	SeedHotelName     string
	SeedHotelAddress  string
	SeedHotelPhone    string
	SeedHotelEmail    string
	SeedAdminEmail    string
	SeedAdminPass     string
	SeedAdminFName    string
	SeedAdminLName    string
}

func Load() (Config, error) {
	cfg := Config{
		Addr:              getEnv("APP_ADDR", ":8080"),
		DBPath:            getEnv("DB_PATH", "postgres://hotel_user:hotel_password@localhost:5432/hotel_db?sslmode=disable"),
		SessionCookie:     getEnv("SESSION_COOKIE", "auth_session"),
		HotelCookie:       getEnv("HOTEL_COOKIE", "hotel_id"),
		SessionTTL:        getDuration("SESSION_TTL_HOURS", 24*time.Hour),
		RequestTimeout:    getDuration("REQUEST_TIMEOUT_SECONDS", 15*time.Second),
		ShutdownTimeout:   getDuration("SHUTDOWN_TIMEOUT_SECONDS", 20*time.Second),
		ReadTimeout:       getDuration("READ_TIMEOUT_SECONDS", 15*time.Second),
		WriteTimeout:      getDuration("WRITE_TIMEOUT_SECONDS", 15*time.Second),
		IdleTimeout:       getDuration("IDLE_TIMEOUT_SECONDS", 60*time.Second),
		SeedHotelCodeName: getEnv("SEED_HOTEL_CODE_NAME", "default"),
		SeedHotelName:     getEnv("SEED_HOTEL_NAME", "Default Hotel"),
		SeedHotelAddress:  getEnv("SEED_HOTEL_ADDRESS", "123 Main St"),
		SeedHotelPhone:    getEnv("SEED_HOTEL_PHONE", ""),
		SeedHotelEmail:    getEnv("SEED_HOTEL_EMAIL", ""),
		SeedAdminEmail:    getEnv("SEED_ADMIN_EMAIL", "admin@hotel.local"),
		SeedAdminPass:     getEnv("SEED_ADMIN_PASSWORD", "admin123"),
		SeedAdminFName:    getEnv("SEED_ADMIN_FIRST_NAME", "System"),
		SeedAdminLName:    getEnv("SEED_ADMIN_LAST_NAME", "Admin"),
	}

	if cfg.SeedAdminPass == "" {
		return Config{}, fmt.Errorf("SEED_ADMIN_PASSWORD cannot be empty")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	v, err := strconv.Atoi(raw)
	if err != nil || v <= 0 {
		return fallback
	}
	if fallback%time.Hour == 0 {
		return time.Duration(v) * time.Hour
	}
	return time.Duration(v) * time.Second
}
