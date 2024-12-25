package cli

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "auth-cli",
	Short: "A CLI for interacting with the AuthService",
	Long:  "A Command Line Interface for managing user authentication with the AuthService.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}