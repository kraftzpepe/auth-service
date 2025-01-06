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

var resetPasswordCmd = &cobra.Command{
	Use:   "reset-password",
	Short: "Reset a user's password using a reset token",
	Long:  "Reset a user's password using a valid password reset token and a new password.",
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags for token and new password
		token, _ := cmd.Flags().GetString("token")
		newPassword, _ := cmd.Flags().GetString("new-password")

		// Connect to the gRPC server
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to server: %v", err)
		}
		defer conn.Close()

		client := pb.NewAuthServiceClient(conn)

		// Create the ResetPasswordRequest
		req := &pb.ResetPasswordRequest{
			Token:       token,
			NewPassword: newPassword,
		}

		// Set a timeout for the request
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Call the ResetPassword method
		res, err := client.ResetPassword(ctx, req)
		if err != nil {
			log.Fatalf("Failed to reset password: %v", err)
		}

		// Print the confirmation message
		fmt.Printf("Password reset successful:\n")
		fmt.Printf("%s\n", res.GetMessage())
	},
}

func init() {
	// Add flags for token and new password
	resetPasswordCmd.Flags().String("token", "", "Password reset token")
	resetPasswordCmd.Flags().String("new-password", "", "New password for the user")
	resetPasswordCmd.MarkFlagRequired("token")
	resetPasswordCmd.MarkFlagRequired("new-password")

	// Add the command to the root command
	rootCmd.AddCommand(resetPasswordCmd)
}
