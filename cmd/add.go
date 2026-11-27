package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MayR-Labs/secureflow-go/internal/config"
	"github.com/MayR-Labs/secureflow-go/internal/gitutil"
	"github.com/MayR-Labs/secureflow-go/internal/manifest"
	"github.com/MayR-Labs/secureflow-go/internal/utils"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <file>...",
	Short: "Add files to secureflow.yaml, .gitignore them, and encrypt them",
	Long: `Adds one or more files to secureflow.yaml, appends them to .gitignore so the
plaintext is never committed, and encrypts them with your existing password.

The password is checked against the files that are already encrypted, so
everything stays under a single password.

Examples:
  secureflow add .env.prod --copy-to .env
  secureflow add android/app/keystore.jks android/key.properties
  secureflow add ios/GoogleService-Info.plist --no-encrypt`,
	Args: cobra.MinimumNArgs(1),
	RunE: runAdd,
}

var (
	addOutput      string
	addCopyTo      string
	addNoGitignore bool
	addNoEncrypt   bool
)

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVar(&addOutput, "output", "", "encrypted filename (default: inferred from the path; single file only)")
	addCmd.Flags().StringVar(&addCopyTo, "copy-to", "", "copy the decrypted file here on decrypt (single file only)")
	addCmd.Flags().BoolVar(&addNoGitignore, "no-gitignore", false, "don't add the files to .gitignore")
	addCmd.Flags().BoolVar(&addNoEncrypt, "no-encrypt", false, "only update the config; don't encrypt now")
}

func runAdd(cmd *cobra.Command, args []string) error {
	if len(args) > 1 && (addOutput != "" || addCopyTo != "") {
		return fmt.Errorf("--output and --copy-to can only be used when adding a single file")
	}

	if !utils.FileExists(cfgFile) {
		return fmt.Errorf("%s not found; run 'secureflow init' first", cfgFile)
	}
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Build the new entries
	existing := map[string]bool{}
	for _, in := range cfg.Inputs() {
		existing[filepath.ToSlash(filepath.Clean(in))] = true
	}
	var entries []config.FileMapping
	for _, arg := range args {
		in := filepath.ToSlash(filepath.Clean(arg))
		if existing[in] {
			fmt.Printf("%s ℹ️  %s is already in %s, skipping%s\n", utils.ColorBlue, in, cfgFile, utils.ColorReset)
			continue
		}
		if strings.HasPrefix(in, "../") || filepath.IsAbs(arg) {
			return fmt.Errorf("%s is outside the project directory; paths must be relative to the project root", arg)
		}
		if info, err := os.Stat(in); err == nil && info.IsDir() {
			return fmt.Errorf("%s is a directory; add individual files", in)
		} else if err != nil && !addNoEncrypt {
			return fmt.Errorf("%s not found (use --no-encrypt to add it to the config anyway)", in)
		}
		existing[in] = true
		entries = append(entries, config.FileMapping{Input: in, Output: addOutput, CopyTo: addCopyTo})
	}
	if len(entries) == 0 {
		return nil
	}

	// Validate the combined config (e.g. duplicate encrypted names) before changing anything
	combined := *cfg
	combined.Files = append(append([]config.FileMapping{}, cfg.Files...), entries...)
	if err := combined.Validate(); err != nil {
		return err
	}

	// Encrypt in memory first so a wrong password changes nothing
	var out []*sealed
	var m *manifest.Manifest
	if !addNoEncrypt {
		if m, err = manifest.Load(cfg.OutputDir); err != nil {
			return err
		}
		pwd, err := promptForPassword("🔐 Enter your SecureFlow password: ", m.PasswordHint)
		if err != nil {
			return err
		}

		hasEncrypted := false
		for _, fm := range cfg.Files {
			if utils.FileExists(encryptedPath(cfg, fm)) {
				hasEncrypted = true
				break
			}
		}
		if hasEncrypted {
			if err := verifyPassword(cfg, pwd, nil); err != nil {
				printHint(m.PasswordHint)
				return err
			}
		} else if passwordFromFlagOrEnv() == "" {
			// First encrypted file: confirm the password and offer a hint
			if err := confirmPassword(pwd); err != nil {
				return err
			}
			if m.PasswordHint, err = promptHint(m.PasswordHint); err != nil {
				return err
			}
		}

		for _, fm := range entries {
			s, err := sealInput(&combined, fm, pwd)
			if err != nil {
				return fmt.Errorf("failed to encrypt %s: %w", fm.Input, err)
			}
			out = append(out, s)
		}
	}

	// Now apply the changes
	if err := config.AppendFiles(cfgFile, entries); err != nil {
		return err
	}
	for _, e := range entries {
		fmt.Printf("%s ➕ Added %s -> %s%s\n", utils.ColorGreen, e.Input, e.EncryptedName(), utils.ColorReset)
	}

	if !addNoGitignore {
		var paths []string
		for _, e := range entries {
			paths = append(paths, e.Input)
			if e.CopyTo != "" {
				paths = append(paths, e.CopyTo)
			}
		}
		added, err := gitutil.AddToGitignore(".gitignore", paths)
		if err != nil {
			fmt.Printf("%s ⚠️  Warning: could not update .gitignore: %v%s\n", utils.ColorYellow, err, utils.ColorReset)
		} else if len(added) > 0 {
			fmt.Printf("%s 🙈 Added to .gitignore: %s%s\n", utils.ColorGreen, strings.Join(added, ", "), utils.ColorReset)
		}
	}

	if len(out) > 0 {
		if err := utils.EnsureDir(cfg.OutputDir); err != nil {
			return err
		}
		if err := writeSealed(out); err != nil {
			return err
		}
		for _, s := range out {
			m.Upsert(s.entry)
			fmt.Printf("%s 🔒 Encrypted %s -> %s%s\n", utils.ColorGreen, s.entry.Input, s.path, utils.ColorReset)
		}
		if err := m.Save(cfg.OutputDir); err != nil {
			return err
		}
	}

	// Plaintext that is already committed stays in git history even once ignored
	for _, e := range entries {
		if gitutil.IsTracked(e.Input) {
			fmt.Printf("\n%s ⚠️  %s is already tracked by git. Untrack it with:%s\n", utils.ColorYellow, e.Input, utils.ColorReset)
			fmt.Printf("     git rm --cached %q\n", e.Input)
			fmt.Println("   It also remains in your git history: rotate the secrets it contains.")
		}
	}

	if addNoEncrypt {
		fmt.Println("\nRun 'secureflow encrypt' when the files are ready.")
	}
	return nil
}
