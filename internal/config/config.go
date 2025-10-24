package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// FileMapping represents a file to be encrypted or decrypted
type FileMapping struct {
	Input  string `yaml:"input"`
	Output string `yaml:"output"`
}

// Config represents the secureflow.yaml configuration
type Config struct {
	OutputDir     string        `yaml:"output_dir"`
	TestOutputDir string        `yaml:"test_output_dir"`
	Files         []FileMapping `yaml:"files"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		OutputDir:     "enc_keys",
		TestOutputDir: "test_dec_keys",
		Files: []FileMapping{
			{Input: ".env.prod", Output: ".env.prod.encrypted"},
			{Input: "android/app/keystore.jks", Output: "keystore.jks.encrypted"},
			{Input: "android/key.properties", Output: "key.properties.encrypted"},
			{Input: "android/service-key.json", Output: "service-key.json.encrypted"},
		},
	}
}

// Load reads and parses a secureflow.yaml file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

// Save writes the configuration to a YAML file
func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
