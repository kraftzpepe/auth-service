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

var requestPasswordResetCmd = &cobra.Command{
	Use:   "request-password-reset",
	Short: "Request a password reset for a user",
	Long:  "Request a password reset by sending a reset token to the user's email address.",
	Run: func(cmd *cobra.Command, args []string) {
		// Get email flag
		email, _ := cmd.Flags().GetString("email")

		// Connect to the gRPC server
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to server: %v", err)
		}
		defer conn.Close()

		client := pb.NewAuthServiceClient(conn)

		// Create the RequestPasswordResetRequest
		req := &pb.RequestPasswordResetRequest{
			Email: email,
		}

		// Set a timeout for the request
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Call the RequestPasswordReset method
		res, err := client.RequestPasswordReset(ctx, req)
		if err != nil {
			log.Fatalf("Failed to request password reset: %v", err)
		}

		// Print the confirmation message
		fmt.Printf("Password reset request successful:\n")
		fmt.Printf("%s\n", res.GetMessage())
	},
}

func init() {
	// Add flag for the email
	requestPasswordResetCmd.Flags().String("email", "", "Email address of the user requesting password reset")
	requestPasswordResetCmd.MarkFlagRequired("email")

	// Add the command to the root command
	rootCmd.AddCommand(requestPasswordResetCmd)
}
