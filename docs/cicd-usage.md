# CI/CD Usage Guide

This guide covers how to integrate SecureFlow into your CI/CD pipelines across different platforms.

## Table of Contents

- [Overview](#overview)
- [Platform-Specific Guides](#platform-specific-guides)
  - [GitHub Actions](#github-actions)
  - [GitLab CI](#gitlab-ci)
  - [Bitbucket Pipelines](#bitbucket-pipelines)
  - [Jenkins](#jenkins)
  - [CircleCI](#circleci)
  - [Azure Pipelines](#azure-pipelines)
- [Best Practices](#best-practices)
- [Security Considerations](#security-considerations)
- [Troubleshooting](#troubleshooting)

## Overview

SecureFlow is designed to work seamlessly in CI/CD environments with its non-interactive mode. The general workflow is:

1. **Store encryption password** as a secret in your CI/CD platform
2. **Install SecureFlow** in your pipeline
3. **Decrypt secrets** before building/deploying
4. **Use decrypted files** in subsequent steps
5. **Optionally clean up** decrypted files after use

## Platform-Specific Guides

### GitHub Actions

#### Basic Setup

**.github/workflows/deploy.yml**:

```yaml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v3
      
      - name: Install SecureFlow
        run: |
          wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
          chmod +x secureflow-linux-amd64
          sudo mv secureflow-linux-amd64 /usr/local/bin/secureflow
      
      - name: Decrypt secrets
        run: |
          secureflow decrypt --password "${{ secrets.SECUREFLOW_PASSWORD }}" --non-interactive
        
      - name: Build application
        run: |
          # Your build commands here
          npm run build
      
      - name: Deploy
        run: |
          # Your deployment commands here
          echo "Deploying with decrypted secrets..."
```

#### Advanced: Multiple Environments

```yaml
name: Deploy Multi-Environment

on:
  push:
    branches: [main, staging, development]

jobs:
  deploy:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Install SecureFlow
        run: |
          wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
          chmod +x secureflow-linux-amd64
          sudo mv secureflow-linux-amd64 /usr/local/bin/secureflow
      
      - name: Decrypt Production Secrets
        if: github.ref == 'refs/heads/main'
        run: |
          secureflow decrypt --config secureflow.prod.yaml --password "${{ secrets.PROD_PASSWORD }}" --non-interactive
      
      - name: Decrypt Staging Secrets
        if: github.ref == 'refs/heads/staging'
        run: |
          secureflow decrypt --config secureflow.staging.yaml --password "${{ secrets.STAGING_PASSWORD }}" --non-interactive
      
      - name: Deploy
        run: |
          ./deploy.sh
```

#### Using Reusable Workflow

**.github/workflows/decrypt-secrets.yml**:

```yaml
name: Decrypt Secrets

on:
  workflow_call:
    inputs:
      config-file:
        required: false
        type: string
        default: 'secureflow.yaml'
    secrets:
      password:
        required: true

jobs:
  decrypt:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Install SecureFlow
        run: |
          wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
          chmod +x secureflow-linux-amd64
          sudo mv secureflow-linux-amd64 /usr/local/bin/secureflow
      
      - name: Decrypt
        run: |
          secureflow decrypt --config ${{ inputs.config-file }} --password "${{ secrets.password }}" --non-interactive
```

Then use it in other workflows:

```yaml
jobs:
  decrypt:
    uses: ./.github/workflows/decrypt-secrets.yml
    secrets:
      password: ${{ secrets.SECUREFLOW_PASSWORD }}
  
  build:
    needs: decrypt
    runs-on: ubuntu-latest
    steps:
      - name: Build
        run: npm run build
```

#### Storing Secrets

1. Go to your repository on GitHub
2. Navigate to **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Name: `SECUREFLOW_PASSWORD`
5. Value: Your encryption password
6. Click **Add secret**

### GitLab CI

#### Basic Setup

**.gitlab-ci.yml**:

```yaml
stages:
  - decrypt
  - build
  - deploy

variables:
  SECUREFLOW_VERSION: "latest"

before_script:
  - apt-get update -qq

decrypt_secrets:
  stage: decrypt
  image: ubuntu:22.04
  before_script:
    - apt-get update && apt-get install -y wget
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

build:
  stage: build
  image: node:18
  dependencies:
    - decrypt_secrets
  script:
    - npm install
    - npm run build
  artifacts:
    paths:
      - dist/

deploy:
  stage: deploy
  dependencies:
    - decrypt_secrets
    - build
  script:
    - echo "Deploying with decrypted secrets..."
    - ./deploy.sh
  only:
    - main
```

#### Environment-Specific Configuration

```yaml
.decrypt_template:
  before_script:
    - apt-get update && apt-get install -y wget
    - wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
    - chmod +x secureflow-linux-amd64
    - mv secureflow-linux-amd64 /usr/local/bin/secureflow

decrypt_production:
  extends: .decrypt_template
  stage: decrypt
  script:
    - secureflow decrypt --config secureflow.prod.yaml --password "$PROD_PASSWORD" --non-interactive
  only:
    - main
  artifacts:
    paths:
      - .env.production
    expire_in: 1 hour

decrypt_staging:
  extends: .decrypt_template
  stage: decrypt
  script:
    - secureflow decrypt --config secureflow.staging.yaml --password "$STAGING_PASSWORD" --non-interactive
  only:
    - staging
  artifacts:
    paths:
      - .env.staging
    expire_in: 1 hour
```

#### Storing Secrets

1. Go to your project on GitLab
2. Navigate to **Settings** → **CI/CD** → **Variables**
3. Click **Add variable**
4. Key: `SECUREFLOW_PASSWORD`
5. Value: Your encryption password
6. Check **Mask variable** and **Protect variable**
7. Click **Add variable**

### Bitbucket Pipelines

#### Basic Setup

**bitbucket-pipelines.yml**:

```yaml
pipelines:
  default:
    - step:
        name: Decrypt and Deploy
        image: ubuntu:22.04
        script:
          # Install SecureFlow
          - apt-get update && apt-get install -y wget
          - wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
          - chmod +x secureflow-linux-amd64
          - mv secureflow-linux-amd64 /usr/local/bin/secureflow
          
          # Decrypt secrets
          - secureflow decrypt --password "$SECUREFLOW_PASSWORD" --non-interactive
          
          # Build and deploy
          - echo "Building application..."
          - ./build.sh
          - echo "Deploying..."
          - ./deploy.sh
```

#### Branch-Specific Deployment

```yaml
pipelines:
  branches:
    main:
      - step:
          name: Decrypt Production Secrets
          image: ubuntu:22.04
          script:
            - apt-get update && apt-get install -y wget
            - wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
            - chmod +x secureflow-linux-amd64
            - mv secureflow-linux-amd64 /usr/local/bin/secureflow
            - secureflow decrypt --config secureflow.prod.yaml --password "$PROD_PASSWORD" --non-interactive
          artifacts:
            - .env.production
            - android/app/keystore.jks
      
      - step:
          name: Deploy Production
          script:
            - ./deploy-production.sh
    
    staging:
      - step:
          name: Decrypt Staging Secrets
          image: ubuntu:22.04
          script:
            - apt-get update && apt-get install -y wget
            - wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
            - chmod +x secureflow-linux-amd64
            - mv secureflow-linux-amd64 /usr/local/bin/secureflow
            - secureflow decrypt --config secureflow.staging.yaml --password "$STAGING_PASSWORD" --non-interactive
          artifacts:
            - .env.staging
      
      - step:
          name: Deploy Staging
          script:
            - ./deploy-staging.sh
```

#### Storing Secrets

1. Go to your repository on Bitbucket
2. Navigate to **Repository settings** → **Pipelines** → **Repository variables**
3. Click **Add variable**
4. Name: `SECUREFLOW_PASSWORD`
5. Value: Your encryption password
6. Check **Secured** to mask the value in logs
7. Click **Add**

### Jenkins

#### Declarative Pipeline

**Jenkinsfile**:

```groovy
pipeline {
    agent any
    
    environment {
        SECUREFLOW_PASSWORD = credentials('secureflow-password')
    }
    
    stages {
        stage('Install SecureFlow') {
            steps {
                sh '''
                    wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
                    chmod +x secureflow-linux-amd64
                    sudo mv secureflow-linux-amd64 /usr/local/bin/secureflow || mv secureflow-linux-amd64 /usr/local/bin/secureflow
                '''
            }
        }
        
        stage('Decrypt Secrets') {
            steps {
                sh 'secureflow decrypt --password "$SECUREFLOW_PASSWORD" --non-interactive'
            }
        }
        
        stage('Build') {
            steps {
                sh 'npm install'
                sh 'npm run build'
            }
        }
        
        stage('Deploy') {
            when {
                branch 'main'
            }
            steps {
                sh './deploy.sh'
            }
        }
    }
    
    post {
        always {
            // Clean up decrypted files
            sh 'rm -f .env.prod android/app/keystore.jks || true'
        }
    }
}
```

#### Scripted Pipeline with Multiple Environments

```groovy
node {
    def password
    
    stage('Checkout') {
        checkout scm
    }
    
    stage('Install SecureFlow') {
        sh '''
            wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
            chmod +x secureflow-linux-amd64
            mv secureflow-linux-amd64 /usr/local/bin/secureflow
        '''
    }
    
    stage('Decrypt Secrets') {
        if (env.BRANCH_NAME == 'main') {
            password = credentials('prod-secureflow-password')
            sh "secureflow decrypt --config secureflow.prod.yaml --password '${password}' --non-interactive"
        } else if (env.BRANCH_NAME == 'staging') {
            password = credentials('staging-secureflow-password')
            sh "secureflow decrypt --config secureflow.staging.yaml --password '${password}' --non-interactive"
        } else {
            password = credentials('dev-secureflow-password')
            sh "secureflow decrypt --config secureflow.dev.yaml --password '${password}' --non-interactive"
        }
    }
    
    stage('Build') {
        sh 'npm install && npm run build'
    }
    
    stage('Deploy') {
        sh './deploy.sh'
    }
}
```

#### Storing Secrets in Jenkins

1. Go to **Manage Jenkins** → **Manage Credentials**
2. Select the appropriate domain (usually "Global")
3. Click **Add Credentials**
4. Kind: **Secret text**
5. Secret: Your encryption password
6. ID: `secureflow-password`
7. Description: SecureFlow encryption password
8. Click **OK**

### CircleCI

#### Basic Setup

**.circleci/config.yml**:

```yaml
version: 2.1

jobs:
  decrypt-and-deploy:
    docker:
      - image: ubuntu:22.04
    
    steps:
      - checkout
      
      - run:
          name: Install Dependencies
          command: |
            apt-get update
            apt-get install -y wget
      
      - run:
          name: Install SecureFlow
          command: |
            wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
            chmod +x secureflow-linux-amd64
            mv secureflow-linux-amd64 /usr/local/bin/secureflow
      
      - run:
          name: Decrypt Secrets
          command: |
            secureflow decrypt --password "$SECUREFLOW_PASSWORD" --non-interactive
      
      - run:
          name: Build
          command: |
            # Your build commands
            ./build.sh
      
      - run:
          name: Deploy
          command: |
            # Your deployment commands
            ./deploy.sh

workflows:
  version: 2
  build-and-deploy:
    jobs:
      - decrypt-and-deploy:
          filters:
            branches:
              only: main
```

#### Using Contexts for Multiple Environments

```yaml
version: 2.1

jobs:
  decrypt:
    docker:
      - image: ubuntu:22.04
    
    parameters:
      config-file:
        type: string
        default: "secureflow.yaml"
    
    steps:
      - checkout
      
      - run:
          name: Install Dependencies
          command: apt-get update && apt-get install -y wget
      
      - run:
          name: Install SecureFlow
          command: |
            wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
            chmod +x secureflow-linux-amd64
            mv secureflow-linux-amd64 /usr/local/bin/secureflow
      
      - run:
          name: Decrypt Secrets
          command: |
            secureflow decrypt --config << parameters.config-file >> --password "$SECUREFLOW_PASSWORD" --non-interactive
      
      - persist_to_workspace:
          root: .
          paths:
            - .env.*
            - android/app/keystore.jks

  deploy:
    docker:
      - image: node:18
    
    steps:
      - checkout
      
      - attach_workspace:
          at: .
      
      - run:
          name: Deploy
          command: ./deploy.sh

workflows:
  production:
    jobs:
      - decrypt:
          name: decrypt-production
          config-file: "secureflow.prod.yaml"
          context: production
          filters:
            branches:
              only: main
      
      - deploy:
          name: deploy-production
          requires:
            - decrypt-production
          filters:
            branches:
              only: main
  
  staging:
    jobs:
      - decrypt:
          name: decrypt-staging
          config-file: "secureflow.staging.yaml"
          context: staging
          filters:
            branches:
              only: staging
      
      - deploy:
          name: deploy-staging
          requires:
            - decrypt-staging
          filters:
            branches:
              only: staging
```

#### Storing Secrets

1. Go to your project on CircleCI
2. Navigate to **Project Settings** → **Environment Variables**
3. Click **Add Environment Variable**
4. Name: `SECUREFLOW_PASSWORD`
5. Value: Your encryption password
6. Click **Add Variable**

For context-based secrets:
1. Go to **Organization Settings** → **Contexts**
2. Create contexts (e.g., "production", "staging")
3. Add `SECUREFLOW_PASSWORD` to each context

### Azure Pipelines

#### Basic Setup

**azure-pipelines.yml**:

```yaml
trigger:
  - main

pool:
  vmImage: 'ubuntu-latest'

variables:
  - group: secureflow-secrets

steps:
  - checkout: self
  
  - bash: |
      wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
      chmod +x secureflow-linux-amd64
      sudo mv secureflow-linux-amd64 /usr/local/bin/secureflow
    displayName: 'Install SecureFlow'
  
  - bash: |
      secureflow decrypt --password "$(SECUREFLOW_PASSWORD)" --non-interactive
    displayName: 'Decrypt Secrets'
  
  - bash: |
      npm install
      npm run build
    displayName: 'Build Application'
  
  - bash: |
      ./deploy.sh
    displayName: 'Deploy'
```

#### Multi-Stage with Environments

```yaml
trigger:
  - main
  - staging

stages:
  - stage: Build
    jobs:
      - job: BuildJob
        pool:
          vmImage: 'ubuntu-latest'
        steps:
          - bash: |
              wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
              chmod +x secureflow-linux-amd64
              sudo mv secureflow-linux-amd64 /usr/local/bin/secureflow
            displayName: 'Install SecureFlow'
          
          - bash: |
              if [ "$(Build.SourceBranchName)" = "main" ]; then
                secureflow decrypt --config secureflow.prod.yaml --password "$(PROD_PASSWORD)" --non-interactive
              else
                secureflow decrypt --config secureflow.staging.yaml --password "$(STAGING_PASSWORD)" --non-interactive
              fi
            displayName: 'Decrypt Secrets'
          
          - bash: npm install && npm run build
            displayName: 'Build'
          
          - publish: $(System.DefaultWorkingDirectory)/dist
            artifact: build-output

  - stage: Deploy
    dependsOn: Build
    jobs:
      - deployment: DeployJob
        environment: production
        strategy:
          runOnce:
            deploy:
              steps:
                - download: current
                  artifact: build-output
                
                - bash: ./deploy.sh
                  displayName: 'Deploy'
```

#### Storing Secrets

1. Go to your project on Azure DevOps
2. Navigate to **Pipelines** → **Library**
3. Click **+ Variable group**
4. Name: `secureflow-secrets`
5. Add variable:
   - Name: `SECUREFLOW_PASSWORD`
   - Value: Your encryption password
   - Check the lock icon to make it a secret
6. Click **Save**

## Best Practices

### 1. Use Non-Interactive Mode

Always use `--non-interactive` flag in CI/CD:

```bash
secureflow decrypt --password "$PASSWORD" --non-interactive
```

### 2. Store Password as Secret

Never hardcode passwords in your pipeline files. Use your platform's secret management:

- GitHub Actions: Repository secrets
- GitLab CI: CI/CD variables (masked and protected)
- Bitbucket: Repository variables (secured)
- Jenkins: Credentials store
- CircleCI: Environment variables or contexts
- Azure Pipelines: Variable groups

### 3. Use Artifacts for Decrypted Files

Pass decrypted files between pipeline steps using artifacts:

```yaml
# GitLab CI example
artifacts:
  paths:
    - .env.prod
  expire_in: 1 hour
```

### 4. Clean Up After Use

Remove decrypted files after deployment to minimize exposure:

```bash
# Add to your pipeline
rm -f .env.prod android/app/keystore.jks
```

### 5. Separate Configs per Environment

Use different config files for different environments:

```bash
secureflow decrypt --config secureflow.prod.yaml --password "$PROD_PASSWORD" --non-interactive
secureflow decrypt --config secureflow.staging.yaml --password "$STAGING_PASSWORD" --non-interactive
```

### 6. Version Pin SecureFlow

Instead of `latest`, consider pinning to a specific version:

```bash
wget https://github.com/MayR-Labs/secureflow-go/releases/download/v1.0.0/secureflow-linux-amd64
```

### 7. Cache SecureFlow Binary

Cache the binary to speed up pipeline runs:

```yaml
# GitHub Actions example
- name: Cache SecureFlow
  uses: actions/cache@v3
  with:
    path: /usr/local/bin/secureflow
    key: secureflow-${{ runner.os }}-latest
```

### 8. Verify Decryption Success

Check that decryption succeeded before proceeding:

```bash
secureflow decrypt --password "$PASSWORD" --non-interactive || exit 1
```

### 9. Use Test Command First

In development pipelines, use `test` command to verify without overwriting:

```bash
secureflow test --password "$PASSWORD" --non-interactive
```

### 10. Monitor Pipeline Logs

Ensure passwords are masked in logs. Test by running a pipeline and checking output.

## Security Considerations

### Password Management

1. **Rotate passwords regularly**: Update your encryption password periodically
2. **Use strong passwords**: At least 16 characters with mixed case, numbers, and symbols
3. **Limit access**: Restrict who can view/edit secrets in your CI/CD platform
4. **Audit access**: Regularly review who has access to secrets

### Pipeline Security

1. **Protect branches**: Require approval for deployments to production
2. **Use protected variables**: Mark variables as protected in your CI/CD platform
3. **Minimize secret exposure**: Only decrypt in the steps that need the secrets
4. **Clean up artifacts**: Set short expiration times for artifacts containing secrets

### Network Security

1. **Use HTTPS**: Always download SecureFlow over HTTPS
2. **Verify checksums**: Consider verifying downloaded binary checksums
3. **Private runners**: Use self-hosted runners for sensitive workloads

## Troubleshooting

### "Wrong password" Error in CI/CD

```bash
Error: decryption failed (wrong password?)
```

**Solutions**:
- Verify the secret value in your CI/CD platform
- Check for extra spaces or line breaks in the password
- Test locally with the same password
- Ensure you're using the correct config file

### Binary Not Found

```bash
secureflow: command not found
```

**Solutions**:
- Check the installation step succeeded
- Verify the binary was moved to a directory in PATH
- Try using the full path: `/usr/local/bin/secureflow`

### Permission Denied

```bash
Permission denied: /usr/local/bin/secureflow
```

**Solutions**:
- Remove `sudo` if running in Docker
- Use a user-writable location: `~/.local/bin/secureflow`
- Add to PATH: `export PATH=$PATH:~/.local/bin`

### Artifacts Not Persisting

**Solutions**:
- Check artifact paths match your file locations
- Verify artifact retention settings
- Ensure previous step completed successfully
- Check workspace attachment in dependent jobs

### Timeout Issues

If decryption takes too long:

**Solutions**:
- Check file sizes (large files take longer)
- Verify network connectivity
- Increase timeout in pipeline configuration
- Use faster runner/agent

## Example Complete Workflow

Here's a complete example combining best practices:

### GitHub Actions Complete Example

```yaml
name: Production Deployment

on:
  push:
    branches: [main]

env:
  SECUREFLOW_VERSION: v1.0.0

jobs:
  decrypt:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Cache SecureFlow
        id: cache-secureflow
        uses: actions/cache@v3
        with:
          path: /usr/local/bin/secureflow
          key: secureflow-${{ runner.os }}-${{ env.SECUREFLOW_VERSION }}
      
      - name: Install SecureFlow
        if: steps.cache-secureflow.outputs.cache-hit != 'true'
        run: |
          wget https://github.com/MayR-Labs/secureflow-go/releases/download/${{ env.SECUREFLOW_VERSION }}/secureflow-linux-amd64
          chmod +x secureflow-linux-amd64
          sudo mv secureflow-linux-amd64 /usr/local/bin/secureflow
      
      - name: Decrypt Secrets
        run: |
          secureflow decrypt --config secureflow.prod.yaml --password "${{ secrets.PROD_PASSWORD }}" --non-interactive
        
      - name: Upload Decrypted Files
        uses: actions/upload-artifact@v3
        with:
          name: decrypted-secrets
          path: |
            .env.production
            android/app/keystore.jks
          retention-days: 1

  build:
    needs: decrypt
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Download Decrypted Files
        uses: actions/download-artifact@v3
        with:
          name: decrypted-secrets
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
      
      - name: Install Dependencies
        run: npm ci
      
      - name: Build
        run: npm run build
      
      - name: Upload Build Artifact
        uses: actions/upload-artifact@v3
        with:
          name: build-output
          path: dist/

  deploy:
    needs: build
    runs-on: ubuntu-latest
    environment: production
    steps:
      - uses: actions/checkout@v3
      
      - name: Download Build
        uses: actions/download-artifact@v3
        with:
          name: build-output
          path: dist/
      
      - name: Download Secrets
        uses: actions/download-artifact@v3
        with:
          name: decrypted-secrets
      
      - name: Deploy
        run: |
          ./deploy.sh
      
      - name: Clean Up
        if: always()
        run: |
          rm -f .env.production android/app/keystore.jks
```

## See Also

- [Configuration Guide](./configuration.md) - Detailed configuration options
- [Security Guide](./security.md) - Security best practices
- [Troubleshooting](./troubleshooting.md) - Common issues and solutions
