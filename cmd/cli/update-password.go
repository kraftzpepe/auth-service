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

var updatePasswordCmd = &cobra.Command{
	Use:   "update-password",
	Short: "Update the user's password",
	Long:  "Update the user's password by providing the old password and a new password.",
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags for email, oldPassword, and newPassword
		email, _ := cmd.Flags().GetString("email")
		oldPassword, _ := cmd.Flags().GetString("old-password")
		newPassword, _ := cmd.Flags().GetString("new-password")

		// Connect to the gRPC server
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to server: %v", err)
		}
		defer conn.Close()

		client := pb.NewAuthServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Authenticate user
		loginReq := &pb.LoginRequest{
			Email:    email,
			Password: oldPassword,
		}

		_, err = client.Login(ctx, loginReq)
		if err != nil {
			log.Fatalf("Authentication failed: %v", err)
		}

		// Reset password
		resetReq := &pb.ResetPasswordRequest{
			Token:       oldPassword, // Assuming you use the old password as a "token"
			NewPassword: newPassword,
		}

		_, err = client.ResetPassword(ctx, resetReq)
		if err != nil {
			log.Fatalf("Failed to update password: %v", err)
		}

		fmt.Println("Password updated successfully.")
	},
}

func init() {
	// Add flags for the update-password command
	updatePasswordCmd.Flags().String("email", "", "Email for the user")
	updatePasswordCmd.Flags().String("old-password", "", "User's current password")
	updatePasswordCmd.Flags().String("new-password", "", "User's new password")
	updatePasswordCmd.MarkFlagRequired("email")
	updatePasswordCmd.MarkFlagRequired("old-password")
	updatePasswordCmd.MarkFlagRequired("new-password")

	rootCmd.AddCommand(updatePasswordCmd)
}
