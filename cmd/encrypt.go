package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/YoungMayor/secureflow-go/internal/config"
	"github.com/YoungMayor/secureflow-go/internal/crypto"
	"github.com/YoungMayor/secureflow-go/internal/utils"
	"github.com/spf13/cobra"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt files specified in the configuration",
	Long: `Encrypts all files listed in secureflow.yaml using AES-256-CBC encryption.
Generates a report file with metadata about encrypted files.`,
	RunE: runEncrypt,
}

func init() {
	rootCmd.AddCommand(encryptCmd)
}

func runEncrypt(cmd *cobra.Command, args []string) error {
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
		pwd, err = utils.ReadPassword(utils.Colorize(utils.ColorBlue, "🔐 Enter password to encrypt your secrets: "))
		if err != nil {
			return err
		}
		if pwd == "" {
			return fmt.Errorf("password cannot be empty")
		}
	}

	// Get optional password hint
	var passwordHint string
	if !nonInteractive {
		passwordHint, err = utils.ReadLine(utils.Colorize(utils.ColorBlue, "🔑 (Optional) Enter a password hint (leave blank to skip): "))
		if err != nil {
			return err
		}
	}

	// Get optional note
	var note string
	if !nonInteractive {
		note, err = utils.ReadLine(utils.Colorize(utils.ColorBlue, "📝 (Optional) Enter a short note (leave blank for default): "))
		if err != nil {
			return err
		}
	}
	if note == "" {
		note = "Encrypted secrets for CI/CD"
	}

	fmt.Println()

	// Ensure output directory exists
	if err := utils.EnsureDir(cfg.OutputDir); err != nil {
		return err
	}

	// Create report file
	reportPath := filepath.Join(cfg.OutputDir, "report.txt")
	reportFile, err := os.Create(reportPath)
	if err != nil {
		return fmt.Errorf("failed to create report file: %w", err)
	}
	defer reportFile.Close()

	// Write report header
	fmt.Fprintf(reportFile, "Encryption Report\n")
	fmt.Fprintf(reportFile, "=================\n")
	fmt.Fprintf(reportFile, "\n")
	fmt.Fprintf(reportFile, "Note: %s\n", note)
	if passwordHint != "" {
		fmt.Fprintf(reportFile, "Password Hint: %s\n", passwordHint)
	} else {
		fmt.Fprintf(reportFile, "Password Hint: N/A\n")
	}
	fmt.Fprintf(reportFile, "Created at: %s\n", time.Now().Format("2006-01-02"))
	fmt.Fprintf(reportFile, "=================\n")
	fmt.Fprintf(reportFile, "\n")

	// Encrypt each file
	successCount := 0
	for _, fileMapping := range cfg.Files {
		fmt.Printf("%s 📦 Encrypting %s...\n", utils.ColorYellow, fileMapping.Input)

		// Check if input file exists
		if !utils.FileExists(fileMapping.Input) {
			fmt.Printf("%s ⚠️  Warning: %s not found, skipping\n\n", utils.ColorYellow, fileMapping.Input)
			continue
		}

		// Get file info before encryption
		fileInfo, err := utils.GetFileInfo(fileMapping.Input)
		if err != nil {
			fmt.Printf("%s ⚠️  Warning: Could not get file info for %s: %v\n\n", utils.ColorYellow, fileMapping.Input, err)
			continue
		}

		// Encrypt file
		outputPath := filepath.Join(cfg.OutputDir, fileMapping.Output)
		if err := crypto.EncryptFile(fileMapping.Input, outputPath, pwd); err != nil {
			fmt.Printf("%s ❌ Failed to encrypt %s: %v\n\n", utils.ColorRed, fileMapping.Input, err)
			continue
		}

		fmt.Printf("%s ✅ %s encrypted successfully -> %s\n\n", utils.ColorGreen, fileMapping.Input, outputPath)

		// Write to report
		fmt.Fprintf(reportFile, "File:           %s\n", fileMapping.Input)
		fmt.Fprintf(reportFile, "Encrypted As:   %s\n", fileMapping.Output)
		fmt.Fprintf(reportFile, "Size (bytes):   %d\n", fileInfo.Size)
		fmt.Fprintf(reportFile, "Lines:          %d\n", fileInfo.Lines)
		fmt.Fprintf(reportFile, "Last Modified:  %s\n", fileInfo.LastModified.Format("2006-01-02 15:04:05"))
		fmt.Fprintf(reportFile, "----------------------------------------\n")
		fmt.Fprintf(reportFile, "\n")

		successCount++
	}

	if successCount == 0 {
		return fmt.Errorf("no files were encrypted")
	}

	fmt.Printf("%s ✅ Encryption complete. %d file(s) saved to %s\n", utils.ColorGreen, successCount, cfg.OutputDir)
	fmt.Printf("📄 Report saved to %s\n", reportPath)

	return nil
}
