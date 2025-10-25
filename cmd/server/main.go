package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	db "github.com/org/ranas-bdi-backend/internal/adapters/db/postgres"
	"github.com/org/ranas-bdi-backend/internal/adapters/http"
)

func loadConfig() (db.Config, error) {
	requiredVars := []string{
		"PG_HOST",
		"PG_PORT",
		"PG_DB",
		"PG_USER",
		"PG_PASSWORD",
	}

	for _, envVar := range requiredVars {
		if os.Getenv(envVar) == "" {
			return db.Config{}, fmt.Errorf("variable de entorno requerida no encontrada: %s", envVar)
		}
	}

	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")
	dbname := os.Getenv("PG_DB")
	sslmode := getEnvOrDefault("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode)

	maxConns := getEnvAsIntOrDefault("DB_MAX_CONNS", 10)
	minConns := getEnvAsIntOrDefault("DB_MIN_CONNS", 0)
	maxConnIdleTime := getEnvAsDurationOrDefault("DB_MAX_CONN_IDLE_TIME", "5m")
	appName := getEnvOrDefault("APP_NAME", "ranas-bdi-backend")

	return db.Config{
		DSN:             dsn,
		MaxConns:        int32(maxConns),
		MinConns:        int32(minConns),
		MaxConnIdleTime: maxConnIdleTime,
		AppName:         appName,
	}, nil
}

// getEnvOrDefault obtiene una variable de entorno o retorna un valor por defecto
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsIntOrDefault obtiene una variable de entorno como entero o retorna un valor por defecto
func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsDurationOrDefault obtiene una variable de entorno como duración o retorna un valor por defecto
func getEnvAsDurationOrDefault(key, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}

func main() {
	_ = godotenv.Load()

	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}

	if err := db.Init(context.Background(), cfg); err != nil {
		log.Fatalf("Error inicializando base de datos: %v", err)
	}
	defer db.Close()
	log.Println("✅ Connected to Postgres!")

	handler := http.NewHandler()
	router := handler.GetRouter()

	port := getEnvOrDefault("PORT", "8080")
	log.Printf("Starting server on :%s", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
