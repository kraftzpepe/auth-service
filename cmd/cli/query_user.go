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

var queryUserCmd = &cobra.Command{
	Use:   "query-user",
	Short: "Query a user's details",
	Long:  "Query a user's details by email, UUID, or username.",
	Run: func(cmd *cobra.Command, args []string) {
		// Local variable declarations
		email, _ := cmd.Flags().GetString("email")
		uuid, _ := cmd.Flags().GetString("uuid")
		username, _ := cmd.Flags().GetString("username")

		// Connect to the gRPC server
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to server: %v", err)
		}
		defer conn.Close()

		client := pb.NewAuthServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var res *pb.GetUserResponse

		switch {
		case email != "":
			req := &pb.GetUserRequest{Identifier: email}
			res, err = client.GetUserByEmail(ctx, req)
		case uuid != "":
			req := &pb.GetUserRequest{Identifier: uuid}
			res, err = client.GetUserByUUID(ctx, req)
		case username != "":
			req := &pb.GetUserRequest{Identifier: username}
			res, err = client.GetUserByUsername(ctx, req)
		default:
			log.Fatalf("Please provide one of --email, --uuid, or --username")
		}

		if err != nil {
			log.Fatalf("Failed to query user: %v", err)
		}

		// Display the user details
		fmt.Printf("User Details:\n")
		fmt.Printf("ID: %s\n", res.GetId())
		fmt.Printf("Username: %s\n", res.GetUsername())
		fmt.Printf("Email: %s\n", res.GetEmail())
		fmt.Printf("Created At: %s\n", res.GetCreatedAt())
		fmt.Printf("Updated At: %s\n", res.GetUpdatedAt())
	},
}

func init() {
	// Add flags for querying user details
	queryUserCmd.Flags().String("email", "", "Query user by email")
	queryUserCmd.Flags().String("uuid", "", "Query user by UUID")
	queryUserCmd.Flags().String("username", "", "Query user by username")

	rootCmd.AddCommand(queryUserCmd)
}
