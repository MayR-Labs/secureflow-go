package cmd

import (
	"fmt"

	"github.com/MayR-Labs/secureflow-go/internal/config"
	"github.com/MayR-Labs/secureflow-go/internal/crypto"
	"github.com/MayR-Labs/secureflow-go/internal/manifest"
	"github.com/MayR-Labs/secureflow-go/internal/utils"
	"github.com/spf13/cobra"
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt files specified in the configuration",
	Long: `Decrypts all encrypted files listed in secureflow.yaml back to their
original locations. Useful for local development and CI/CD pipelines.

Decrypted files are written with owner-only permissions (0600).`,
	RunE: runDecrypt,
}

func init() {
	rootCmd.AddCommand(decryptCmd)
}

// opened is a decrypted file held in memory
type opened struct {
	fm        config.FileMapping
	plaintext []byte
	format    crypto.Format
}

// openAll decrypts every available encrypted file in memory. Nothing is
// written to disk, so a wrong password leaves the working tree untouched.
func openAll(cfg *config.Config, pwd, hint string) ([]opened, error) {
	var result []opened
	for _, fm := range cfg.Files {
		path := encryptedPath(cfg, fm)
		if !utils.FileExists(path) {
			fmt.Printf("%s ⚠️  Warning: %s not found, skipping%s\n", utils.ColorYellow, path, utils.ColorReset)
			continue
		}
		fmt.Printf("%s 📄 Decrypting %s...%s\n", utils.ColorYellow, path, utils.ColorReset)

		plaintext, format, err := openEncrypted(cfg, fm, pwd)
		if err != nil {
			fmt.Printf("%s ❌ Failed to decrypt %s: %v%s\n\n", utils.ColorRed, path, err, utils.ColorReset)
			if isWrongPassword(err) {
				printHint(hint)
				return nil, fmt.Errorf("decryption failed (wrong password?)")
			}
			return nil, err
		}
		result = append(result, opened{fm: fm, plaintext: plaintext, format: format})
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no encrypted files found in %s", cfg.OutputDir)
	}
	return result, nil
}

// warnLegacy suggests migrating files still in the v1 (AES-CBC) format
func warnLegacy(files []opened) {
	legacy := 0
	for _, f := range files {
		if f.format == crypto.FormatLegacy {
			legacy++
		}
	}
	if legacy > 0 {
		fmt.Printf("\n%s ℹ️  %d file(s) use the legacy AES-CBC format. Run 'secureflow refresh' to upgrade them to AES-GCM.%s\n",
			utils.ColorBlue, legacy, utils.ColorReset)
	}
}

func runDecrypt(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Show the hint saved at encryption time to help users remember the password
	hint := manifest.ReadHint(cfg.OutputDir)

	pwd, err := promptForPassword("🔐 Enter password to decrypt your secrets: ", hint)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf("%s 🔐 Starting decryption process...%s\n\n", utils.ColorYellow, utils.ColorReset)

	files, err := openAll(cfg, pwd, hint)
	if err != nil {
		return err
	}

	fmt.Println()
	for _, f := range files {
		if err := writePlaintext(f.fm.Input, f.plaintext); err != nil {
			return fmt.Errorf("failed to write %s: %w", f.fm.Input, err)
		}
		fmt.Printf("%s ✅ %s%s\n", utils.ColorGreen, f.fm.Input, utils.ColorReset)

		if f.fm.CopyTo != "" {
			if err := writePlaintext(f.fm.CopyTo, f.plaintext); err != nil {
				fmt.Printf("%s ⚠️  Warning: failed to copy to %s: %v%s\n", utils.ColorYellow, f.fm.CopyTo, err, utils.ColorReset)
			} else {
				fmt.Printf("%s 📋 Copied to %s%s\n", utils.ColorGreen, f.fm.CopyTo, utils.ColorReset)
			}
		}
	}

	warnLegacy(files)

	fmt.Println()
	fmt.Printf("%s 🎉 All secrets decrypted successfully! (%d file(s))%s\n", utils.ColorGreen, len(files), utils.ColorReset)
	return nil
}
