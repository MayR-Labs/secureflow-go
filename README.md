# 🧩 SecureFlow CLI

**SecureFlow** is a lightweight, Go-based CLI for securely encrypting and decrypting sensitive files like environment variables, keystores, and service credentials for local and CI/CD use.

It’s designed to replace fragile Bash scripts with a fast, cross-platform executable that works seamlessly across Linux, macOS, and Windows.

---

## ⚙️ Features

* AES-256 encryption (using OpenSSL-compatible standards)
* Interactive and non-interactive modes
* Automatically generates and uses a `secureflow.yaml` configuration file
* Encrypted file metadata reports
* Safe testing mode (`secureflow test`)
* Clean error handling and consistent output

---

## 🚀 Installation

SecureFlow can be installed **locally** (recommended) or **globally**.

### Local Installation (Recommended)

Install SecureFlow in your project directory. This allows team members and CI/CD pipelines to use the executable without needing global installation.

#### Quick Install (Linux/macOS)

```bash
curl -sSL https://raw.githubusercontent.com/MayR-Labs/secureflow-go/main/install.sh | bash
```

Or with wget:

```bash
wget -qO- https://raw.githubusercontent.com/MayR-Labs/secureflow-go/main/install.sh | bash
```

This downloads the `secureflow` binary to your current directory. 

**Should you commit the binary to your repository?**

- **✅ Recommended:** Commit the binary for easy team onboarding and CI/CD. Team members can use it immediately after cloning.
- **Alternative:** Add it to `.gitignore` if you prefer each developer to install it separately.

You can then:

1. Run it with `./secureflow`
2. Commit it to your repository: `git add secureflow && git commit -m "Add SecureFlow binary"`
3. Or add it to `.gitignore` if you prefer: `echo "secureflow" >> .gitignore`

#### Manual Local Installation

