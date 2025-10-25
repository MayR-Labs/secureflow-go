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
	CopyTo string `yaml:"copy_to,omitempty"` // Optional: copy decrypted file to this path
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
			{Input: ".env.prod", Output: ".env.prod.encrypted", CopyTo: ".env"},
			{Input: "android/app/keystore.jks", Output: "keystore.jks.encrypted"},
			{Input: "android/key.properties", Output: "key.properties.encrypted"},
			{Input: "android/service-key.json", Output: "service-key.json.encrypted"},
		},
	}
}

// TemplateConfig returns a configuration based on the template name
func TemplateConfig(templateName string) *Config {
	switch templateName {
	case "reactnative", "react-native":
		return ReactNativeConfig()
	case "flutter":
		return FlutterConfig()
	case "web":
		return WebConfig()
	case "docker":
		return DockerConfig()
	case "k8s", "kubernetes":
		return K8sConfig()
	case "microservices":
		return MicroservicesConfig()
	default:
		return DefaultConfig()
	}
}

// ReactNativeConfig returns a React Native project configuration
func ReactNativeConfig() *Config {
	return &Config{
		OutputDir:     "enc_keys",
		TestOutputDir: "test_dec_keys",
		Files: []FileMapping{
			{Input: ".env.prod", Output: ".env.prod.encrypted", CopyTo: ".env"},
			{Input: ".env.staging", Output: ".env.staging.encrypted"},
			{Input: "android/app/keystore.jks", Output: "keystore.jks.encrypted"},
			{Input: "android/key.properties", Output: "key.properties.encrypted"},
			{Input: "android/service-key.json", Output: "service-key.json.encrypted"},
			{Input: "ios/GoogleService-Info.plist", Output: "GoogleService-Info.plist.encrypted"},
		},
	}
}

// FlutterConfig returns a Flutter project configuration
func FlutterConfig() *Config {
	return &Config{
		OutputDir:     "enc_keys",
		TestOutputDir: "test_dec_keys",
		Files: []FileMapping{
			{Input: ".env.prod", Output: ".env.prod.encrypted", CopyTo: ".env"},
			{Input: ".env.staging", Output: ".env.staging.encrypted"},
			{Input: "android/app/keystore.jks", Output: "keystore.jks.encrypted"},
			{Input: "android/key.properties", Output: "key.properties.encrypted"},
			{Input: "android/app/google-services.json", Output: "google-services.json.encrypted"},
			{Input: "ios/Runner/GoogleService-Info.plist", Output: "GoogleService-Info.plist.encrypted"},
		},
	}
}

// WebConfig returns a web application configuration
func WebConfig() *Config {
	return &Config{
		OutputDir:     "enc_keys",
		TestOutputDir: "test_dec_keys",
		Files: []FileMapping{
			{Input: ".env.prod", Output: ".env.prod.encrypted", CopyTo: ".env"},
			{Input: ".env.production", Output: ".env.production.encrypted"},
			{Input: ".env.staging", Output: ".env.staging.encrypted"},
			{Input: "config/database.yml", Output: "database.yml.encrypted"},
			{Input: "config/secrets.yml", Output: "secrets.yml.encrypted"},
		},
	}
}

// DockerConfig returns a Docker deployment configuration
func DockerConfig() *Config {
	return &Config{
		OutputDir:     "docker/secrets/encrypted",
		TestOutputDir: "docker/secrets/test",
		Files: []FileMapping{
			{Input: ".env.prod", Output: ".env.prod.encrypted", CopyTo: ".env"},
			{Input: "docker/.env.production", Output: "docker-env.production.encrypted"},
			{Input: "docker/compose/.env.db", Output: "docker-env.db.encrypted"},
			{Input: "docker/nginx/ssl/private.key", Output: "nginx-ssl-private.key.encrypted"},
		},
	}
}

// K8sConfig returns a Kubernetes configuration
func K8sConfig() *Config {
	return &Config{
		OutputDir:     "k8s/encrypted-secrets",
		TestOutputDir: "k8s/test-secrets",
		Files: []FileMapping{
			{Input: ".env.prod", Output: ".env.prod.encrypted", CopyTo: ".env"},
			{Input: "k8s/secrets/database-credentials.yaml", Output: "database-credentials.yaml.encrypted"},
			{Input: "k8s/secrets/api-keys.yaml", Output: "api-keys.yaml.encrypted"},
			{Input: "k8s/secrets/tls-cert.yaml", Output: "tls-cert.yaml.encrypted"},
		},
	}
}

// MicroservicesConfig returns a microservices architecture configuration
func MicroservicesConfig() *Config {
	return &Config{
		OutputDir:     "encrypted",
		TestOutputDir: "decrypted_test",
		Files: []FileMapping{
			{Input: ".env.prod", Output: ".env.prod.encrypted", CopyTo: ".env"},
			{Input: "services/auth/.env.prod", Output: "auth-env.prod.encrypted", CopyTo: "services/auth/.env"},
			{Input: "services/api/.env.prod", Output: "api-env.prod.encrypted", CopyTo: "services/api/.env"},
			{Input: "services/worker/.env.prod", Output: "worker-env.prod.encrypted", CopyTo: "services/worker/.env"},
			{Input: "shared/redis.conf", Output: "shared-redis.conf.encrypted"},
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
