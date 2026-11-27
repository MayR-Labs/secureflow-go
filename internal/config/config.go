package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// DefaultTestOutputDir is used when test_output_dir is not set
const DefaultTestOutputDir = "test_dec_keys"

// EncryptedSuffix is appended to inferred encrypted filenames
const EncryptedSuffix = ".encrypted"

// FileMapping represents a file to be encrypted or decrypted
type FileMapping struct {
	Input  string `yaml:"input"`
	Output string `yaml:"output,omitempty"`  // Optional: inferred from Input when empty
	CopyTo string `yaml:"copy_to,omitempty"` // Optional: copy decrypted file to this path
}

// EncryptedName returns the filename of the encrypted file inside output_dir.
// If Output is set it is returned as-is; otherwise it is inferred from Input.
func (f FileMapping) EncryptedName() string {
	if strings.TrimSpace(f.Output) != "" {
		return f.Output
	}
	return InferOutputName(f.Input)
}

// InferOutputName derives an encrypted filename from an input path by
// flattening its directories with "-" and appending ".encrypted".
//
//	android/app/keystore.jks -> android-app-keystore.jks.encrypted
//	.env.prod                -> .env.prod.encrypted
//	./config/db.yml          -> config-db.yml.encrypted
func InferOutputName(input string) string {
	// Normalise Windows separators so configs behave the same on every OS
	p := strings.ReplaceAll(input, "\\", "/")
	p = filepath.ToSlash(filepath.Clean(p))

	var parts []string
	for _, part := range strings.Split(p, "/") {
		switch part {
		case "", ".", "..":
			continue
		}
		// Drop a Windows drive letter such as "C:"
		if len(part) == 2 && part[1] == ':' {
			continue
		}
		parts = append(parts, part)
	}

	return strings.Join(parts, "-") + EncryptedSuffix
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

	if strings.TrimSpace(cfg.TestOutputDir) == "" {
		cfg.TestOutputDir = DefaultTestOutputDir
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Inputs returns the input path of every file entry
func (c *Config) Inputs() []string {
	inputs := make([]string, len(c.Files))
	for i, f := range c.Files {
		inputs[i] = f.Input
	}
	return inputs
}

// Validate checks the configuration for common mistakes, such as missing
// inputs or two files that would be encrypted to the same output name.
func (c *Config) Validate() error {
	if strings.TrimSpace(c.OutputDir) == "" {
		return fmt.Errorf("invalid config: output_dir is required")
	}
	seen := make(map[string]string, len(c.Files))
	for i, f := range c.Files {
		if strings.TrimSpace(f.Input) == "" {
			return fmt.Errorf("invalid config: files[%d] is missing 'input'", i)
		}
		name := f.EncryptedName()
		if prev, ok := seen[name]; ok {
			return fmt.Errorf("invalid config: %q and %q both encrypt to %q; set a unique 'output' for one of them", prev, f.Input, name)
		}
		seen[name] = f.Input
	}

	return nil
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

// AppendFiles adds entries to the `files` list of the config at path. It edits
// the YAML node tree rather than re-marshalling the struct, so comments and
// the order of existing keys are preserved.
func AppendFiles(path string, entries []FileMapping) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}
	if doc.Kind != yaml.DocumentNode || len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		return fmt.Errorf("config file %s is not a YAML mapping", path)
	}
	root := doc.Content[0]

	var files *yaml.Node
	for i := 0; i+1 < len(root.Content); i += 2 {
		if root.Content[i].Value == "files" {
			files = root.Content[i+1]
			break
		}
	}
	if files == nil {
		files = &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		root.Content = append(root.Content, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "files"}, files)
	}
	if files.Kind == yaml.ScalarNode && (files.Tag == "!!null" || files.Value == "") {
		*files = yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	}
	if files.Kind != yaml.SequenceNode {
		return fmt.Errorf("'files' in %s must be a list", path)
	}
	files.Style = 0 // force block style, e.g. if it was written as `files: []`

	str := func(v string) *yaml.Node { return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: v} }
	for _, e := range entries {
		item := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
		item.Content = append(item.Content, str("input"), str(e.Input))
		if e.Output != "" {
			item.Content = append(item.Content, str("output"), str(e.Output))
		}
		if e.CopyTo != "" {
			item.Content = append(item.Content, str("copy_to"), str(e.CopyTo))
		}
		files.Content = append(files.Content, item)
	}

	var buf strings.Builder
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(&doc); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}
	if err := enc.Close(); err != nil {
		return err
	}

	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte(buf.String()), info.Mode().Perm())
}
