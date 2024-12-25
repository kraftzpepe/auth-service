package config

import (
	"log"
	"os"
)

type Config struct {
	DatabaseURL string
	GRPCPort    string
}

func LoadConfig() *Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051" // Default gRPC port
	}

	return &Config{
		DatabaseURL: dbURL,
		GRPCPort:    grpcPort,
	}
}
