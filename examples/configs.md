# Sample Configuration Examples

This directory contains sample `secureflow.yaml` configuration files for different use cases.

## Basic Configuration

**File**: `basic-secureflow.yaml`

```yaml
output_dir: enc_keys
test_output_dir: test_dec_keys

files:
  - input: .env.prod
    output: .env.prod.encrypted
```

Minimal configuration for encrypting a single environment file.

## Flutter/React Native Mobile App

**File**: `mobile-app-secureflow.yaml`

```yaml
output_dir: encrypted_secrets
test_output_dir: test_decrypted

files:
  - input: .env.production
    output: .env.production.encrypted
  - input: android/app/keystore.jks
    output: keystore.jks.encrypted
  - input: android/key.properties
    output: key.properties.encrypted
  - input: android/service-account.json
    output: service-account.json.encrypted
  - input: ios/Runner/GoogleService-Info.plist
    output: GoogleService-Info.plist.encrypted
```

Configuration for mobile app projects with Android and iOS credentials.

## Web Application

**File**: `webapp-secureflow.yaml`

```yaml
output_dir: enc_keys
test_output_dir: test_dec_keys

files:
  - input: .env.production
    output: .env.production.encrypted
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

Configuration for web applications with multiple environments and SSL certificates.

## Microservices

**File**: `microservices-secureflow.yaml`

```yaml
output_dir: encrypted
test_output_dir: decrypted_test

files:
  # Auth Service
  - input: services/auth/.env.prod
    output: auth-env.prod.encrypted
  - input: services/auth/jwt-keys/private.pem
    output: auth-jwt-private.pem.encrypted
  
  # API Service
  - input: services/api/.env.prod
    output: api-env.prod.encrypted
  - input: services/api/config/database.json
    output: api-database.json.encrypted
  
  # Worker Service
  - input: services/worker/.env.prod
    output: worker-env.prod.encrypted
  
  # Shared Secrets
  - input: shared/redis.conf
    output: shared-redis.conf.encrypted
  - input: shared/rabbitmq.json
    output: shared-rabbitmq.json.encrypted
```

Configuration for microservices architecture with multiple services.

## Docker Deployment

**File**: `docker-secureflow.yaml`

```yaml
output_dir: docker/secrets/encrypted
test_output_dir: docker/secrets/test

files:
  - input: docker/.env.production
    output: docker-env.production.encrypted
  - input: docker/compose/.env.db
    output: docker-env.db.encrypted
  - input: docker/nginx/ssl/private.key
    output: nginx-ssl-private.key.encrypted
  - input: docker/nginx/ssl/certificate.crt
    output: nginx-ssl-certificate.crt.encrypted
```

Configuration for Docker deployments with environment files and SSL certificates.

## Kubernetes Secrets

**File**: `kubernetes-secureflow.yaml`

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

Configuration for Kubernetes secret files.

## Multi-Environment Setup

**File**: `multi-env-secureflow.yaml`

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

Configuration for projects with multiple environments.

## Usage

To use any of these configurations:

1. Copy the desired configuration to your project root as `secureflow.yaml`
2. Adjust the file paths to match your project structure
3. Run `secureflow encrypt` to encrypt your files

Example:

```bash
# Copy the mobile app configuration
cp examples/configs/mobile-app-secureflow.yaml ./secureflow.yaml

# Edit to match your project
vim secureflow.yaml

# Encrypt files
secureflow encrypt
```

Or use a custom config file directly:

```bash
secureflow encrypt --config examples/configs/webapp-secureflow.yaml
```
