package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/YoungMayor/secureflow-go/internal/config"
	"github.com/YoungMayor/secureflow-go/internal/crypto"
	"github.com/YoungMayor/secureflow-go/internal/utils"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test decryption without overwriting existing files",
	Long: `Decrypts files into a separate test directory to verify the encryption 
password is correct without overwriting existing secrets.`,
	RunE: runTest,
}

func init() {
	rootCmd.AddCommand(testCmd)
}

func runTest(cmd *cobra.Command, args []string) error {
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
		pwd, err = utils.ReadPassword(utils.Colorize(utils.ColorBlue, "🔐 [TEST] Enter password to test decrypt your secrets: "))
		if err != nil {
			return err
		}
		if pwd == "" {
			return fmt.Errorf("password cannot be empty")
		}
	}

	// Ensure test output directory exists
	if err := utils.EnsureDir(cfg.TestOutputDir); err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf("%s 🔐 [TEST] Starting decryption process...\n\n", utils.ColorYellow)

	// Decrypt each file to test directory
	successCount := 0
	for _, fileMapping := range cfg.Files {
		encryptedPath := filepath.Join(cfg.OutputDir, fileMapping.Output)
		
		fmt.Printf("%s 📄 Decrypting %s...\n", utils.ColorYellow, encryptedPath)

		// Check if encrypted file exists
		if !utils.FileExists(encryptedPath) {
			fmt.Printf("%s ⚠️  Warning: %s not found, skipping\n\n", utils.ColorYellow, encryptedPath)
			continue
		}

		// Get just the filename for test output
		testOutputPath := filepath.Join(cfg.TestOutputDir, filepath.Base(fileMapping.Input))

		// Decrypt file
		if err := crypto.DecryptFile(encryptedPath, testOutputPath, pwd); err != nil {
			fmt.Printf("%s ❌ Failed to decrypt %s: %v\n\n", utils.ColorRed, encryptedPath, err)
			return fmt.Errorf("test decryption failed (wrong password?)")
		}

		fmt.Printf("%s ✅ %s decrypted successfully -> %s\n\n", utils.ColorGreen, encryptedPath, testOutputPath)
		successCount++
	}

	if successCount == 0 {
		return fmt.Errorf("no files were decrypted")
	}

	fmt.Println()
	fmt.Printf("%s 🎉 Test decryption successful! (%d file(s))\n", utils.ColorGreen, successCount)
	fmt.Printf("📁 Test files saved to: %s\n", cfg.TestOutputDir)

	return nil
}
