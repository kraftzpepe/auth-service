package cli

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/kraftzpepe/auth-service/proto/generated"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in as a user",
	Long:  "Log in as a user by providing an email and password to receive access and refresh tokens.",
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags for email and password
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")

		// Connect to the gRPC server
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to server: %v", err)
		}
		defer conn.Close()

		client := pb.NewAuthServiceClient(conn)

		// Create the LoginRequest
		req := &pb.LoginRequest{
			Email:    email,
			Password: password,
		}

		// Set a timeout for the request
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Call the Login method
		res, err := client.Login(ctx, req)
		if err != nil {
			log.Fatalf("Login failed: %v", err)
		}

		// Print the tokens
		fmt.Printf("Login successful:\n")
		fmt.Printf("AccessToken: %s\n", res.GetAccessToken())
		fmt.Printf("RefreshToken: %s\n", res.GetRefreshToken())
	},
}

func init() {
	// Add flags for the login command
	loginCmd.Flags().String("email", "", "Email for the user")
	loginCmd.Flags().String("password", "", "Password for the user")
	loginCmd.MarkFlagRequired("email")
	loginCmd.MarkFlagRequired("password")

	rootCmd.AddCommand(loginCmd)
}
