package main

import (
	"log"
	"os"
	"strconv"

	"github.com/ivan-ca97/life/internal/infrastructure/postgres"
	"github.com/ivan-ca97/life/internal/server"
)

var version = "dev"

func main() {
	dbConfig := postgres.ConnectionConfig{
		Host:     mustEnv("POSTGRES_HOST"),
		Port:     mustEnv("POSTGRES_PORT"),
		User:     mustEnv("POSTGRES_USER"),
		Password: mustEnv("POSTGRES_PASSWORD"),
		Database: mustEnv("POSTGRES_DB"),
	}

	database, err := postgres.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = postgres.RunMigrations(dbConfig)
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	portStr := mustEnv("PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("invalid PORT: %v", err)
	}

	corsOrigins := os.Getenv("CORS_ORIGINS")
	seedEmail := os.Getenv("SEED_ADMIN_EMAIL")
	seedPassword := os.Getenv("SEED_ADMIN_PASSWORD")
	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")

	r2AccountId := os.Getenv("R2_ACCOUNT_ID")
	r2AccessKeyId := os.Getenv("R2_ACCESS_KEY_ID")
	r2SecretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	r2Bucket := os.Getenv("R2_BUCKET")
	r2PublicURL := os.Getenv("R2_PUBLIC_URL")

	s, err := server.NewServer(database, port, version, corsOrigins, seedEmail, seedPassword, googleClientId, r2AccountId, r2AccessKeyId, r2SecretAccessKey, r2Bucket, r2PublicURL)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	err = s.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing required environment variable: %s", key)
	}
	return v
}
