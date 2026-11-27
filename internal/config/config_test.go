package config

import (
	"os"
	"path/filepath"
	"strings"
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

func TestCopyToField(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "copy_to.yaml")

	// Create config with copy_to field
	testConfig := &Config{
		OutputDir:     "enc",
		TestOutputDir: "dec",
		Files: []FileMapping{
			{Input: ".env.prod", Output: ".env.prod.encrypted", CopyTo: ".env"},
			{Input: "config.json", Output: "config.json.encrypted"},
		},
	}

	// Save and load
	if err := testConfig.Save(configPath); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	loadedConfig, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify first file has copy_to
	if loadedConfig.Files[0].CopyTo != ".env" {
		t.Errorf("Expected CopyTo '.env', got '%s'", loadedConfig.Files[0].CopyTo)
	}

	// Verify second file has empty copy_to
	if loadedConfig.Files[1].CopyTo != "" {
		t.Errorf("Expected empty CopyTo, got '%s'", loadedConfig.Files[1].CopyTo)
	}
}

func TestTemplateConfig(t *testing.T) {
	tests := []struct {
		name         string
		template     string
		expectFiles  int
		expectCopyTo bool
	}{
		{"Default", "default", 4, true},
		{"ReactNative", "reactnative", 6, true},
		{"React-Native", "react-native", 6, true},
		{"Flutter", "flutter", 6, true},
		{"Web", "web", 5, true},
		{"Docker", "docker", 4, true},
		{"K8s", "k8s", 4, true},
		{"Kubernetes", "kubernetes", 4, true},
		{"Microservices", "microservices", 5, true},
		{"Invalid", "invalid-template", 4, true}, // Should default
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := TemplateConfig(tt.template)

			if cfg == nil {
				t.Fatal("TemplateConfig returned nil")
			}

			if len(cfg.Files) != tt.expectFiles {
				t.Errorf("Expected %d files, got %d", tt.expectFiles, len(cfg.Files))
			}

			if tt.expectCopyTo {
				// Check if at least one file has copy_to set
				hasCopyTo := false
				for _, f := range cfg.Files {
					if f.CopyTo != "" {
						hasCopyTo = true
						break
					}
				}
				if !hasCopyTo {
					t.Error("Expected at least one file to have copy_to set")
				}
			}
		})
	}
}

func TestReactNativeConfig(t *testing.T) {
	cfg := ReactNativeConfig()
	if cfg == nil {
		t.Fatal("ReactNativeConfig returned nil")
	}
	if len(cfg.Files) == 0 {
		t.Error("Expected non-empty files list")
	}
	// First file should have copy_to for .env
	if cfg.Files[0].CopyTo != ".env" {
		t.Errorf("Expected first file to have copy_to '.env', got '%s'", cfg.Files[0].CopyTo)
	}
}

func TestFlutterConfig(t *testing.T) {
	cfg := FlutterConfig()
	if cfg == nil {
		t.Fatal("FlutterConfig returned nil")
	}
	if len(cfg.Files) == 0 {
		t.Error("Expected non-empty files list")
	}
}

func TestWebConfig(t *testing.T) {
	cfg := WebConfig()
	if cfg == nil {
		t.Fatal("WebConfig returned nil")
	}
	if len(cfg.Files) == 0 {
		t.Error("Expected non-empty files list")
	}
}

func TestDockerConfig(t *testing.T) {
	cfg := DockerConfig()
	if cfg == nil {
		t.Fatal("DockerConfig returned nil")
	}
	if len(cfg.Files) == 0 {
		t.Error("Expected non-empty files list")
	}
}

func TestK8sConfig(t *testing.T) {
	cfg := K8sConfig()
	if cfg == nil {
		t.Fatal("K8sConfig returned nil")
	}
	if len(cfg.Files) == 0 {
		t.Error("Expected non-empty files list")
	}
}

func TestMicroservicesConfig(t *testing.T) {
	cfg := MicroservicesConfig()
	if cfg == nil {
		t.Fatal("MicroservicesConfig returned nil")
	}
	if len(cfg.Files) == 0 {
		t.Error("Expected non-empty files list")
	}
}

