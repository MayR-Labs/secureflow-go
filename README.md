# 🔐 SecureFlow CLI

**SecureFlow** is a lightweight, cross-platform CLI tool for securely encrypting and decrypting sensitive files using AES-256 encryption. Built with Go, it provides a simple, reliable way to manage secrets in your projects and CI/CD pipelines.

**Perfect for**: Environment variables, keystores, SSL certificates, service credentials, API keys, database passwords, and any sensitive configuration files.

---

## ✨ Why SecureFlow?

- **🔒 Strong Security**: AES-256-GCM authenticated encryption, PBKDF2-SHA256 (600,000 iterations)
- **🚀 Fast & Lightweight**: Single binary, no dependencies, instant startup
- **🌍 Cross-Platform**: Works on Linux, macOS, and Windows
- **🤖 CI/CD Ready**: Non-interactive mode designed for automation
- **📝 Configuration-Based**: Simple YAML config for managing multiple files
- **🧪 Safe Testing**: Test decryption without overwriting existing files
- **📊 JSON Manifest**: Track what's encrypted, when, and the password hint
- **🔄 Easy Rotation**: `secureflow refresh` re-encrypts everything with a new password
- **➕ One-step adds**: `secureflow add <file>` updates the config, `.gitignore`, and encrypts
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

### Local Installation for CI/CD

SecureFlow provides a convenient `install-local` command that sets up a portable installation without requiring system-wide installation. This is perfect for CI/CD pipelines where you can't or don't want to install binaries globally.

```bash
# Run this once to set up local installation
secureflow install-local
```

This command will:
- Download platform-specific executables (Linux, macOS, Windows) to `.secureflow/` directory
- Create a `secureflow.sh` launcher script in the current directory
- The launcher automatically detects your platform and runs the correct executable

**Usage in CI/CD:**

```bash
# In your CI/CD pipeline
./secureflow.sh decrypt --password "$SECUREFLOW_PASSWORD" --non-interactive
./secureflow.sh encrypt --password "$SECUREFLOW_PASSWORD" --non-interactive
```

**Benefits:**
- ✅ No system installation required
- ✅ Works across different platforms automatically
- ✅ Portable - commit `.secureflow/` and `secureflow.sh` to your repo
- ✅ Perfect for CI/CD environments with restricted permissions

**Example: Using in GitHub Actions**

**Option A: Install during CI/CD run**
```yaml
steps:
  - uses: actions/checkout@v4
  
  - name: Setup SecureFlow locally
    run: |
      wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
      chmod +x secureflow-linux-amd64
      ./secureflow-linux-amd64 install-local
  
  - name: Decrypt secrets
    run: ./secureflow.sh decrypt --password "${{ secrets.SECUREFLOW_PASSWORD }}" --non-interactive
```

**Option B: Commit installation to repo (no CI/CD setup needed)**
```bash
# Run locally once
secureflow install-local
git add .secureflow/ secureflow.sh
git commit -m "Add portable SecureFlow"
```

Then your CI/CD workflow is simply:
```yaml
steps:
  - uses: actions/checkout@v4
  
  - name: Decrypt secrets
    run: ./secureflow.sh decrypt --password "${{ secrets.SECUREFLOW_PASSWORD }}" --non-interactive
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

Interactively you're asked to confirm the password, and the previous hint is offered as the default.

**Non-interactive mode** (for scripts and CI/CD):

```bash
secureflow encrypt --password "your_password" --non-interactive --hint "the usual" --note "Q4 rotation"
```

**Custom config file**:

```bash
secureflow encrypt --config ./path/to/custom-config.yaml
```

All encrypted files are saved to `enc_keys/` (or your configured `output_dir`).

### Decrypt Files

**Interactive mode** (shows the password hint saved at encryption time, if any):

```bash
secureflow decrypt
# 💡 Password hint: For the wise only
# 🔐 Enter password to decrypt your secrets:
```

The hint is shown again if decryption fails.

**Non-interactive mode** (for CI/CD):

```bash
secureflow decrypt --password "$ENCRYPTION_PASSWORD" --non-interactive

