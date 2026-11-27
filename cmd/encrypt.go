package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/MayR-Labs/secureflow-go/internal/config"
	"github.com/MayR-Labs/secureflow-go/internal/manifest"
	"github.com/MayR-Labs/secureflow-go/internal/utils"
	"github.com/spf13/cobra"
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt files specified in the configuration",
	Long: `Encrypts all files listed in secureflow.yaml using AES-256-GCM with a
PBKDF2-SHA256 derived key (600,000 iterations).

Writes <output_dir>/manifest.json with the password hint, note and metadata
about each encrypted file.

If encrypted files already exist, the password must match them. To change the
password, use 'secureflow refresh' (or --force to re-encrypt from the plaintext
files with a new password).`,
	RunE: runEncrypt,
}

var (
	encryptHint  string
	encryptNote  string
	encryptForce bool
)

func init() {
	rootCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().StringVar(&encryptHint, "hint", "", "password hint saved in the manifest and shown on decrypt")
	encryptCmd.Flags().StringVar(&encryptNote, "note", "", "short note saved in the manifest")
	encryptCmd.Flags().BoolVar(&encryptForce, "force", false, "re-encrypt with a password that differs from the existing encrypted files")
}

func runEncrypt(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	m, err := manifest.Load(cfg.OutputDir)
	if err != nil {
		return err
	}

	// Get password (confirm it interactively so a typo doesn't lock secrets away)
	pwd, err := promptForPassword("🔐 Enter password to encrypt your secrets: ", m.PasswordHint)
	if err != nil {
		return err
	}
	if passwordFromFlagOrEnv() == "" {
		if err := confirmPassword(pwd); err != nil {
			return err
		}
	}

	// Work out which inputs are present
	present := map[string]bool{}
	var missing []config.FileMapping
	for _, fm := range cfg.Files {
		if utils.FileExists(fm.Input) {
			present[fm.Input] = true
		} else {
			missing = append(missing, fm)
		}
	}
	if len(present) == 0 {
		return fmt.Errorf("none of the input files exist; nothing to encrypt")
	}

	// Make sure we don't end up with files encrypted under different passwords
	if err := checkPasswordConsistency(cfg, pwd, present, len(missing) > 0); err != nil {
		return err
	}

	// Hint and note
	if cmd.Flags().Changed("hint") {
		m.PasswordHint = encryptHint
	} else if !nonInteractive {
		if m.PasswordHint, err = promptHint(m.PasswordHint); err != nil {
			return err
		}
	}
	if cmd.Flags().Changed("note") {
		m.Note = encryptNote
	} else if !nonInteractive {
		def := m.Note
		if def == "" {
			def = manifest.DefaultNote
		}
		note, err := utils.ReadLine(utils.Colorize(utils.ColorBlue, fmt.Sprintf("📝 (Optional) Enter a short note (leave blank for %q): ", def)))
		if err != nil {
			return err
		}
		if note != "" {
			m.Note = note
		}
	}

	fmt.Println()

	// Encrypt everything in memory first, then write all files together
	var out []*sealed
	for _, fm := range cfg.Files {
		if !present[fm.Input] {
			fmt.Printf("%s ⚠️  Warning: %s not found, skipping%s\n", utils.ColorYellow, fm.Input, utils.ColorReset)
			continue
		}
		fmt.Printf("%s 📦 Encrypting %s...%s\n", utils.ColorYellow, fm.Input, utils.ColorReset)
		s, err := sealInput(cfg, fm, pwd)
		if err != nil {
			return fmt.Errorf("failed to encrypt %s: %w", fm.Input, err)
		}
		out = append(out, s)
	}

	if err := writeSealed(out); err != nil {
		return err
	}
	for _, s := range out {
		m.Upsert(s.entry)
		fmt.Printf("%s ✅ %s -> %s%s\n", utils.ColorGreen, s.entry.Input, s.path, utils.ColorReset)
	}

	m.Prune(cfg.Inputs())
	if err := m.Save(cfg.OutputDir); err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf("%s ✅ Encryption complete. %d file(s) saved to %s%s\n", utils.ColorGreen, len(out), cfg.OutputDir, utils.ColorReset)
	fmt.Printf("📄 Manifest saved to %s\n", manifest.Path(cfg.OutputDir))
	return nil
}

// checkPasswordConsistency stops encrypt from silently mixing passwords.
// If any existing encrypted file can't be opened with pwd:
//   - and some inputs are missing, those old files would keep the old password -> error
//   - otherwise every file is re-encrypted, so allow it with --force or confirmation
func checkPasswordConsistency(cfg *config.Config, pwd string, present map[string]bool, someMissing bool) error {
	err := verifyPassword(cfg, pwd, nil)
	if err == nil {
		return nil
	}
	if !errors.Is(err, errPasswordMismatch) {
		return err
	}

	if someMissing {
		return fmt.Errorf("%w\n(some input files are missing, so re-encrypting would leave files under two different passwords)", err)
	}
	if encryptForce {
		return nil
	}
	if nonInteractive {
		return fmt.Errorf("%w\n(pass --force to re-encrypt everything with the new password)", err)
	}

	fmt.Fprintln(os.Stderr, utils.Colorize(utils.ColorYellow, "⚠️  This password doesn't match your existing encrypted files."))
	ok, err := utils.Confirm("Re-encrypt every file with the new password? (y/N): ")
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("aborted; nothing was changed")
	}
	return nil
}