func TestInferOutputName(t *testing.T) {
	tests := map[string]string{
		"android/app/keystore.jks": "android-app-keystore.jks.encrypted",
		".env.prod":                ".env.prod.encrypted",
		"./config/database.yml":    "config-database.yml.encrypted",
		"services//auth/.env.prod": "services-auth-.env.prod.encrypted",
		"../shared/secret.key":     "shared-secret.key.encrypted",
		"/etc/app/secret.pem":      "etc-app-secret.pem.encrypted",
		`android\app\keystore.jks`: "android-app-keystore.jks.encrypted",
	}
	for input, want := range tests {
		if got := InferOutputName(input); got != want {
			t.Errorf("InferOutputName(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestEncryptedNamePrefersExplicitOutput(t *testing.T) {
	f := FileMapping{Input: "android/app/keystore.jks", Output: "keystore.jks.encrypted"}
	if got := f.EncryptedName(); got != "keystore.jks.encrypted" {
		t.Errorf("EncryptedName() = %q, want explicit output", got)
	}
	f.Output = ""
	if got := f.EncryptedName(); got != "android-app-keystore.jks.encrypted" {
		t.Errorf("EncryptedName() = %q, want inferred output", got)
	}
}

func TestValidate(t *testing.T) {
	valid := &Config{OutputDir: "enc", Files: []FileMapping{{Input: "a/.env"}, {Input: "b/.env"}}}
	if err := valid.Validate(); err != nil {
		t.Errorf("expected valid config, got %v", err)
	}

	cases := map[string]*Config{
		"missing output_dir": {Files: []FileMapping{{Input: ".env"}}},
		"missing input":      {OutputDir: "enc", Files: []FileMapping{{Output: "x.encrypted"}}},
		"duplicate output": {OutputDir: "enc", Files: []FileMapping{
			{Input: "android/app/keystore.jks"},
			{Input: "other.jks", Output: "android-app-keystore.jks.encrypted"},
		}},
	}
	for name, cfg := range cases {
		if err := cfg.Validate(); err == nil {
			t.Errorf("%s: expected validation error", name)
		}
	}
}

func TestLoadWithoutOutput(t *testing.T) {
	path := filepath.Join(t.TempDir(), "secureflow.yaml")
	yaml := "output_dir: enc_keys\nfiles:\n  - input: android/app/keystore.jks\n"
	if err := os.WriteFile(path, []byte(yaml), 0644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if got := cfg.Files[0].EncryptedName(); got != "android-app-keystore.jks.encrypted" {
		t.Errorf("EncryptedName() = %q", got)
	}
}

func TestAppendFilesPreservesComments(t *testing.T) {
	path := filepath.Join(t.TempDir(), "secureflow.yaml")
	original := `# Project secrets
output_dir: enc_keys # committed
test_output_dir: test_dec_keys

files:
  # production env
  - input: .env.prod
    copy_to: .env
`
	if err := os.WriteFile(path, []byte(original), 0644); err != nil {
		t.Fatal(err)
	}
	if err := AppendFiles(path, []FileMapping{{Input: "android/app/keystore.jks"}, {Input: "a.json", Output: "x.encrypted", CopyTo: "b.json"}}); err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(path)
	got := string(data)
	for _, want := range []string{"# Project secrets", "# committed", "# production env", "input: android/app/keystore.jks", "output: x.encrypted"} {
		if !strings.Contains(got, want) {
			t.Errorf("result missing %q:\n%s", want, got)
		}
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Files) != 3 || cfg.Files[1].EncryptedName() != "android-app-keystore.jks.encrypted" || cfg.Files[2].CopyTo != "b.json" {
		t.Errorf("unexpected files: %+v", cfg.Files)
	}
}

func TestAppendFilesEmptyList(t *testing.T) {
	path := filepath.Join(t.TempDir(), "secureflow.yaml")
	if err := os.WriteFile(path, []byte("output_dir: enc\nfiles: []\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := AppendFiles(path, []FileMapping{{Input: ".env"}}); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(path)
	if err != nil || len(cfg.Files) != 1 {
		t.Fatalf("got %+v, %v", cfg, err)
	}
}