# or keep the password off the command line
SECUREFLOW_PASSWORD="$ENCRYPTION_PASSWORD" secureflow decrypt --non-interactive
```

**Custom config file**:

```bash
secureflow decrypt --config ./custom-config.yaml
```

### Add a File

Add a file to `secureflow.yaml`, ignore it in `.gitignore`, and encrypt it with your existing password, all in one step:

```bash
secureflow add .env.prod --copy-to .env
secureflow add android/app/keystore.jks android/key.properties
secureflow add ios/GoogleService-Info.plist --no-encrypt   # config + .gitignore only
```

The password is checked against your existing encrypted files, so everything stays under one password. Comments in `secureflow.yaml` are preserved. If the file is already tracked by git, SecureFlow tells you how to untrack it.

### Refresh / Rotate the Password

Re-encrypt every file, either with a new password or with the same one:

```bash
secureflow refresh                       # interactive: old password, then optional new one
secureflow refresh --non-interactive \
  --old-password "$OLD" --new-password "$NEW" --hint "Vault: Q1 2027"
```

- Decrypted content stays in memory only; nothing is written to disk in plaintext.
- All-or-nothing: if any file fails to decrypt with the old password, no file is changed.
- Leaving out the new password re-encrypts with the current one. This is how you upgrade legacy AES-CBC files to AES-GCM.
- `secureflow rotate` is an alias.

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
secureflow init --help
secureflow encrypt --help
secureflow decrypt --help
secureflow test --help
secureflow add --help
secureflow refresh --help
secureflow install-local --help
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
    output: .env.production.encrypted  # (Optional) Encrypted filename
    copy_to: .env                # (Optional) Copy decrypted file to this path

  - input: config/database.yml   # output inferred: config-database.yml.encrypted

  - input: android/app/keystore.jks  # output inferred: android-app-keystore.jks.encrypted
```

### Configuration Options

- **`output_dir`**: Directory for encrypted files (committed to git)
- **`test_output_dir`**: Directory for test decryption (added to .gitignore)
- **`files`**: Array of file entries
  - **`input`**: Path to source file (relative to project root)
  - **`output`**: *(Optional)* Encrypted filename (just filename, not path). When omitted it's inferred from `input` by replacing `/` with `-` and appending `.encrypted`, e.g. `android/app/keystore.jks` → `android-app-keystore.jks.encrypted`. SecureFlow refuses configs where two files would end up with the same encrypted name.
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

### GitHub Action (Recommended for GitHub)

Use the SecureFlow action as a step in the same job that builds your app. It downloads the right binary for the runner (Linux/macOS/Windows, x64/ARM64), verifies its SHA-256 checksum, adds `secureflow` to `PATH`, and decrypts your secrets:

```yaml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Decrypt secrets
        uses: MayR-Labs/secureflow-go@v1.3.0
        with:
          password: ${{ secrets.SECUREFLOW_PASSWORD }}

      - name: Build & Deploy
        run: |
          npm run build
          ./deploy.sh
```

| Input | Default | Description |
|---|---|---|
| `password` | — | Encryption password (from a secret). Required unless `command: install`. |
| `command` | `decrypt` | `decrypt`, `test`, `encrypt`, or `install` (install only). |
| `config` | `secureflow.yaml` | Config file path, relative to `working-directory`. |
| `working-directory` | `.` | Directory to run in (useful for monorepos). |
| `version` | pinned tag, else `latest` | SecureFlow release to install, e.g. `v1.3.0`. |
| `args` | — | Extra arguments for the command. |
| `verify-provenance` | `true` | Verify the binary's GitHub build-provenance attestation (v1.3.0+), on top of the SHA-256 checksum. |

Outputs: `version` (installed version) and `path` (path to the binary).

> **Why an action and not a reusable workflow?** A reusable workflow (`workflow_call`) runs as a separate job on a separate machine, so the decrypted files would not exist in your build job. The action runs as a step inside your job.

Pin to a release tag (`@v1.3.0`) or a commit SHA rather than `@main` so a change on `main` can't alter your pipeline.

### Method 1: Local Installation

The easiest approach for CI/CD is to use `install-local` which creates a portable installation.

**Approach A: Download and install during each CI/CD run**

```yaml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      # Download and setup SecureFlow (no sudo required)
      - name: Setup SecureFlow
        run: |
          wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
          chmod +x secureflow-linux-amd64
          ./secureflow-linux-amd64 install-local
      
      # Use the launcher script
      - name: Decrypt secrets
        run: |
          ./secureflow.sh decrypt --password "${{ secrets.SECUREFLOW_PASSWORD }}" --non-interactive
      
      # Build and deploy
      - name: Build & Deploy
        run: |
          npm run build
          ./deploy.sh
```