Download the precompiled binary for your OS from the [Releases](https://github.com/MayR-Labs/secureflow-go/releases) page.

**Linux (AMD64):**
```bash
wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
chmod +x secureflow-linux-amd64
mv secureflow-linux-amd64 secureflow
```

**macOS (Intel):**
```bash
wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-darwin-amd64
chmod +x secureflow-darwin-amd64
mv secureflow-darwin-amd64 secureflow
```

**macOS (Apple Silicon):**
```bash
wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-darwin-arm64
chmod +x secureflow-darwin-arm64
mv secureflow-darwin-arm64 secureflow
```

**Windows:**
Download `secureflow-windows-amd64.exe` from the [Releases](https://github.com/MayR-Labs/secureflow-go/releases) page and rename it to `secureflow.exe` in your project directory.

### Global Installation

If you prefer to install SecureFlow globally (available system-wide):

```bash
curl -sSL https://raw.githubusercontent.com/MayR-Labs/secureflow-go/main/install.sh | bash -s -- --global
```

Or manually:

**Linux (AMD64):**
```bash
wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
chmod +x secureflow-linux-amd64
sudo mv secureflow-linux-amd64 /usr/local/bin/secureflow
```

**macOS (Intel):**
```bash
wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-darwin-amd64
chmod +x secureflow-darwin-amd64
sudo mv secureflow-darwin-amd64 /usr/local/bin/secureflow
```

**macOS (Apple Silicon):**
```bash
wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-darwin-arm64
chmod +x secureflow-darwin-arm64
sudo mv secureflow-darwin-arm64 /usr/local/bin/secureflow
```

**Windows:**
Download `secureflow-windows-amd64.exe` from the [Releases](https://github.com/MayR-Labs/secureflow-go/releases) page and add it to your PATH.

### From Source

```bash
git clone https://github.com/MayR-Labs/secureflow-go.git
cd secureflow-go
go build -o secureflow
# For global install: sudo mv secureflow /usr/local/bin/
```

### Verify Installation

**Local installation:**
```bash
./secureflow --version
```

**Global installation:**
```bash
secureflow --version
```

---

## 🧰 Usage

> **Note:** If you installed SecureFlow locally, prefix all commands with `./` (e.g., `./secureflow init`). For global installations, use `secureflow` directly.

### 1. Initialise a new config

**Local installation:**
```bash
./secureflow init
```

**Global installation:**
```bash
secureflow init
```

Creates a default `secureflow.yaml` file in your current directory:

```yaml
# secureflow.yaml
output_dir: enc_keys
test_output_dir: test_dec_keys

files:
  - input: .env.prod
    output: .env.prod.encrypted
  - input: android/app/keystore.jks
    output: keystore.jks.encrypted
  - input: android/key.properties
    output: key.properties.encrypted
  - input: android/service-key.json
    output: service-key.json.encrypted
```

You can modify this file to fit your project structure.

---

### 2. Encrypt Files

```bash
secureflow encrypt
```

You’ll be prompted for an encryption password, password hint, and optional note.

For non-interactive mode (CI/CD):

```bash
secureflow encrypt --password "your_password" --non-interactive
```

To specify a custom config file:

```bash
secureflow encrypt --config ./path/to/secureflow.yaml
```

All encrypted files will be saved to the directory specified in the YAML file (default: `enc_keys`).

---

### 3. Decrypt Files

Decrypt files for local development or CI pipelines:

**Local installation:**
```bash
./secureflow decrypt --password "your_password"
```

**Global installation:**
```bash
secureflow decrypt --password "your_password"
```

For non-interactive mode:

```bash
./secureflow decrypt --password "$ENCRYPTION_PASSWORD" --non-interactive
```

To use a custom config:

```bash
./secureflow decrypt --config ./custom/secureflow.yaml
```

---

### 4. Test Decryption

This mode decrypts files into a separate test directory without overwriting existing secrets.

**Local installation:**
```bash
./secureflow test
```

**Global installation:**
```bash
secureflow test
```

Non-interactive version:

```bash
./secureflow test --password "your_password" --non-interactive
```

---

## 🧾 Report File

Each encryption run generates a detailed `report.txt` inside the output directory.

Example:

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

---

## 💡 Common Use Cases

### CI/CD Decryption Step

#### With Local Executable (Recommended)

If you commit the `secureflow` binary to your repository, CI/CD pipelines can use it directly:

**GitHub Actions:**
```yaml
- name: Make secureflow executable
  run: chmod +x ./secureflow

- name: Decrypt secrets
  run: ./secureflow decrypt --password ${{ secrets.SECUREFLOW_PASSWORD }} --non-interactive
```

**GitLab CI:**
```yaml
decrypt_secrets:
  script:
    - chmod +x ./secureflow
    - ./secureflow decrypt --password "$SECUREFLOW_PASSWORD" --non-interactive
```

#### Without Committed Binary

If you prefer not to commit the binary, download it during CI/CD:

**GitHub Actions:**
```yaml
- name: Download SecureFlow
  run: |
    wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
    chmod +x secureflow-linux-amd64
    mv secureflow-linux-amd64 secureflow

- name: Decrypt secrets
  run: ./secureflow decrypt --password ${{ secrets.SECUREFLOW_PASSWORD }} --non-interactive
```

**GitLab CI:**
```yaml
before_script:
  - wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
  - chmod +x secureflow-linux-amd64
  - mv secureflow-linux-amd64 secureflow

decrypt_secrets:
  script:
    - ./secureflow decrypt --password "$SECUREFLOW_PASSWORD" --non-interactive
```

### Local Development Workflow

**Initial setup:**
```bash
# Install SecureFlow locally in your project
curl -sSL https://raw.githubusercontent.com/MayR-Labs/secureflow-go/main/install.sh | bash

# Initialize configuration
./secureflow init

# Encrypt your sensitive files
./secureflow encrypt
```

**Team member cloning the repository:**
```bash
git clone https://github.com/yourorg/your-project.git
cd your-project

# If secureflow is committed, just decrypt
./secureflow decrypt

# If secureflow is not committed, install it first
curl -sSL https://raw.githubusercontent.com/MayR-Labs/secureflow-go/main/install.sh | bash
./secureflow decrypt
```

### Local Encryption

```bash
./secureflow encrypt
```

This helps ensure your sensitive files never end up in plaintext in version control.

---

## 🧱 Error Handling

* **Missing files** → Skipped with warning, process continues
* **Wrong password** → Graceful failure with exit code `1`
* **Malformed YAML** → Detailed message showing offending line
* **Output directory missing** → Automatically created
* **Interrupts (Ctrl+C)** → Gracefully exits with cleanup notice

---

## 🔐 Security Model

* Encryption: AES-256-CBC with PBKDF2 key derivation (OpenSSL-compatible)
* Passwords are never stored or logged
* Non-interactive password injection supported for CI/CD pipelines
* Compatible with encrypted artefacts from the earlier Bash version

---

## 🧩 Project Structure (Go)

```
secureflow/
│
├── cmd/
│   ├── root.go          # CLI entrypoint (Cobra)
│   ├── encrypt.go
│   ├── decrypt.go
│   ├── test.go
│   └── init.go
│
├── internal/
│   ├── crypto/          # Encryption/decryption logic
│   ├── config/          # secureflow.yaml handling
│   └── utils/           # File handling, error logging
│
├── go.mod
├── main.go
├── LICENSE                 # MIT License
├── install.sh              # Installation script
├── examples/               # Usage examples and configs
└── README.md
```

---

## 📚 Examples

Check out the [examples directory](./examples/README.md) for:
- Practical usage examples
- CI/CD integration guides (GitHub Actions, GitLab CI, Bitbucket Pipelines)
- Sample configuration files for different use cases
- Best practices and troubleshooting

---

## 🧪 Development

Run locally without installing:

```bash
go run main.go encrypt
```

Run all tests:

```bash
go test ./...
```

Build binary:

```bash
go build -o secureflow
```

---

## 🧭 Future Roadmap

* Support for `.env` key filtering (only encrypt certain variables)
* Optional GPG-based encryption backend
* Integration with Flutter build runners
* Progress bars for large files

---

## 📄 License

SecureFlow is licensed under the [MIT License](LICENSE).

Copyright (c) 2025 [MayR Labs](https://mayrlabs.com)

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## 🔗 Links

- [GitHub Repository](https://github.com/MayR-Labs/secureflow-go)
- [Latest Releases](https://github.com/MayR-Labs/secureflow-go/releases)
- [MayR Labs](https://mayrlabs.com)
- [Report Issues](https://github.com/MayR-Labs/secureflow-go/issues)

---
