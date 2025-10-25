# Security Guide

This guide covers security best practices for using SecureFlow to protect your sensitive data.

## Table of Contents

- [Encryption Overview](#encryption-overview)
- [Password Security](#password-security)
- [File Management](#file-management)
- [CI/CD Security](#cicd-security)
- [Best Practices](#best-practices)
- [Security Checklist](#security-checklist)

## Encryption Overview

### Encryption Algorithm

SecureFlow uses industry-standard encryption:

- **Algorithm**: AES-256-CBC (Advanced Encryption Standard, 256-bit key, Cipher Block Chaining mode)
- **Key Derivation**: PBKDF2 (Password-Based Key Derivation Function 2)
- **Compatibility**: OpenSSL-compatible format
- **Salt**: 8 random bytes per encrypted file
- **Format**: Salted__ header + salt + encrypted data

### How It Works

1. **Encryption Process**:
   ```
   Password → PBKDF2 → Key + IV → AES-256-CBC → Encrypted File
   ```

2. **File Format**:
   ```
   [Salted__][8-byte salt][encrypted data]
   ```

3. **Key Derivation**:
   - Uses PBKDF2 with SHA-256
   - Derives 48 bytes: 32 for key, 16 for IV
   - Salt is randomly generated per file

### OpenSSL Compatibility

SecureFlow is compatible with OpenSSL commands:

**Encrypt with OpenSSL**:
```bash
openssl enc -aes-256-cbc -salt -pbkdf2 -in file.txt -out file.txt.encrypted -k "password"
```

**Decrypt with SecureFlow**:
```bash
secureflow decrypt --password "password"
```

And vice versa - files encrypted with SecureFlow can be decrypted with OpenSSL.

## Password Security

### Password Requirements

**Minimum Requirements**:
- At least 16 characters
- Mix of uppercase and lowercase letters
- Include numbers
- Include special characters

**Good Password Example**:
```
T7$kL9@mP2!qR5#nW8
```

**Bad Password Examples**:
```
password123          # Too simple
MyP@ssw0rd          # Common pattern
project2024         # Too short, predictable
```

### Password Generation

Use a password manager or generator:

```bash
# Generate strong password on Linux/macOS
openssl rand -base64 24

# Generate with special characters
< /dev/urandom tr -dc 'A-Za-z0-9!@#$%^&*' | head -c 24 ; echo
```

### Password Storage

**✅ DO**:
- Use a password manager (1Password, LastPass, Bitwarden)
- Store in CI/CD platform's secret management
- Share via secure channels (encrypted messaging)
- Keep separate passwords for different projects
- Document password hints in `report.txt`

**❌ DON'T**:
- Commit passwords to version control
- Share passwords via email or Slack
- Write passwords in plain text files
- Reuse passwords across projects
- Store passwords in browser history

### Password Rotation

**When to Rotate**:
- Every 3-6 months (recommended)
- When team member leaves
- After suspected compromise
- Before major releases

**How to Rotate**:
1. Generate new password
2. Re-encrypt all files with new password:
   ```bash
   secureflow encrypt --password "new_password"
   ```
3. Update password in CI/CD secrets
4. Commit new encrypted files
5. Notify team of password change
6. Update password manager

## File Management

### Protecting Source Files

**Always add plaintext secrets to `.gitignore`**:

```gitignore
# Environment files
.env
.env.*
!.env.example

# Mobile app secrets
android/app/keystore.jks
android/key.properties
android/service-account.json
ios/Runner/GoogleService-Info.plist

# SSL certificates
*.key
*.pem
*.p12

# Database credentials
config/database.yml
config/secrets.yml

# Test decryption directory
test_dec_keys/
```

### Encrypted File Management

**✅ DO**:
- Commit encrypted files to version control
- Commit `secureflow.yaml` configuration
- Commit `report.txt` for reference
- Use descriptive encrypted filenames
- Review encryption reports regularly

**❌ DON'T**:
- Commit plaintext secrets
- Delete encrypted files accidentally
- Modify encrypted files manually
- Share encrypted files without config

### Directory Structure

**Recommended structure**:
```
your-project/
├── .gitignore                    # Ignore plaintext secrets
├── secureflow.yaml              # Config (commit this)
├── enc_keys/                    # Encrypted files (commit this)
│   ├── .env.prod.encrypted
│   ├── keystore.jks.encrypted
│   └── report.txt               # Report (commit this)
├── test_dec_keys/               # Test directory (don't commit)
│   └── .env.prod
├── .env.prod                    # Plaintext (don't commit)
└── android/
    └── app/
        └── keystore.jks         # Plaintext (don't commit)
```

### File Permissions

**On Linux/macOS**:
```bash
# Encrypted files (can be readable)
chmod 644 enc_keys/*.encrypted

# Decrypted sensitive files (restrict access)
chmod 600 .env.prod
chmod 600 android/app/keystore.jks

# Directories
chmod 755 enc_keys
chmod 700 test_dec_keys
```

**On Windows**:
Use File Explorer → Properties → Security to restrict access to decrypted files.

## CI/CD Security

### Secret Management

**Platform-Specific Best Practices**:

#### GitHub Actions
- Use **Repository secrets** for passwords
- Use **Environment secrets** for deployment-specific passwords
- Enable **branch protection** for production branches
- Use **required reviewers** for sensitive workflows
- Enable **CODEOWNERS** for workflow files

#### GitLab CI
- Use **masked variables** to hide values in logs
- Use **protected variables** for production branches only
- Set **environment-specific variables**
- Limit variable scope to specific branches/tags
- Use **file variables** for multi-line secrets

#### Bitbucket Pipelines
- Use **secured variables** to mask in logs
- Limit variables to specific branches
- Use **deployment variables** for environment-specific secrets
- Enable **required approvals** for production deploys

#### Jenkins
- Use **Credentials** store for passwords
- Scope credentials to specific jobs/folders
- Enable **Mask Passwords** plugin
- Use **Role-Based Access Control** (RBAC)
- Audit credential access regularly

### Pipeline Security

**Secure Pipeline Design**:

```yaml
# Example: Secure GitHub Actions workflow
name: Secure Deploy

on:
  push:
    branches: [main]

permissions:
  contents: read    # Minimal permissions

jobs:
  decrypt:
    runs-on: ubuntu-latest
    environment: production  # Requires manual approval
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Decrypt
        run: |
          secureflow decrypt --password "${{ secrets.PROD_PASSWORD }}" --non-interactive
        
      - name: Upload artifacts with short expiry
        uses: actions/upload-artifact@v3
        with:
          name: secrets
          path: .env.production
          retention-days: 1  # Short retention
  
  deploy:
    needs: decrypt
    runs-on: ubuntu-latest
    
    steps:
      - name: Download secrets
        uses: actions/download-artifact@v3
      
      - name: Deploy
        run: ./deploy.sh
      
      - name: Clean up
        if: always()
        run: rm -f .env.production
```

**Key Security Practices**:
1. Use minimal permissions
2. Require manual approval for production
3. Short artifact retention
4. Clean up secrets after use
5. Use specific branch triggers

### Audit and Monitoring

**What to Monitor**:
- Secret access in CI/CD platform
- Pipeline execution history
- Failed decryption attempts
- Unauthorized config changes
- Encrypted file modifications

**Audit Checklist**:
- [ ] Review who has access to secrets
- [ ] Check pipeline logs for anomalies
- [ ] Verify encrypted files haven't been tampered with
- [ ] Ensure `.gitignore` is properly configured
- [ ] Review recent commits for accidentally committed secrets

## Best Practices

### 1. Defense in Depth

Don't rely solely on encryption:

- Use VPN for sensitive operations
- Enable 2FA on all accounts
- Use network segmentation
- Implement least privilege access
- Regular security audits

### 2. Separate Environments

Use different passwords for different environments:

```bash
# Development
secureflow encrypt --config secureflow.dev.yaml --password "$DEV_PASSWORD"

# Staging
secureflow encrypt --config secureflow.staging.yaml --password "$STAGING_PASSWORD"

# Production
secureflow encrypt --config secureflow.prod.yaml --password "$PROD_PASSWORD"
```

### 3. Minimize Secret Scope

**Only encrypt what's necessary**:

```yaml
# Good: Only encrypt secrets
files:
  - input: .env.production
  - input: android/app/keystore.jks

# Avoid: Don't encrypt public configs
files:
  - input: .env.production
  - input: config/public_api.json      # Not necessary
```

### 4. Regular Security Reviews

**Monthly checklist**:
- [ ] Review who has access to passwords
- [ ] Check for exposed secrets in logs
- [ ] Verify `.gitignore` is complete
- [ ] Test decryption process
- [ ] Review encryption reports

**Quarterly checklist**:
- [ ] Rotate passwords
- [ ] Re-encrypt all files
- [ ] Update team access
- [ ] Review pipeline security
- [ ] Audit secret usage

### 5. Incident Response Plan

**If password is compromised**:

1. **Immediate Actions**:
   ```bash
   # Generate new password
   NEW_PASS=$(openssl rand -base64 24)
   
   # Re-encrypt immediately
   secureflow encrypt --password "$NEW_PASS"
   
   # Update CI/CD secrets
   # Update team password manager
   ```

2. **Investigate**:
   - Review access logs
   - Check for unauthorized deployments
   - Examine recent commits
   - Interview team members

3. **Communicate**:
   - Notify affected team members
   - Document the incident
   - Update security procedures
   - Consider external notification if required

**If plaintext secrets are committed**:

1. **Immediate Actions**:
   ```bash
   # Remove from git history (use with caution)
   git filter-branch --force --index-filter \
     'git rm --cached --ignore-unmatch path/to/secret.env' \
     --prune-empty --tag-name-filter cat -- --all
   
   # Force push (coordinate with team)
   git push origin --force --all
   ```

2. **Rotate All Secrets**:
   - Change all passwords
   - Regenerate API keys
   - Update database credentials
   - Re-issue certificates

3. **Use Tools to Prevent Future Incidents**:
   ```bash
   # Install git-secrets
   git secrets --install
   git secrets --register-aws
   
   # Or use gitleaks
   gitleaks detect --verbose
   ```

### 6. Team Training

**Train team members on**:
- How to use SecureFlow properly
- Never committing plaintext secrets
- Recognizing security threats
- Proper password management
- Incident response procedures

### 7. Verification Testing

**Regular verification**:
```bash
# Test encryption
secureflow test --password "$PASSWORD" --non-interactive

# Verify files can be decrypted
secureflow decrypt --password "$PASSWORD" --non-interactive

# Check encrypted file integrity
ls -lh enc_keys/
cat enc_keys/report.txt
```

### 8. Backup Strategy

**Backup important files**:
```bash
# Backup encrypted files
tar -czf secureflow-backup-$(date +%Y%m%d).tar.gz enc_keys/ secureflow.yaml

# Store in secure location
# - Encrypted cloud storage
# - Secure file server
# - Password-protected archive
```

**What to backup**:
- Encrypted files (`enc_keys/`)
- Configuration (`secureflow.yaml`)
- Encryption reports (`report.txt`)
- Password hints (in secure location)

## Security Checklist

### Initial Setup
- [ ] Generate strong password (16+ characters)
- [ ] Store password in password manager
- [ ] Configure `.gitignore` for plaintext secrets
- [ ] Test encryption/decryption locally
- [ ] Set appropriate file permissions
- [ ] Document password hint in report

### Before First Commit
- [ ] Verify no plaintext secrets in git history
- [ ] Check `.gitignore` is working
- [ ] Confirm encrypted files are in place
- [ ] Review `report.txt` for accuracy
- [ ] Test decryption process

### CI/CD Setup
- [ ] Store password as secret in CI/CD platform
- [ ] Test decryption in pipeline
- [ ] Verify password is masked in logs
- [ ] Set short artifact retention
- [ ] Enable branch protection
- [ ] Require approval for production deploys

### Ongoing Maintenance
- [ ] Rotate passwords quarterly
- [ ] Review access permissions monthly
- [ ] Audit pipeline logs regularly
- [ ] Keep SecureFlow updated
- [ ] Test disaster recovery procedures
- [ ] Train new team members

### Before Deploying to Production
- [ ] Verify encryption password is correct
- [ ] Test decryption process
- [ ] Check all required files are encrypted
- [ ] Review CI/CD pipeline security
- [ ] Confirm rollback procedure
- [ ] Document emergency contacts

## Compliance Considerations

### GDPR (General Data Protection Regulation)
- Use strong encryption (AES-256 ✅)
- Maintain data processing records
- Implement access controls
- Enable audit logging
- Have data breach response plan

### HIPAA (Health Insurance Portability and Accountability Act)
- Encrypt data at rest (✅)
- Implement access controls (✅)
- Maintain audit trails
- Regular security assessments
- Business associate agreements

### SOC 2
- Encryption of sensitive data (✅)
- Access control policies (✅)
- Change management
- Monitoring and logging
- Incident response

### PCI DSS (Payment Card Industry Data Security Standard)
- Strong cryptography (AES-256 ✅)
- Key management procedures
- Access control measures
- Logging and monitoring
- Regular security testing

## Additional Security Resources

### Tools
- **git-secrets**: Prevent committing secrets
- **gitleaks**: Detect secrets in git repos
- **truffleHog**: Find secrets in git history
- **detect-secrets**: Prevent secrets in commits

### Reading
- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [NIST Cybersecurity Framework](https://www.nist.gov/cyberframework)
- [CIS Controls](https://www.cisecurity.org/controls)

## See Also

- [Configuration Guide](./configuration.md) - Configuration best practices
- [CI/CD Usage Guide](./cicd-usage.md) - Secure CI/CD integration
- [Troubleshooting](./troubleshooting.md) - Security-related issues
