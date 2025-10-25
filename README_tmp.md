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

This downloads the `secureflow` binary to your current directory. You can then:

1. Run it with `./secureflow`
2. Commit it to your repository (optional but recommended for team workflows)
3. Add it to `.gitignore` if you prefer team members to install it individually

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

