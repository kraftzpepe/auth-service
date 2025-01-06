package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/kraftzpepe/auth-service/config"
	"github.com/kraftzpepe/auth-service/db"
	"github.com/kraftzpepe/auth-service/internal/handler"
	"github.com/kraftzpepe/auth-service/internal/repositories"
	"github.com/kraftzpepe/auth-service/internal/service"

	pb "github.com/kraftzpepe/auth-service/proto/generated"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load database connection
	database, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(database)
	refreshTokenRepo := repositories.NewRefreshTokenRepository(database)
	passwordResetTokenRepo := repositories.NewPasswordResetTokenRepository(database)

	// Initialize services
	authService := service.NewAuthService(userRepo, refreshTokenRepo, passwordResetTokenRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)

	// Start gRPC server
	grpcPort := config.LoadConfig().GRPCPort
	listener, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authHandler) // Ensure the AuthServiceServer is registered

	// Graceful shutdown setup
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting gRPC server on port %s", grpcPort)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	<-stopChan // Wait for termination signal
	log.Println("Shutting down gracefully...")

	grpcServer.GracefulStop()
	log.Println("gRPC server stopped")
	log.Println("Service shutdown complete")
}
