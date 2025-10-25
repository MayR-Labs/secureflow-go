package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.OutputDir != "enc_keys" {
		t.Errorf("Expected output_dir 'enc_keys', got '%s'", cfg.OutputDir)
	}

	if cfg.TestOutputDir != "test_dec_keys" {
		t.Errorf("Expected test_output_dir 'test_dec_keys', got '%s'", cfg.TestOutputDir)
	}

	if len(cfg.Files) == 0 {
		t.Error("Expected default config to have file mappings")
	}

	// Verify first file mapping
	if len(cfg.Files) > 0 {
		firstFile := cfg.Files[0]
		if firstFile.Input == "" {
			t.Error("Expected first file mapping to have input")
		}
		if firstFile.Output == "" {
			t.Error("Expected first file mapping to have output")
		}
	}
}

func TestConfigSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "secureflow.yaml")

	// Create test config
	testConfig := &Config{
		OutputDir:     "test_enc",
		TestOutputDir: "test_dec",
		Files: []FileMapping{
			{Input: ".env.test", Output: ".env.test.encrypted"},
			{Input: "test.txt", Output: "test.txt.encrypted"},
		},
	}

	// Test save
	t.Run("Save", func(t *testing.T) {
		err := testConfig.Save(configPath)
		if err != nil {
			t.Fatalf("Save failed: %v", err)
		}

		// Verify file exists
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			t.Fatal("Config file was not created")
		}
	})

	// Test load
	t.Run("Load", func(t *testing.T) {
		loadedConfig, err := Load(configPath)
		if err != nil {
			t.Fatalf("Load failed: %v", err)
		}

		if loadedConfig.OutputDir != testConfig.OutputDir {
			t.Errorf("Expected OutputDir '%s', got '%s'", testConfig.OutputDir, loadedConfig.OutputDir)
		}

		if loadedConfig.TestOutputDir != testConfig.TestOutputDir {
			t.Errorf("Expected TestOutputDir '%s', got '%s'", testConfig.TestOutputDir, loadedConfig.TestOutputDir)
		}

		if len(loadedConfig.Files) != len(testConfig.Files) {
			t.Errorf("Expected %d files, got %d", len(testConfig.Files), len(loadedConfig.Files))
		}

		// Verify file mappings
		for i := range testConfig.Files {
			if i >= len(loadedConfig.Files) {
				break
			}

			if loadedConfig.Files[i].Input != testConfig.Files[i].Input {
				t.Errorf("File %d: Expected input '%s', got '%s'", i, testConfig.Files[i].Input, loadedConfig.Files[i].Input)
			}

			if loadedConfig.Files[i].Output != testConfig.Files[i].Output {
				t.Errorf("File %d: Expected output '%s', got '%s'", i, testConfig.Files[i].Output, loadedConfig.Files[i].Output)
			}
		}
	})
}

func TestLoadNonExistentFile(t *testing.T) {
	_, err := Load("/nonexistent/path/secureflow.yaml")
	if err == nil {
		t.Fatal("Expected error when loading non-existent file")
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	invalidPath := filepath.Join(tmpDir, "invalid.yaml")

	// Create invalid YAML
	invalidYAML := `
output_dir: enc_keys
files:
  - input: .env
    output: .env.encrypted
  - this is invalid yaml syntax
`
	if err := os.WriteFile(invalidPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("Failed to create invalid YAML file: %v", err)
	}

	_, err := Load(invalidPath)
	if err == nil {
		t.Fatal("Expected error when loading invalid YAML")
	}
}

func TestSaveToInvalidPath(t *testing.T) {
	cfg := DefaultConfig()
	err := cfg.Save("/nonexistent/directory/secureflow.yaml")
	if err == nil {
		t.Fatal("Expected error when saving to non-existent directory")
	}
}

func TestEmptyFilesList(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "empty.yaml")

	emptyConfig := &Config{
		OutputDir:     "enc",
		TestOutputDir: "dec",
		Files:         []FileMapping{},
	}

	// Save and load config with empty files list
	if err := emptyConfig.Save(configPath); err != nil {
		t.Fatalf("Failed to save empty config: %v", err)
	}

	loadedConfig, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load empty config: %v", err)
	}

	if len(loadedConfig.Files) != 0 {
		t.Errorf("Expected 0 files, got %d", len(loadedConfig.Files))
	}
}
