# Configuration Guide

This guide provides comprehensive information about configuring SecureFlow for your project.

## Table of Contents

- [Configuration File Structure](#configuration-file-structure)
- [Configuration Options](#configuration-options)
- [Example Configurations](#example-configurations)
- [Best Practices](#best-practices)

## Configuration File Structure

SecureFlow uses a YAML configuration file (default: `secureflow.yaml`) to define which files to encrypt/decrypt and where to store them.

### Basic Structure

```yaml
output_dir: enc_keys           # Directory for encrypted files
test_output_dir: test_dec_keys # Directory for test decryption

files:
  - input: .env.prod           # Source file path
    output: .env.prod.encrypted # Encrypted filename
  - input: path/to/secret.key
    output: secret.key.encrypted
```

## Configuration Options

### Global Options

#### `output_dir`
- **Type**: String
- **Required**: Yes
- **Default**: `enc_keys`
- **Description**: Directory where encrypted files will be stored
- **Example**: `output_dir: encrypted_secrets`

#### `test_output_dir`
- **Type**: String
- **Required**: Yes
- **Default**: `test_dec_keys`
- **Description**: Directory for test decryption output (used with `secureflow test` command)
- **Example**: `test_output_dir: test_decrypted`

### File Entries

Each file entry in the `files` array requires the `input` and `output` fields, and optionally supports the `copy_to` field:

#### `input`
- **Type**: String
- **Required**: Yes
- **Description**: Path to the source file to encrypt (relative to project root)
- **Example**: `input: config/database.yml`

#### `output`
- **Type**: String
- **Required**: Yes
- **Description**: Filename for the encrypted file (stored in `output_dir`)
- **Example**: `output: database.yml.encrypted`

#### `copy_to`
- **Type**: String
- **Required**: No
- **Description**: After decryption, copy the decrypted file to this path. Useful when applications expect `.env` but you store `.env.prod`
- **Example**: `copy_to: .env`

**Note**: The `output` is just a filename, not a path. All encrypted files are stored in the `output_dir`.

## Example Configurations

### Basic Configuration

Minimal setup for encrypting a single environment file:

```yaml
output_dir: enc_keys
test_output_dir: test_dec_keys

files:
  - input: .env.prod
    output: .env.prod.encrypted
    copy_to: .env  # Automatically creates .env after decryption
```

### Using `copy_to` for Environment Files

The `copy_to` field is particularly useful for environment files. Many applications expect a `.env` file, but you might store different versions like `.env.prod`, `.env.staging`, etc.

```yaml
output_dir: enc_keys
test_output_dir: test_dec_keys

files:
  - input: .env.prod
    output: .env.prod.encrypted
    copy_to: .env  # Production env becomes .env

  - input: .env.staging
    output: .env.staging.encrypted
    # No copy_to - staging env stays as .env.staging
```

After running `secureflow decrypt`, you'll have:
- `.env.prod` (decrypted)
- `.env` (copy of .env.prod)
- `.env.staging` (decrypted)

### Mobile App (Flutter/React Native)

Configuration for mobile apps with Android and iOS credentials:

```yaml
output_dir: encrypted_secrets
test_output_dir: test_decrypted

files:
  - input: .env.production
    output: .env.production.encrypted
    copy_to: .env
  - input: android/app/keystore.jks
    output: keystore.jks.encrypted
  - input: android/key.properties
    output: key.properties.encrypted
  - input: android/service-account.json
    output: service-account.json.encrypted
  - input: ios/Runner/GoogleService-Info.plist
    output: GoogleService-Info.plist.encrypted
```

### Web Application

Configuration for web apps with multiple environments:

```yaml
output_dir: enc_keys
test_output_dir: test_dec_keys

files:
  - input: .env.production
    output: .env.production.encrypted
    copy_to: .env
  - input: .env.staging
    output: .env.staging.encrypted
  - input: config/database.yml
    output: database.yml.encrypted
  - input: config/secrets.yml
    output: secrets.yml.encrypted
  - input: ssl/private.key
    output: ssl-private.key.encrypted
  - input: ssl/certificate.crt
    output: ssl-certificate.crt.encrypted
```

### Microservices Architecture

Configuration for projects with multiple services:

```yaml
output_dir: encrypted
test_output_dir: decrypted_test

files:
  # Auth Service
  - input: services/auth/.env.prod
    output: auth-env.prod.encrypted
    copy_to: services/auth/.env
  - input: services/auth/jwt-keys/private.pem
    output: auth-jwt-private.pem.encrypted
  
  # API Service
  - input: services/api/.env.prod
    output: api-env.prod.encrypted
    copy_to: services/api/.env
  - input: services/api/config/database.json
    output: api-database.json.encrypted
  
  # Worker Service
  - input: services/worker/.env.prod
    output: worker-env.prod.encrypted
    copy_to: services/worker/.env
  
  # Shared Secrets
  - input: shared/redis.conf
    output: shared-redis.conf.encrypted
  - input: shared/rabbitmq.json
    output: shared-rabbitmq.json.encrypted
```

### Docker Deployment

Configuration for Docker-based deployments:

```yaml
output_dir: docker/secrets/encrypted
test_output_dir: docker/secrets/test

files:
  - input: docker/.env.production
    output: docker-env.production.encrypted
    copy_to: docker/.env
  - input: docker/compose/.env.db
    output: docker-env.db.encrypted
  - input: docker/nginx/ssl/private.key
    output: nginx-ssl-private.key.encrypted
  - input: docker/nginx/ssl/certificate.crt
    output: nginx-ssl-certificate.crt.encrypted
```

### Kubernetes Secrets

Configuration for Kubernetes secret files:

```yaml
output_dir: k8s/encrypted-secrets
test_output_dir: k8s/test-secrets

files:
  - input: k8s/secrets/database-credentials.yaml
    output: database-credentials.yaml.encrypted
  - input: k8s/secrets/api-keys.yaml
    output: api-keys.yaml.encrypted
  - input: k8s/secrets/tls-cert.yaml
    output: tls-cert.yaml.encrypted
  - input: k8s/config/app-config.json
    output: app-config.json.encrypted
```

### Multi-Environment Setup

Configuration for projects with multiple environments:

```yaml
output_dir: secrets/encrypted
test_output_dir: secrets/test

files:
  # Production
  - input: envs/production/.env
    output: production.env.encrypted
  - input: envs/production/database.yml
    output: production.database.yml.encrypted
  
  # Staging
  - input: envs/staging/.env
    output: staging.env.encrypted
  - input: envs/staging/database.yml
    output: staging.database.yml.encrypted
  
  # Development (shared team secrets)
  - input: envs/development/.env.shared
    output: development.env.shared.encrypted
```

## Best Practices

### 1. Using `copy_to` for Application Compatibility

When your application expects files at specific paths (like `.env`) but you store environment-specific versions (like `.env.prod`), use the `copy_to` field:

```yaml
files:
  - input: .env.prod
    output: .env.prod.encrypted
    copy_to: .env  # App reads from .env
```

**Benefits:**
- Keep your source files clearly named (`.env.prod`, `.env.staging`)
- Application automatically gets the file it expects (`.env`)
- No need to manually copy files after decryption
- Works great in CI/CD pipelines

**Add to .gitignore:**
```
.env       # The copied file (decrypted)
# Note: .env.prod is your encrypted source, commit it to git
```

### 2. Naming Conventions

- **Use descriptive output names**: Make it clear what each encrypted file contains
  ```yaml
  output: prod-database-credentials.encrypted  # Good
  output: file1.encrypted                       # Bad
  ```

- **Include environment in filename**: When encrypting files for different environments
  ```yaml
  output: staging.env.encrypted
  output: production.env.encrypted
  ```

### 3. Directory Organization

- **Keep encrypted files in a dedicated directory**: Use `output_dir` consistently
- **Add encrypted directories to git**: `git add enc_keys/`
- **Ignore decrypted directories**: Add to `.gitignore`:
  ```
  test_dec_keys/
  .env.production
  .env.staging
  .env
  *.key
  *.jks
  ```

### 4. Configuration Management

- **Commit the config file**: Your `secureflow.yaml` should be in version control
- **Use custom configs per environment**: For different deployment scenarios
  ```bash
  secureflow encrypt --config secureflow.production.yaml
  secureflow encrypt --config secureflow.staging.yaml
  ```

### 5. File Organization

- **Group related files**: Keep similar files together in the configuration
- **Use comments**: YAML supports comments to document your configuration
  ```yaml
  files:
    # Database credentials
    - input: .env.db
      output: database.env.encrypted
      copy_to: .env.database
    
    # API keys
    - input: .env.api
      output: api.env.encrypted
  ```

### 6. Security Considerations

- **Never commit plaintext secrets**: Always add sensitive files to `.gitignore`
- **Verify encrypted files**: After encryption, check that:
  - Encrypted files exist in `output_dir`
  - `report.txt` was generated with correct metadata
  - Original files are still in place (encryption doesn't delete them)

### 7. Path Specifications

- **Use relative paths**: All input paths are relative to project root
- **Use forward slashes**: Works on all platforms (Windows, Linux, macOS)
  ```yaml
  input: config/secrets.json      # Good (cross-platform)
  input: config\secrets.json      # Bad (Windows only)
  ```

### 8. Testing Configuration

Always test your configuration before using it in production:

```bash
# Test with the test command
secureflow test --password "test_password"

# Verify test output
ls test_dec_keys/
```

## Using Custom Configuration Files

You can use multiple configuration files for different purposes:

```bash
# Development encryption
secureflow encrypt --config secureflow.dev.yaml

# Production encryption
secureflow encrypt --config secureflow.prod.yaml

# CI/CD encryption
secureflow encrypt --config secureflow.ci.yaml
```

## Initializing Configuration

Create a default configuration file:

```bash
secureflow init
```

### Using Configuration Templates

SecureFlow provides pre-configured templates for common project types. You can either select a template interactively or specify one directly:

**Interactive Selection:**
```bash
secureflow init
# Prompts you to choose from:
# 1. Default (React Native/Mobile App)
# 2. React Native
# 3. Flutter
# 4. Web Application
# 5. Docker Deployment
# 6. Kubernetes (K8s)
# 7. Microservices
```

**Direct Template Selection:**
```bash
secureflow init --template flutter
secureflow init --template reactnative
secureflow init --template web
secureflow init --template docker
secureflow init --template k8s
secureflow init --template microservices
```

**Available Templates:**

- **`default`** - React Native/Mobile app with Android/iOS files and `.env.prod` → `.env` copy
- **`reactnative`** - React Native specific with staging and production environments
- **`flutter`** - Flutter mobile app with Android and iOS configurations
- **`web`** - Web application with multiple environment files and database/secrets configurations
- **`docker`** - Docker deployment with compose and nginx configurations
- **`k8s`** (or `kubernetes`) - Kubernetes secrets and configuration files
- **`microservices`** - Multi-service architecture with per-service environment files

All templates include the `copy_to` field for environment files to automatically create `.env` from `.env.prod` after decryption.

If an invalid template name is provided, SecureFlow will use the default template.

This generates a `secureflow.yaml` with example entries that you can customize.

## Configuration Validation

SecureFlow validates your configuration automatically:

- **Missing files**: Warns but continues with other files
- **Invalid YAML**: Shows syntax error with line number
- **Missing directories**: Creates them automatically
- **Duplicate outputs**: Allowed (last one wins)

## Environment-Specific Configurations

### Development

```yaml
output_dir: dev_encrypted
test_output_dir: dev_test

files:
  - input: .env.development
    output: dev.env.encrypted
```

### Staging

```yaml
output_dir: staging_encrypted
test_output_dir: staging_test

files:
  - input: .env.staging
    output: staging.env.encrypted
```

### Production

```yaml
output_dir: prod_encrypted
test_output_dir: prod_test

files:
  - input: .env.production
    output: prod.env.encrypted
```

## Troubleshooting Configuration Issues

### "Config file not found"

```bash
Error: config file not found: secureflow.yaml
```

**Solution**: Run `secureflow init` or specify config path:
```bash
secureflow encrypt --config ./path/to/config.yaml
```

### "Input file not found"

```bash
⚠️  Warning: .env.prod not found, skipping
```

**Solution**: 
- Verify the file exists at the specified path
- Check the path is relative to project root
- Remove the entry from config if file is optional

### "Failed to parse YAML"

```bash
Error: yaml: line 5: mapping values are not allowed in this context
```

**Solution**:
- Check YAML syntax (proper indentation, colons, dashes)
- Validate YAML at [yamllint.com](https://www.yamllint.com/)
- Ensure proper spacing after colons

## See Also

- [CI/CD Usage Guide](./cicd-usage.md) - Using SecureFlow in CI/CD pipelines
- [Security Guide](./security.md) - Security best practices
- [Troubleshooting](./troubleshooting.md) - Common issues and solutions
