# 🔐 SecureFlow CLI

**SecureFlow** is a lightweight, cross-platform CLI tool for securely encrypting and decrypting sensitive files using AES-256 encryption. Built with Go, it provides a simple, reliable way to manage secrets in your projects and CI/CD pipelines.

**Perfect for**: Environment variables, keystores, SSL certificates, service credentials, API keys, database passwords, and any sensitive configuration files.

---

## ✨ Why SecureFlow?

- **🔒 Strong Security**: AES-256-CBC encryption with PBKDF2 key derivation (OpenSSL-compatible)
- **🚀 Fast & Lightweight**: Single binary, no dependencies, instant startup
- **🌍 Cross-Platform**: Works on Linux, macOS, and Windows
- **🤖 CI/CD Ready**: Non-interactive mode designed for automation
- **📝 Configuration-Based**: Simple YAML config for managing multiple files
- **🧪 Safe Testing**: Test decryption without overwriting existing files
- **📊 Detailed Reports**: Track what's encrypted, when, and metadata
- **💪 Production-Ready**: Battle-tested, reliable error handling

---

## 📋 Table of Contents

- [Quick Start](#-quick-start)
- [Installation](#-installation)
- [Usage](#-usage)
- [Configuration](#-configuration)
- [CI/CD Integration](#-cicd-integration)
- [Documentation](#-documentation)
- [Security](#-security)
- [Common Use Cases](#-common-use-cases)
- [Contributing](#-contributing)

---

## 🚀 Quick Start

```bash
# 1. Install SecureFlow (Linux/macOS)
curl -sSL https://raw.githubusercontent.com/MayR-Labs/secureflow-go/main/install.sh | bash

# 2. Initialize configuration
cd your-project
secureflow init

# 3. Edit secureflow.yaml to list your sensitive files
vim secureflow.yaml

# 4. Encrypt your files
secureflow encrypt

# 5. Commit encrypted files to git
git add enc_keys/ secureflow.yaml
git commit -m "Add encrypted secrets"

# 6. Add originals to .gitignore
echo ".env.prod" >> .gitignore
```

**That's it!** Your secrets are now encrypted and safe to commit. Your team can decrypt them with:

```bash
secureflow decrypt
```

---

## 🔧 Installation

### Quick Install (Linux/macOS)

The easiest way to install SecureFlow:

```bash
curl -sSL https://raw.githubusercontent.com/MayR-Labs/secureflow-go/main/install.sh | bash
```

Or with wget:

```bash
wget -qO- https://raw.githubusercontent.com/MayR-Labs/secureflow-go/main/install.sh | bash
```

### Manual Installation

#### Download Pre-built Binary

Download the appropriate binary for your OS from the [Releases](https://github.com/MayR-Labs/secureflow-go/releases) page.

**Linux (AMD64)**:
```bash
wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
chmod +x secureflow-linux-amd64
sudo mv secureflow-linux-amd64 /usr/local/bin/secureflow
```

**macOS (Intel)**:
```bash
wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-darwin-amd64
chmod +x secureflow-darwin-amd64
sudo mv secureflow-darwin-amd64 /usr/local/bin/secureflow
```

**macOS (Apple Silicon)**:
```bash
wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-darwin-arm64
chmod +x secureflow-darwin-arm64
sudo mv secureflow-darwin-arm64 /usr/local/bin/secureflow
```

**Windows**:
Download `secureflow-windows-amd64.exe` from the [Releases](https://github.com/MayR-Labs/secureflow-go/releases) page and add it to your PATH.

#### Build from Source

```bash
git clone https://github.com/MayR-Labs/secureflow-go.git
cd secureflow-go
go build -o secureflow
sudo mv secureflow /usr/local/bin/
```

### Verify Installation

```bash
secureflow --version
```

---

## 💻 Usage

### Initialize Configuration

Create a default `secureflow.yaml` configuration file:

```bash
secureflow init
```

Use a project template for quick setup:

```bash
# Interactive template selection
secureflow init

# Or specify a template directly
secureflow init --template flutter
secureflow init --template reactnative
secureflow init --template web
secureflow init --template docker
secureflow init --template k8s
secureflow init --template microservices
```

**Available templates:**
- `default` - React Native/Mobile app with Android/iOS files
- `reactnative` - React Native specific configuration
- `flutter` - Flutter mobile app configuration
- `web` - Web application with multiple environments
- `docker` - Docker deployment configuration
- `k8s` - Kubernetes secrets configuration
- `microservices` - Microservices architecture with multiple services

All templates include `copy_to: .env` for `.env.prod` files to automatically create `.env` after decryption.

This generates a configuration file like:

```yaml
# secureflow.yaml (Flutter template example)
output_dir: enc_keys
test_output_dir: test_dec_keys

files:
  - input: .env.prod
    output: .env.prod.encrypted
    copy_to: .env
  - input: android/app/keystore.jks
    output: keystore.jks.encrypted
```

**Edit this file** to match your project's sensitive files.

### Encrypt Files

**Interactive mode** (prompts for password and optional hint):

```bash
secureflow encrypt
```

**Non-interactive mode** (for scripts and CI/CD):

```bash
secureflow encrypt --password "your_password" --non-interactive
```

**Custom config file**:

```bash
secureflow encrypt --config ./path/to/custom-config.yaml
```

All encrypted files are saved to `enc_keys/` (or your configured `output_dir`).

### Decrypt Files

**Interactive mode**:

```bash
secureflow decrypt
```

**Non-interactive mode** (for CI/CD):

```bash
secureflow decrypt --password "$ENCRYPTION_PASSWORD" --non-interactive
```

**Custom config file**:

```bash
secureflow decrypt --config ./custom-config.yaml
```

### Test Decryption

Test decryption without overwriting existing files (decrypts to `test_dec_keys/`):

```bash
secureflow test
```

Or non-interactively:

```bash
secureflow test --password "your_password" --non-interactive
```

### View Help

```bash
secureflow --help
secureflow encrypt --help
secureflow decrypt --help
secureflow test --help
```

---

## 📝 Configuration

SecureFlow uses a YAML configuration file (`secureflow.yaml`) to define which files to encrypt/decrypt.

### Basic Configuration

```yaml
output_dir: enc_keys           # Where encrypted files are stored
test_output_dir: test_dec_keys # Where test decryption outputs go

files:
  - input: .env.production       # Source file (relative to project root)
    output: .env.production.encrypted  # Encrypted filename
    copy_to: .env                # (Optional) Copy decrypted file to this path

  - input: config/database.yml
    output: database.yml.encrypted
    
  - input: ssl/private.key
    output: ssl-private.key.encrypted
```

### Configuration Options

- **`output_dir`**: Directory for encrypted files (committed to git)
- **`test_output_dir`**: Directory for test decryption (added to .gitignore)
- **`files`**: Array of file entries
  - **`input`**: Path to source file (relative to project root)
  - **`output`**: Encrypted filename (just filename, not path)
  - **`copy_to`**: *(Optional)* Copy decrypted file to this path - useful when apps expect `.env` but you store `.env.prod`

### The `copy_to` Feature

When decrypting files, SecureFlow can automatically copy the decrypted file to another location. This is particularly useful for environment files where your application expects `.env` but you store `.env.prod` or `.env.production`:

```yaml
files:
  - input: .env.prod
    output: .env.prod.encrypted
    copy_to: .env  # Automatically creates .env from .env.prod after decryption
```

After decryption, both `.env.prod` and `.env` will exist with identical content.

### Example Configurations

See the [Configuration Guide](./docs/configuration.md) for detailed examples including:
- Mobile apps (Flutter, React Native)
- Web applications
- Microservices
- Docker deployments
- Kubernetes secrets
- Multi-environment setups

**Quick example configs**:
- [Basic](./docs/configuration.md#basic-configuration)
- [Mobile App](./docs/configuration.md#mobile-app-flutterreact-native)
- [Web Application](./docs/configuration.md#web-application)
- [Microservices](./docs/configuration.md#microservices-architecture)

---

## 🤖 CI/CD Integration

SecureFlow is designed to work seamlessly in CI/CD pipelines with its non-interactive mode.

### Quick CI/CD Example (GitHub Actions)

```yaml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v3
      
      # Install SecureFlow
      - name: Install SecureFlow
        run: |
          wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
          chmod +x secureflow-linux-amd64
          sudo mv secureflow-linux-amd64 /usr/local/bin/secureflow
      
      # Decrypt secrets
      - name: Decrypt secrets
        run: |
          secureflow decrypt --password "${{ secrets.SECUREFLOW_PASSWORD }}" --non-interactive
      
      # Build and deploy
      - name: Build & Deploy
        run: |
          npm run build
          ./deploy.sh
```

### Supported Platforms

SecureFlow works with all major CI/CD platforms:

- **GitHub Actions** ✅
- **GitLab CI** ✅
- **Bitbucket Pipelines** ✅
- **Jenkins** ✅
- **CircleCI** ✅
- **Azure Pipelines** ✅
- **Travis CI** ✅
- **Any CI/CD platform** that can run shell commands ✅

### Complete CI/CD Guide

See our comprehensive [CI/CD Usage Guide](./docs/cicd-usage.md) for:
- Platform-specific setup instructions
- Best practices for secrets management
- Multi-environment configurations
- Caching strategies
- Security considerations
- Troubleshooting CI/CD issues

---

## 📚 Documentation

Comprehensive documentation is available in the `docs/` directory:

### Core Documentation

- **[Configuration Guide](./docs/configuration.md)** - Detailed configuration options, examples, and best practices
- **[CI/CD Usage Guide](./docs/cicd-usage.md)** - Complete guide for integrating SecureFlow into your CI/CD pipelines
- **[Security Guide](./docs/security.md)** - Security best practices, encryption details, and compliance considerations
- **[Troubleshooting](./docs/troubleshooting.md)** - Common issues, solutions, and debugging tips

### Additional Resources

- **[Examples Directory](./examples/)** - Practical usage examples and sample configurations
  - [Usage Examples](./examples/README.md) - Real-world usage scenarios
  - [Sample Configs](./examples/configs.md) - Configuration templates for different project types

---

## 🔐 Security

### Encryption Details

- **Algorithm**: AES-256-CBC (Advanced Encryption Standard, 256-bit)
- **Key Derivation**: PBKDF2 with SHA-256
- **Format**: OpenSSL-compatible (Salted__ header + 8-byte salt + encrypted data)
- **Compatibility**: Files can be encrypted/decrypted with OpenSSL

### Security Model

- **Passwords are never stored or logged**
- **Strong encryption** (AES-256-CBC) with proper key derivation (PBKDF2)
- **OpenSSL compatible** - industry-standard format
- **Non-interactive mode** for secure CI/CD integration
- **Clean password handling** - passwords cleared from memory after use

### Best Practices

✅ **DO**:
- Use strong passwords (16+ characters, mixed case, numbers, symbols)
- Store passwords in password managers
- Rotate passwords periodically (every 3-6 months)
- Use different passwords for different environments
- Keep plaintext secrets in `.gitignore`
- Commit only encrypted files to version control

❌ **DON'T**:
- Commit plaintext secrets to git
- Share passwords via email or Slack
- Reuse passwords across projects
- Log or print passwords
- Store passwords in code or config files

For comprehensive security guidance, see our [Security Guide](./docs/security.md).

---

## 📊 Encryption Reports

Each encryption run generates a `report.txt` in your output directory with detailed metadata:

```
Encryption Report
=================
Note: Encrypted secrets for CI/CD
Password Hint: For the wise only
Created at: 2025-10-24
=================

File:           .env.prod
Encrypted As:   .env.prod.encrypted
Size (bytes):   348
Lines:          17
Last Modified:  2025-10-22 11:24:09
----------------------------------------
```

This helps you track:
- What files are encrypted
- When they were encrypted
- Password hints for team reference
- File sizes and metadata

---

## 🎯 Common Use Cases

### Local Development

Decrypt secrets when setting up a new project:

```bash
# Clone project
git clone https://github.com/yourorg/your-project.git
cd your-project

# Install SecureFlow
curl -sSL https://raw.githubusercontent.com/MayR-Labs/secureflow-go/main/install.sh | bash

# Decrypt secrets (get password from team)
secureflow decrypt
```

### CI/CD Deployment

Add to your CI/CD pipeline (example: GitHub Actions):

```yaml
- name: Decrypt secrets
  run: secureflow decrypt --password ${{ secrets.SECUREFLOW_PASSWORD }} --non-interactive
```

### Team Onboarding

New team member joining:

1. Install SecureFlow
2. Clone repository
3. Get decryption password from team lead (via secure channel)
4. Run `secureflow decrypt`
5. Start working!

### Environment Management

Manage secrets across multiple environments:

```bash
# Production
secureflow encrypt --config secureflow.prod.yaml --password "$PROD_PASS"

# Staging
secureflow encrypt --config secureflow.staging.yaml --password "$STAGING_PASS"

# Development
secureflow encrypt --config secureflow.dev.yaml --password "$DEV_PASS"
```

---

## 🧩 Project Structure

```
secureflow-go/
│
├── cmd/                    # CLI commands (Cobra)
│   ├── root.go            # Root command and global flags
│   ├── encrypt.go         # Encryption command
│   ├── decrypt.go         # Decryption command
│   ├── test.go            # Test decryption command
│   └── init.go            # Initialize config command
│
├── internal/              # Internal packages
│   ├── crypto/           # Encryption/decryption logic
│   ├── config/           # Configuration handling
│   └── utils/            # Utilities (file ops, logging)
│
├── docs/                 # Comprehensive documentation
│   ├── configuration.md  # Configuration guide
│   ├── cicd-usage.md    # CI/CD integration guide
│   ├── security.md      # Security best practices
│   └── troubleshooting.md # Troubleshooting guide
│
├── examples/             # Usage examples
│   ├── README.md        # Practical examples
│   └── configs.md       # Sample configurations
│
├── main.go              # Application entry point
├── install.sh           # Installation script
├── go.mod               # Go module definition
├── LICENSE              # MIT License
└── README.md            # This file
```

---

## 🧱 Error Handling

SecureFlow provides clear, actionable error messages:

- **Missing files** → Warns and skips, continues with other files
- **Wrong password** → Clear error message with exit code 1
- **Invalid YAML** → Shows line number and syntax error
- **Missing directories** → Automatically creates them
- **Interrupts (Ctrl+C)** → Graceful exit with cleanup notice

---

## 🧪 Development

### Run Locally

```bash
go run main.go encrypt
go run main.go decrypt
go run main.go test
go run main.go init
```

### Run Tests

```bash
go test ./...
```

With verbose output:

```bash
go test ./... -v
```

### Build Binary

```bash
go build -o secureflow
```

Cross-compile for different platforms:

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o secureflow-linux-amd64

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o secureflow-darwin-amd64

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o secureflow-darwin-arm64

# Windows
GOOS=windows GOARCH=amd64 go build -o secureflow-windows-amd64.exe
```

---

## 🧭 Roadmap

Future enhancements we're considering:

- [ ] Support for `.env` key filtering (only encrypt certain variables)
- [ ] Optional GPG-based encryption backend
- [ ] Progress bars for large files
- [ ] Batch re-encryption command
- [ ] Integration with Flutter build runners
- [ ] Support for multiple encryption backends
- [ ] Vault/secret manager integration
- [ ] Shell completion scripts

---

## 🤝 Contributing

Contributions are welcome! Here's how you can help:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/AmazingFeature`
3. **Make your changes**
4. **Run tests**: `go test ./...`
5. **Commit your changes**: `git commit -m 'Add some AmazingFeature'`
6. **Push to the branch**: `git push origin feature/AmazingFeature`
7. **Open a Pull Request**

### Development Guidelines

- Follow Go best practices and conventions
- Add tests for new features
- Update documentation as needed
- Keep commits focused and descriptive
- Ensure all tests pass before submitting PR

---

## 📄 License

SecureFlow is licensed under the [MIT License](LICENSE).

Copyright (c) 2025 [MayR Labs](https://mayrlabs.com)

---

## 🔗 Links

- **[GitHub Repository](https://github.com/MayR-Labs/secureflow-go)** - Source code and issues
- **[Latest Releases](https://github.com/MayR-Labs/secureflow-go/releases)** - Download binaries
- **[Documentation](./docs/)** - Comprehensive guides
- **[Examples](./examples/)** - Practical usage examples
- **[MayR Labs](https://mayrlabs.com)** - Our website
- **[Report Issues](https://github.com/MayR-Labs/secureflow-go/issues)** - Bug reports and feature requests

---

## 🙏 Acknowledgments

SecureFlow is built with these excellent libraries:

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [yaml.v3](https://gopkg.in/yaml.v3) - YAML parsing
- Go standard library - Crypto, file I/O, and more

Special thanks to all contributors and users who help make SecureFlow better!

---

## ⭐ Star Us!

If you find SecureFlow useful, please consider giving us a star on GitHub! It helps others discover the project.

[![GitHub stars](https://img.shields.io/github/stars/MayR-Labs/secureflow-go?style=social)](https://github.com/MayR-Labs/secureflow-go/stargazers)

---

**Made with ❤️ by [MayR Labs](https://mayrlabs.com)**
