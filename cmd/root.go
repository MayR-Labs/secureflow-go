package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version information
	Version = "1.0.0"

	// Global flags
	cfgFile        string
	nonInteractive bool
	password       string
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "secureflow",
	Short: "SecureFlow - Secure file encryption/decryption CLI",
	Long: `SecureFlow is a lightweight, Go-based CLI for securely encrypting 
and decrypting sensitive files like environment variables, keystores, 
and service credentials for local and CI/CD use.`,
	Version: Version,
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "secureflow.yaml", "config file path")
	rootCmd.PersistentFlags().BoolVar(&nonInteractive, "non-interactive", false, "run in non-interactive mode")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "encryption/decryption password (for non-interactive mode)")
}
