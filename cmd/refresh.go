package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/MayR-Labs/secureflow-go/internal/config"
	"github.com/MayR-Labs/secureflow-go/internal/crypto"
	"github.com/MayR-Labs/secureflow-go/internal/manifest"
	"github.com/MayR-Labs/secureflow-go/internal/utils"
	"github.com/spf13/cobra"
)

var refreshCmd = &cobra.Command{
	Use:     "refresh",
	Aliases: []string{"rotate"},
	Short:   "Re-encrypt all files, optionally with a new password",
	Long: `Decrypts every encrypted file with the current (old) password and
re-encrypts it with a new password, or with the same password if no new one is
given. Plaintext is only held in memory; nothing decrypted is written to disk.

The operation is all-or-nothing: if any file fails to decrypt, no files are
changed. Refreshing also upgrades files from the legacy AES-CBC format to
AES-256-GCM.

Passwords can be supplied with flags or environment variables:
  --old-password  (or --password, SECUREFLOW_PASSWORD)
  --new-password  (or SECUREFLOW_NEW_PASSWORD)

Examples:
  secureflow refresh                       # interactive
  secureflow rotate --non-interactive \
    --old-password "$OLD" --new-password "$NEW" --hint "Q1 2027 vault entry"`,
	RunE: runRefresh,
}

var (
	refreshOldPassword string
	refreshNewPassword string
	refreshHint        string
)

func init() {
	rootCmd.AddCommand(refreshCmd)
	refreshCmd.Flags().StringVar(&refreshOldPassword, "old-password", "", "current password (defaults to --password / "+PasswordEnvVar+")")
	refreshCmd.Flags().StringVar(&refreshNewPassword, "new-password", "", "new password; omit to keep the current one (or set "+NewPasswordEnvVar+")")
	refreshCmd.Flags().StringVar(&refreshHint, "hint", "", "new password hint saved in the manifest")
}

func runRefresh(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	m, err := manifest.Load(cfg.OutputDir)
	if err != nil {
		return err
	}

	// Old password
	oldPwd := refreshOldPassword
	if oldPwd == "" {
		oldPwd, err = promptForPassword("🔐 Enter the CURRENT password: ", m.PasswordHint)
		if err != nil {
			return err
		}
	}

	// Decrypt everything first; abort before touching anything if one fails
	fmt.Println()
	files, err := openAll(cfg, oldPwd, m.PasswordHint)
	if err != nil {
		return err
	}
	fmt.Println()

	// New password (optional)
	newPwd := refreshNewPassword
	if newPwd == "" {
		newPwd = os.Getenv(NewPasswordEnvVar)
	}
	if newPwd == "" && !nonInteractive {
		newPwd, err = utils.ReadPassword(utils.Colorize(utils.ColorBlue, "🔐 Enter a NEW password (leave blank to keep the current one): "))
		if err != nil {
			return err
		}
		if newPwd != "" {
			if err := confirmPassword(newPwd); err != nil {
				return err
			}
		}
	}
	changed := newPwd != "" && newPwd != oldPwd
	if newPwd == "" {
		newPwd = oldPwd
	}

	// Hint: keep it unless changed; prompt when the password changed
	if cmd.Flags().Changed("hint") {
		m.PasswordHint = refreshHint
	} else if changed && !nonInteractive {
		if m.PasswordHint, err = promptHint(m.PasswordHint); err != nil {
			return err
		}
	} else if changed && m.PasswordHint != "" {
		fmt.Printf("%s ⚠️  Password changed but the hint was kept; pass --hint to update it.%s\n", utils.ColorYellow, utils.ColorReset)
	}

	// Re-encrypt in memory, then swap all files in together
	var out []*sealed
	upgraded := 0
	for _, f := range files {
		srcModified := time.Time{}
		if prev, ok := m.Get(f.fm.Input); ok {
			srcModified = prev.SourceModified
		}
		s, err := seal(cfg, f.fm, f.plaintext, newPwd, srcModified)
		if err != nil {
			return fmt.Errorf("failed to re-encrypt %s: %w", f.fm.Input, err)
		}
		out = append(out, s)
		if f.format == crypto.FormatLegacy {
			upgraded++
		}
	}

	if err := writeSealed(out); err != nil {
		return err
	}
	for _, s := range out {
		m.Upsert(s.entry)
		fmt.Printf("%s 🔄 %s%s\n", utils.ColorGreen, s.path, utils.ColorReset)
	}
	m.Prune(cfg.Inputs())
	if err := m.Save(cfg.OutputDir); err != nil {
		return err
	}

	fmt.Println()
	if changed {
		fmt.Printf("%s ✅ Re-encrypted %d file(s) with the new password.%s\n", utils.ColorGreen, len(out), utils.ColorReset)
		fmt.Println("   Remember to update the password in your CI secrets (e.g. SECUREFLOW_PASSWORD) and tell your team.")
	} else {
		fmt.Printf("%s ✅ Re-encrypted %d file(s) with the same password.%s\n", utils.ColorGreen, len(out), utils.ColorReset)
	}
	if upgraded > 0 {
		fmt.Printf("   %d file(s) upgraded from legacy AES-CBC to AES-256-GCM.\n", upgraded)
	}
	return nil
}
