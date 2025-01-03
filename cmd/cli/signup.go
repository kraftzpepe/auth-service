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

var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Sign up a new user",
	Long:  "Sign up a new user by providing a username, email, and password.",
	Run: func(cmd *cobra.Command, args []string) {
		// Local variable declarations
		username, _ := cmd.Flags().GetString("username")
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")

		// Connect to the gRPC server
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to server: %v", err)
		}
		defer conn.Close()

		client := pb.NewAuthServiceClient(conn)

		// Create the RegisterRequest
		req := &pb.RegisterRequest{
			Username: username,
			Email:    email,
			Password: password,
		}

		// Set a timeout for the request
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Call the Register method
		res, err := client.Register(ctx, req)
		if err != nil {
			log.Fatalf("Failed to register user: %v", err)
		}

		// Print the response
		fmt.Printf("User registered successfully:\n")
		fmt.Printf("ID: %s\nUsername: %s\nEmail: %s\nAccessToken: %s\nRefreshToken: %s\n",
			res.GetId(), res.GetUsername(), res.GetEmail(), res.GetAccessToken(), res.GetRefreshToken())
	},
}

func init() {
	// Add flags for the signup command
	signupCmd.Flags().String("username", "", "Username for the new user")
	signupCmd.Flags().String("email", "", "Email for the new user")
	signupCmd.Flags().String("password", "", "Password for the new user")
	signupCmd.MarkFlagRequired("username")
	signupCmd.MarkFlagRequired("email")
	signupCmd.MarkFlagRequired("password")

	rootCmd.AddCommand(signupCmd)
}
