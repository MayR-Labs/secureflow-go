package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/MayR-Labs/secureflow-go/internal/config"
	"github.com/MayR-Labs/secureflow-go/internal/crypto"
	"github.com/MayR-Labs/secureflow-go/internal/utils"
	"github.com/spf13/cobra"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt files specified in the configuration",
	Long: `Decrypts all encrypted files listed in secureflow.yaml back to their 
original locations. Useful for local development and CI/CD pipelines.`,
	RunE: runDecrypt,
}

func init() {
	rootCmd.AddCommand(decryptCmd)
}

func runDecrypt(cmd *cobra.Command, args []string) error {
	// Load config
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Get password
	var pwd string
	if password != "" {
		pwd = password
	} else if nonInteractive {
		return fmt.Errorf("password required in non-interactive mode (use --password flag)")
	} else {
		pwd, err = utils.ReadPassword(utils.Colorize(utils.ColorBlue, "🔐 Enter password to decrypt your secrets: "))
		if err != nil {
			return err
		}
		if pwd == "" {
			return fmt.Errorf("password cannot be empty")
		}
	}

	fmt.Println()
	fmt.Printf("%s 🔐 Starting decryption process...\n\n", utils.ColorYellow)

	// Decrypt each file
	successCount := 0
	for _, fileMapping := range cfg.Files {
		encryptedPath := filepath.Join(cfg.OutputDir, fileMapping.Output)

		fmt.Printf("%s 📄 Decrypting %s...\n", utils.ColorYellow, encryptedPath)

		// Check if encrypted file exists
		if !utils.FileExists(encryptedPath) {
			fmt.Printf("%s ⚠️  Warning: %s not found, skipping\n\n", utils.ColorYellow, encryptedPath)
			continue
		}

		// Ensure output directory exists
		outputDir := filepath.Dir(fileMapping.Input)
		if outputDir != "." && outputDir != "" {
			if err := utils.EnsureDir(outputDir); err != nil {
				fmt.Printf("%s ❌ Failed to create directory %s: %v\n\n", utils.ColorRed, outputDir, err)
				continue
			}
		}

		// Decrypt file
		if err := crypto.DecryptFile(encryptedPath, fileMapping.Input, pwd); err != nil {
			fmt.Printf("%s ❌ Failed to decrypt %s: %v\n\n", utils.ColorRed, encryptedPath, err)
			return fmt.Errorf("decryption failed (wrong password?)")
		}

		fmt.Printf("%s ✅ %s decrypted successfully -> %s\n", utils.ColorGreen, encryptedPath, fileMapping.Input)

		// Handle copy_to if specified
		if fileMapping.CopyTo != "" {
			// Ensure copy_to directory exists
			copyToDir := filepath.Dir(fileMapping.CopyTo)
			if copyToDir != "." && copyToDir != "" {
				if err := utils.EnsureDir(copyToDir); err != nil {
					fmt.Printf("%s ⚠️  Warning: Failed to create directory %s: %v\n", utils.ColorYellow, copyToDir, err)
					fmt.Println()
					continue
				}
			}

			// Copy the decrypted file to the copy_to location
			if err := utils.CopyFile(fileMapping.Input, fileMapping.CopyTo); err != nil {
				fmt.Printf("%s ⚠️  Warning: Failed to copy %s to %s: %v\n", utils.ColorYellow, fileMapping.Input, fileMapping.CopyTo, err)
			} else {
				fmt.Printf("%s 📋 Copied to %s\n", utils.ColorGreen, fileMapping.CopyTo)
			}
		}

		fmt.Println()
		successCount++
	}

	if successCount == 0 {
		return fmt.Errorf("no files were decrypted")
	}

	fmt.Println()
	fmt.Printf("%s 🎉 All secrets decrypted successfully! (%d file(s))\n", utils.ColorGreen, successCount)

	return nil
}
