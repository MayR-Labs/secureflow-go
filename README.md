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

### From Source

```bash
git clone https://github.com/YoungMayor/secureflow-go.git
cd secureflow
go build -o secureflow
```

### From Release (recommended)

Download the precompiled binary for your OS from the [Releases](https://github.com/YoungMayor/secureflow-go/releases) page,
then move it to a directory in your PATH, for example:

```bash
sudo mv secureflow /usr/local/bin/
```

Check installation:

```bash
secureflow --version
```

---

## 🧰 Usage

### 1. Initialise a new config

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

```bash
secureflow decrypt --password "your_password"
```

For non-interactive mode:

```bash
secureflow decrypt --password "$ENCRYPTION_PASSWORD" --non-interactive
```

To use a custom config:

```bash
secureflow decrypt --config ./custom/secureflow.yaml
```

---

### 4. Test Decryption

This mode decrypts files into a separate test directory without overwriting existing secrets.

```bash
secureflow test
```

Non-interactive version:

```bash
secureflow test --password "your_password" --non-interactive
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

In GitHub Actions, for example:

```yaml
- name: Decrypt secrets
  run: secureflow decrypt --password ${{ secrets.SECUREFLOW_PASSWORD }} --non-interactive
```

### Local Encryption

```bash
secureflow encrypt
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
└── README.md
```

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
