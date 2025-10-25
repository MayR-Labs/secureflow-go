# SecureFlow Examples

This directory contains practical examples of using SecureFlow CLI.

## Basic Usage Examples

### 1. Encrypting Environment Variables

```bash
# Initialize configuration
secureflow init

# Encrypt files with interactive mode
secureflow encrypt
# You'll be prompted for:
# - Encryption password
# - Password hint (optional)
# - Note (optional)

# Non-interactive mode (for CI/CD)
secureflow encrypt --password "your_strong_password" --non-interactive
```

### 2. Decrypting for Local Development

```bash
# Interactive mode
secureflow decrypt

# Non-interactive mode
secureflow decrypt --password "your_strong_password" --non-interactive
```

### 3. Testing Decryption

Test decryption without overwriting your working files:

```bash
# Interactive mode
secureflow test

# Non-interactive mode
secureflow test --password "your_password" --non-interactive
```

## Advanced Examples

### Custom Configuration File

Create a custom `secureflow.yaml`:

```yaml
output_dir: encrypted_secrets
test_output_dir: test_decrypted

files:
  - input: .env.production
    output: .env.production.encrypted
  - input: database/credentials.json
    output: credentials.json.encrypted
  - input: ssl/private.key
    output: private.key.encrypted
```

Use it:

```bash
secureflow encrypt --config /path/to/custom/secureflow.yaml
secureflow decrypt --config /path/to/custom/secureflow.yaml
```

### GitHub Actions Integration

`.github/workflows/deploy.yml`:

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
      
      - name: Download SecureFlow
        run: |
          wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
          chmod +x secureflow-linux-amd64
          sudo mv secureflow-linux-amd64 /usr/local/bin/secureflow
      
      - name: Decrypt secrets
        run: |
          secureflow decrypt --password "${{ secrets.SECUREFLOW_PASSWORD }}" --non-interactive
        
      - name: Deploy application
        run: |
          # Your deployment commands here
          echo "Deploying with decrypted secrets..."
```

### GitLab CI Integration

`.gitlab-ci.yml`:

```yaml
stages:
  - decrypt
  - deploy

decrypt_secrets:
  stage: decrypt
  before_script:
    - wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
    - chmod +x secureflow-linux-amd64
    - mv secureflow-linux-amd64 /usr/local/bin/secureflow
  script:
    - secureflow decrypt --password "$SECUREFLOW_PASSWORD" --non-interactive
  artifacts:
    paths:
      - .env.prod
      - android/app/keystore.jks
    expire_in: 1 hour

deploy:
  stage: deploy
  dependencies:
    - decrypt_secrets
  script:
    - echo "Deploying with decrypted secrets..."
```

### Bitbucket Pipelines Integration

`bitbucket-pipelines.yml`:

```yaml
pipelines:
  default:
    - step:
        name: Decrypt and Deploy
        script:
          - wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
          - chmod +x secureflow-linux-amd64
          - mv secureflow-linux-amd64 /usr/local/bin/secureflow
          - secureflow decrypt --password "$SECUREFLOW_PASSWORD" --non-interactive
          - echo "Deploying..."
```

## Example Workflow

### Setting Up a New Project

1. **Initialize SecureFlow**:
```bash
cd /path/to/your/project
secureflow init
```

2. **Edit the generated `secureflow.yaml`** to match your files:
```yaml
output_dir: enc_keys
test_output_dir: test_dec_keys

files:
  - input: .env.production
    output: .env.production.encrypted
  - input: config/database.yml
    output: database.yml.encrypted
```

3. **Encrypt your sensitive files**:
```bash
secureflow encrypt
# Enter a strong password and optional hint
```

4. **Commit encrypted files to version control**:
```bash
git add enc_keys/
git add secureflow.yaml
git commit -m "Add encrypted secrets"
git push
```

5. **Add `.env.production` and other sensitive files to `.gitignore`**:
```
# .gitignore
.env.production
config/database.yml
test_dec_keys/
```

### Local Development Setup

When another developer clones the repo:

1. **Clone repository**:
```bash
git clone https://github.com/yourorg/your-project.git
cd your-project
```

2. **Install SecureFlow** (if not already installed):
```bash
# On Linux/macOS
wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
chmod +x secureflow-linux-amd64
sudo mv secureflow-linux-amd64 /usr/local/bin/secureflow
```

3. **Decrypt secrets** (obtain password from team):
```bash
secureflow decrypt
# Enter the password provided by your team
```

### CI/CD Setup

1. **Store password as secret** in your CI/CD platform:
   - GitHub Actions: Settings → Secrets → Actions → New repository secret
   - GitLab CI: Settings → CI/CD → Variables
   - Bitbucket: Repository settings → Pipelines → Repository variables

2. **Add SecureFlow to your pipeline** (see integration examples above)

3. **Use decrypted files** in subsequent pipeline steps

## Example Project Structure

```
your-project/
├── .gitignore
├── secureflow.yaml
├── enc_keys/                    # Encrypted files (committed to git)
│   ├── .env.production.encrypted
│   ├── keystore.jks.encrypted
│   └── report.txt               # Encryption report
├── test_dec_keys/               # Test decryption output (not in git)
│   └── .env.production
├── .env.production              # Original secret file (not in git)
├── android/
│   └── app/
│       └── keystore.jks        # Original secret file (not in git)
└── src/
    └── ...
```

## Security Best Practices

1. **Use strong passwords**: At least 16 characters with mixed case, numbers, and symbols
2. **Rotate passwords periodically**: Re-encrypt with new passwords every few months
3. **Limit password access**: Share passwords securely (e.g., using a password manager)
4. **Never commit plaintext secrets**: Always add them to `.gitignore`
5. **Use different passwords**: For different projects or environments
6. **Store CI/CD secrets securely**: Use your platform's secret management features
7. **Review the report.txt**: Check what's encrypted and when

## Troubleshooting

### Wrong Password Error

```bash
$ secureflow decrypt
Error: decryption failed (wrong password?)
```

**Solution**: Verify you're using the correct password. Check `enc_keys/report.txt` for password hints.

### File Not Found During Encryption

```bash
⚠️  Warning: .env.production not found, skipping
```

**Solution**: Make sure the file exists before encrypting, or remove it from `secureflow.yaml`.

### Permission Denied

```bash
Error: failed to create directory: permission denied
```

**Solution**: Ensure you have write permissions for the output directory, or run with appropriate permissions.

## Additional Resources

- [Main README](../README.md)
- [GitHub Repository](https://github.com/MayR-Labs/secureflow-go)
- [Latest Releases](https://github.com/MayR-Labs/secureflow-go/releases)
