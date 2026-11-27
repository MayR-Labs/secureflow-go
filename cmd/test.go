package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/MayR-Labs/secureflow-go/internal/config"
	"github.com/MayR-Labs/secureflow-go/internal/manifest"
	"github.com/MayR-Labs/secureflow-go/internal/utils"
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
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	hint := manifest.ReadHint(cfg.OutputDir)

	pwd, err := promptForPassword("🔐 [TEST] Enter password to test decrypt your secrets: ", hint)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf("%s 🔐 [TEST] Starting decryption process...%s\n\n", utils.ColorYellow, utils.ColorReset)

	files, err := openAll(cfg, pwd, hint)
	if err != nil {
		return fmt.Errorf("test %w", err)
	}

	fmt.Println()
	for _, f := range files {
		// Mirror the input path under the test directory so files with the same
		// name in different folders (e.g. services/*/.env.prod) don't collide
		out := filepath.Join(cfg.TestOutputDir, filepath.Clean("/"+filepath.ToSlash(f.fm.Input)))
		if err := writePlaintext(out, f.plaintext); err != nil {
			return fmt.Errorf("failed to write %s: %w", out, err)
		}
		fmt.Printf("%s ✅ %s -> %s%s\n", utils.ColorGreen, f.fm.Input, out, utils.ColorReset)
	}

	warnLegacy(files)

	fmt.Println()
	fmt.Printf("%s 🎉 Test decryption successful! (%d file(s))%s\n", utils.ColorGreen, len(files), utils.ColorReset)
	fmt.Printf("📁 Test files saved to: %s\n", cfg.TestOutputDir)
	return nil
}
