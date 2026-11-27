package cmd

import (
	"bufio"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

//go:embed secureflow.sh
var secureflowShellScript string

var installLocalCmd = &cobra.Command{
	Use:   "install-local",
	Short: "Install secureflow locally for CI/CD use",
	Long: `Downloads platform-specific executables and creates a launcher script.
This allows running secureflow in CI/CD without needing to install it globally.

The command will:
- Download executables for Linux, macOS, and Windows
- Save them to .secureflow/ directory
- Create a secureflow.sh launcher script in the current directory

Usage in CI/CD:
  ./secureflow.sh decrypt --password "$PASSWORD" --non-interactive`,
	RunE: runInstallLocal,
}

func init() {
	rootCmd.AddCommand(installLocalCmd)
}

func runInstallLocal(cmd *cobra.Command, args []string) error {
	fmt.Println("🔧 Setting up local secureflow installation...")
	fmt.Println()

	// Create .secureflow directory
	secureflowDir := ".secureflow"
	if err := os.MkdirAll(secureflowDir, 0755); err != nil {
		return fmt.Errorf("failed to create %s directory: %w", secureflowDir, err)
	}
	fmt.Printf("✅ Created %s directory\n", secureflowDir)

	// Get the latest release version
	version := Version
	if version == "" {
		version = "latest"
	}

	// Define platforms to download
	platforms := []struct {
		os   string
		arch string
	}{
		{"linux", "amd64"},
		{"linux", "arm64"},
		{"darwin", "amd64"},
		{"darwin", "arm64"},
		{"windows", "amd64"},
	}

	// Download executables for each platform
	fmt.Println()
	fmt.Println("📥 Downloading platform-specific executables...")
	fmt.Println()

	baseURL := "https://github.com/MayR-Labs/secureflow-go/releases/latest/download"
	if version != "latest" && version != "" {
		// Add 'v' prefix for tag-based releases
		if version[0] != 'v' {
			version = "v" + version
		}
		baseURL = fmt.Sprintf("https://github.com/MayR-Labs/secureflow-go/releases/download/%s", version)
	}

	// Fetch the release checksums; every binary is verified against them
	checksums, err := fetchChecksums(baseURL + "/checksums.txt")
	if err != nil {
		return fmt.Errorf("failed to download checksums.txt (refusing to install unverified binaries): %w", err)
	}

	successCount := 0
	for _, platform := range platforms {
		binaryName := fmt.Sprintf("secureflow-%s-%s", platform.os, platform.arch)
		if platform.os == "windows" {
			binaryName += ".exe"
		}

		url := fmt.Sprintf("%s/%s", baseURL, binaryName)
		outputPath := filepath.Join(secureflowDir, binaryName)

		fmt.Printf("  Downloading %s...\n", binaryName)
		if err := downloadFile(url, outputPath); err != nil {
			fmt.Printf("  ⚠️  Warning: Failed to download %s: %v\n", binaryName, err)
			continue
		}
		if err := verifyChecksum(outputPath, checksums[binaryName]); err != nil {
			os.Remove(outputPath)
			fmt.Printf("  ❌ %s: %v\n", binaryName, err)
			continue
		}

		// Make executable (Unix systems)
		if platform.os != "windows" {
			if err := os.Chmod(outputPath, 0755); err != nil {
				fmt.Printf("  ⚠️  Warning: Failed to set permissions on %s: %v\n", binaryName, err)
			}
		}

		fmt.Printf("  ✅ Downloaded %s\n", binaryName)
		successCount++
	}

	// Check if at least one binary was downloaded
	if successCount == 0 {
		return fmt.Errorf("failed to download any executables")
	}

	fmt.Println()
	fmt.Printf("✅ Successfully downloaded %d platform executable(s)\n", successCount)

	// Create secureflow.sh launcher script
	fmt.Println()
	fmt.Println("📝 Creating launcher script...")

	launcherPath := "secureflow.sh"
	if err := os.WriteFile(launcherPath, []byte(secureflowShellScript), 0755); err != nil {
		return fmt.Errorf("failed to create launcher script: %w", err)
	}
	fmt.Printf("✅ Created %s\n", launcherPath)

	fmt.Println()
	fmt.Println("🎉 Local installation complete!")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  ./secureflow.sh encrypt")
	fmt.Println("  ./secureflow.sh decrypt --password \"$PASSWORD\" --non-interactive")
	fmt.Println("  ./secureflow.sh --help")
	fmt.Println()
	fmt.Println("The launcher script will automatically select the correct executable for your platform.")

	return nil
}

// fetchChecksums downloads a sha256sum-style checksums.txt into name -> hex digest
func fetchChecksums(url string) (map[string]string, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}
	sums := map[string]string{}
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 2 {
			sums[strings.TrimPrefix(fields[1], "*")] = strings.ToLower(fields[0])
		}
	}
	return sums, scanner.Err()
}

// verifyChecksum checks path's SHA-256 against expected
func verifyChecksum(path, expected string) error {
	if expected == "" {
		return fmt.Errorf("no checksum listed in checksums.txt")
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	if actual := hex.EncodeToString(h.Sum(nil)); actual != expected {
		return fmt.Errorf("checksum mismatch (expected %s, got %s)", expected, actual)
	}
	return nil
}

// downloadFile downloads a file from a URL and saves it to the specified path
func downloadFile(url, outputPath string) error {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create the output file
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Download the file
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the response body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