**Approach B: Commit the local installation to your repo (Recommended for simplicity)**

Run `secureflow install-local` once locally, then commit the files:

```bash
# Run locally (one-time setup)
secureflow install-local
git add .secureflow/ secureflow.sh
git commit -m "Add portable SecureFlow installation"
git push
```

Your CI/CD workflow becomes very simple (no installation step needed):

```yaml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      # No installation needed - files are already in the repo!
      - name: Decrypt secrets
        run: |
          ./secureflow.sh decrypt --password "${{ secrets.SECUREFLOW_PASSWORD }}" --non-interactive
      
      # Build and deploy
      - name: Build & Deploy
        run: |
          npm run build
          ./deploy.sh
```

### Method 2: Global Installation

```yaml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      # Install SecureFlow globally
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

- **Algorithm**: AES-256-GCM (authenticated encryption: tampering or a wrong password is always detected)
- **Key Derivation**: PBKDF2-HMAC-SHA256, 600,000 iterations, 16-byte random salt per file
- **Format**: `SFENC` header (version, iterations, salt, nonce) + ciphertext + 16-byte auth tag; the header is authenticated too
- **Legacy files**: files from SecureFlow ≤ 1.2 (OpenSSL-compatible AES-256-CBC) still decrypt. Run `secureflow refresh` to upgrade them.

> **Note:** since 1.3, newly encrypted files are no longer OpenSSL-compatible. GCM with a SecureFlow header replaces the unauthenticated CBC format, which could be silently modified.

### Security Model

- **Passwords are never stored or logged**
- **Authenticated encryption** (AES-256-GCM) with a slow key derivation (PBKDF2, 600k iterations)
- **Decrypted files are owner-only** (`0600`), and every write is atomic
- **Non-interactive mode** for secure CI/CD integration (`SECUREFLOW_PASSWORD` keeps the password off the command line)
- **Verified downloads**: release binaries ship with SHA-256 checksums and GitHub build-provenance attestations

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

## 📊 Encryption Manifest

`encrypt`, `add` and `refresh` maintain `manifest.json` in your output directory (commit it alongside the encrypted files):

```json
{
  "schema_version": 1,
  "note": "Encrypted secrets for CI/CD",
  "password_hint": "For the wise only",
  "created_at": "2026-11-27T14:39:31Z",
  "updated_at": "2026-11-27T15:02:10Z",
  "files": [
    {
      "input": "android/app/keystore.jks",
      "encrypted": "android-app-keystore.jks.encrypted",
      "format": "v2-aes-256-gcm",
      "size_bytes": 2754,
      "lines": 12,
      "source_modified": "2026-11-27T14:39:30Z",
      "encrypted_at": "2026-11-27T14:39:31Z",
      "encrypted_sha256": "5e0c3cab4af6d36a42d178037909b26989565937d7f9791a6930674584b78d1f"
    }
  ]
}
```

- The **password hint** is shown by `decrypt`, `test` and `refresh` before asking for the password.
- Nothing derived from the plaintext (such as a hash) is stored, since that would let someone check guesses offline. `encrypted_sha256` is the hash of the *encrypted* file.
- A `report.txt` from SecureFlow ≤ 1.2 is still read for its hint, and replaced by `manifest.json` on the next write.

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

**Option 1: Using install-local (recommended)**

```bash
# One-time setup
secureflow install-local

# In your CI/CD pipeline
./secureflow.sh decrypt --password ${{ secrets.SECUREFLOW_PASSWORD }} --non-interactive
```

**Option 2: Global installation**

```yaml
- name: Decrypt secrets
  run: secureflow decrypt --password ${{ secrets.SECUREFLOW_PASSWORD }} --non-interactive
```

### Portable CI/CD Setup

For teams that want a fully portable solution, commit the local installation to your repo:

```bash
# One-time setup in your repository
secureflow install-local
git add .secureflow/ secureflow.sh
git commit -m "Add portable SecureFlow installation"

# Your CI/CD config becomes very simple:
# - name: Decrypt secrets
#   run: ./secureflow.sh decrypt --password "$PASSWORD" --non-interactive
```

This approach:
- ✅ No installation step needed in CI/CD
- ✅ Works across all platforms automatically
- ✅ Consistent versions across all environments
- ✅ Faster CI/CD runs (no download/install time)

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
│   ├── init.go            # Initialize config command
│   ├── install_local.go   # Local installation command
│   └── secureflow.sh      # Launcher script template
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
