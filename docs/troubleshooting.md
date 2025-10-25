# Troubleshooting Guide

This guide helps you diagnose and fix common issues when using SecureFlow.

## Table of Contents

- [Installation Issues](#installation-issues)
- [Encryption Issues](#encryption-issues)
- [Decryption Issues](#decryption-issues)
- [Configuration Issues](#configuration-issues)
- [CI/CD Issues](#cicd-issues)
- [Permission Issues](#permission-issues)
- [Performance Issues](#performance-issues)

## Installation Issues

### Binary Not Found After Installation

**Problem**:
```bash
$ secureflow --version
secureflow: command not found
```

**Solutions**:

1. **Check if binary is in PATH**:
   ```bash
   which secureflow
   ```

2. **Add to PATH if installed in custom location**:
   ```bash
   export PATH=$PATH:/path/to/secureflow
   # Add to ~/.bashrc or ~/.zshrc for persistence
   echo 'export PATH=$PATH:/path/to/secureflow' >> ~/.bashrc
   ```

3. **Use full path**:
   ```bash
   /usr/local/bin/secureflow --version
   ```

4. **Reinstall to standard location**:
   ```bash
   sudo mv secureflow /usr/local/bin/
   ```

### Permission Denied During Installation

**Problem**:
```bash
$ mv secureflow /usr/local/bin/
mv: cannot move 'secureflow' to '/usr/local/bin/secureflow': Permission denied
```

**Solutions**:

1. **Use sudo**:
   ```bash
   sudo mv secureflow /usr/local/bin/
   ```

2. **Install to user directory** (no sudo needed):
   ```bash
   mkdir -p ~/.local/bin
   mv secureflow ~/.local/bin/
   export PATH=$PATH:~/.local/bin
   ```

3. **Change ownership** (if applicable):
   ```bash
   sudo chown $USER:$USER /usr/local/bin/secureflow
   ```

### Download Fails

**Problem**:
```bash
$ wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
--2024-01-15 10:30:45--  https://github.com/...
Resolving github.com... failed: Temporary failure in name resolution.
```

**Solutions**:

1. **Check internet connection**:
   ```bash
   ping github.com
   ```

2. **Use curl instead of wget**:
   ```bash
   curl -L -o secureflow https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
   ```

3. **Use proxy if required**:
   ```bash
   export https_proxy=http://proxy.example.com:8080
   wget https://github.com/...
   ```

4. **Download manually**:
   - Visit https://github.com/MayR-Labs/secureflow-go/releases
   - Download appropriate binary
   - Transfer to target machine

### Wrong Architecture

**Problem**:
```bash
$ ./secureflow --version
bash: ./secureflow: cannot execute binary file: Exec format error
```

**Solutions**:

1. **Check your system architecture**:
   ```bash
   uname -m
   ```
   - `x86_64` or `amd64` → Use `secureflow-linux-amd64` or `secureflow-darwin-amd64`
   - `arm64` or `aarch64` → Use `secureflow-darwin-arm64` or `secureflow-linux-arm64`

2. **Download correct binary**:
   ```bash
   # For Linux ARM64
   wget https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-arm64
   ```

3. **Build from source** if no pre-built binary available:
   ```bash
   git clone https://github.com/MayR-Labs/secureflow-go.git
   cd secureflow-go
   go build -o secureflow
   ```

## Encryption Issues

### File Not Found During Encryption

**Problem**:
```bash
$ secureflow encrypt
⚠️  Warning: .env.prod not found, skipping
```

**Solutions**:

1. **Verify file exists**:
   ```bash
   ls -la .env.prod
   ```

2. **Check path in config**:
   ```yaml
   files:
     - input: .env.prod          # Relative to project root
       output: .env.prod.encrypted
   ```

3. **Create the file if missing**:
   ```bash
   touch .env.prod
   # Add your environment variables
   ```

4. **Remove from config if not needed**:
   ```yaml
   # Comment out or remove
   # - input: .env.prod
   #   output: .env.prod.encrypted
   ```

### Cannot Create Output Directory

**Problem**:
```bash
$ secureflow encrypt
Error: failed to create output directory: permission denied
```

**Solutions**:

1. **Check directory permissions**:
   ```bash
   ls -ld .
   ```

2. **Create directory manually**:
   ```bash
   mkdir -p enc_keys
   chmod 755 enc_keys
   ```

3. **Use different output directory**:
   ```yaml
   output_dir: ./encrypted  # Try different location
   ```

4. **Run with appropriate permissions**:
   ```bash
   sudo secureflow encrypt  # Last resort
   ```

### Encryption Succeeds But No Files Created

**Problem**:
Encryption completes but no files appear in output directory.

**Solutions**:

1. **Check output directory exists**:
   ```bash
   ls -la enc_keys/
   ```

2. **Verify config file is being read**:
   ```bash
   cat secureflow.yaml
   ```

3. **Use explicit config path**:
   ```bash
   secureflow encrypt --config ./secureflow.yaml
   ```

4. **Check for errors in output**:
   ```bash
   secureflow encrypt 2>&1 | tee encryption.log
   ```

### Report File Not Generated

**Problem**:
Encryption succeeds but no `report.txt` is created.

**Solutions**:

1. **Check output directory**:
   ```bash
   ls -la enc_keys/report.txt
   ```

2. **Verify write permissions**:
   ```bash
   chmod 755 enc_keys
   ```

3. **Check for disk space**:
   ```bash
   df -h .
   ```

## Decryption Issues

### Wrong Password Error

**Problem**:
```bash
$ secureflow decrypt
Error: decryption failed (wrong password?)
```

**Solutions**:

1. **Check password hint**:
   ```bash
   cat enc_keys/report.txt
   ```

2. **Verify you're using correct password**:
   - Check password manager
   - Verify no extra spaces
   - Try password in quotes: `--password "your password"`

3. **Test with known password**:
   ```bash
   # Re-encrypt a test file
   echo "test" > test.txt
   secureflow encrypt --password "testpass" --non-interactive
   secureflow decrypt --password "testpass" --non-interactive
   ```

4. **Check encrypted file integrity**:
   ```bash
   # Verify file starts with "Salted__"
   head -c 8 enc_keys/.env.prod.encrypted | od -c
   # Should show: S a l t e d _ _
   ```

### Decrypted Files Not Created

**Problem**:
Decryption succeeds but files don't appear.

**Solutions**:

1. **Check file paths in config**:
   ```yaml
   files:
     - input: .env.prod          # This is where decrypted file goes
       output: .env.prod.encrypted
   ```

2. **Verify parent directories exist**:
   ```bash
   # If input is android/app/keystore.jks
   mkdir -p android/app
   ```

3. **Check permissions**:
   ```bash
   ls -la .env.prod
   chmod 644 .env.prod
   ```

4. **Use test command to verify**:
   ```bash
   secureflow test --password "yourpass" --non-interactive
   ls -la test_dec_keys/
   ```

### Corrupted Encrypted File

**Problem**:
```bash
Error: failed to decrypt file: cipher: message authentication failed
```

**Solutions**:

1. **Check file size**:
   ```bash
   ls -lh enc_keys/*.encrypted
   # Should be non-zero
   ```

2. **Verify file integrity**:
   ```bash
   # Check if file starts with "Salted__"
   head -c 8 enc_keys/.env.prod.encrypted
   ```

3. **Re-encrypt from source**:
   ```bash
   # If you still have the original file
   secureflow encrypt --password "yourpass"
   ```

4. **Restore from backup**:
   ```bash
   git checkout enc_keys/.env.prod.encrypted
   ```

## Configuration Issues

### Config File Not Found

**Problem**:
```bash
$ secureflow encrypt
Error: config file not found: secureflow.yaml
```

**Solutions**:

1. **Create config file**:
   ```bash
   secureflow init
   ```

2. **Specify config path**:
   ```bash
   secureflow encrypt --config ./path/to/secureflow.yaml
   ```

3. **Verify current directory**:
   ```bash
   pwd
   ls -la secureflow.yaml
   ```

### YAML Syntax Error

**Problem**:
```bash
$ secureflow encrypt
Error: yaml: line 5: mapping values are not allowed in this context
```

**Solutions**:

1. **Check YAML syntax**:
   ```yaml
   # Correct
   output_dir: enc_keys
   
   # Wrong
   output_dir:enc_keys  # Missing space after colon
   ```

2. **Validate YAML online**:
   - Visit http://www.yamllint.com/
   - Paste your config
   - Fix reported errors

3. **Check indentation**:
   ```yaml
   # Correct (use spaces, not tabs)
   files:
     - input: .env.prod
       output: .env.prod.encrypted
   
   # Wrong (mixed indentation)
   files:
   	- input: .env.prod  # Tab used here
       output: .env.prod.encrypted  # Spaces used here
   ```

4. **Use YAML linter**:
   ```bash
   # Install yamllint
   pip install yamllint
   
   # Validate config
   yamllint secureflow.yaml
   ```

### Empty or Invalid Config

**Problem**:
```bash
Error: no files specified in configuration
```

**Solutions**:

1. **Add file entries**:
   ```yaml
   output_dir: enc_keys
   test_output_dir: test_dec_keys
   
   files:
     - input: .env.prod
       output: .env.prod.encrypted
   ```

2. **Verify config structure**:
   ```bash
   cat secureflow.yaml
   ```

3. **Reinitialize config**:
   ```bash
   mv secureflow.yaml secureflow.yaml.bak
   secureflow init
   # Copy file entries from backup
   ```

## CI/CD Issues

### Secret Not Available in Pipeline

**Problem**:
```bash
Error: decryption failed (wrong password?)
# In CI/CD logs
```

**Solutions**:

1. **Verify secret is set**:
   - GitHub: Settings → Secrets → Actions
   - GitLab: Settings → CI/CD → Variables
   - Check secret name matches usage: `${{ secrets.SECUREFLOW_PASSWORD }}`

2. **Check secret scope**:
   - Ensure secret is available to branch
   - Check if secret is environment-specific

3. **Test secret value**:
   ```yaml
   # DO NOT leave this in production!
   - run: echo "Password length: ${#SECUREFLOW_PASSWORD}"
   ```

4. **Check for extra characters**:
   - No trailing newlines
   - No extra spaces
   - Copy secret value again

### Binary Download Fails in CI/CD

**Problem**:
```bash
wget: unable to resolve host address 'github.com'
```

**Solutions**:

1. **Check network connectivity**:
   ```yaml
   - run: ping -c 3 github.com
   ```

2. **Use alternative download method**:
   ```yaml
   - run: |
       curl -L -o secureflow https://github.com/MayR-Labs/secureflow-go/releases/latest/download/secureflow-linux-amd64
   ```

3. **Use proxy if required**:
   ```yaml
   - run: |
       export https_proxy=${{ secrets.PROXY_URL }}
       wget https://github.com/...
   ```

4. **Cache the binary**:
   ```yaml
   # GitHub Actions example
   - uses: actions/cache@v3
     with:
       path: /usr/local/bin/secureflow
       key: secureflow-${{ runner.os }}-latest
   ```

### Permission Issues in CI/CD

**Problem**:
```bash
mv: cannot move 'secureflow' to '/usr/local/bin/secureflow': Permission denied
```

**Solutions**:

1. **Remove sudo** (if in Docker):
   ```yaml
   - run: mv secureflow-linux-amd64 /usr/local/bin/secureflow
   ```

2. **Use user directory**:
   ```yaml
   - run: |
       mkdir -p ~/.local/bin
       mv secureflow-linux-amd64 ~/.local/bin/secureflow
       echo "$HOME/.local/bin" >> $GITHUB_PATH
   ```

3. **Use workspace directory**:
   ```yaml
   - run: |
       chmod +x secureflow-linux-amd64
       ./secureflow-linux-amd64 decrypt --password "$PASSWORD" --non-interactive
   ```

### Artifacts Not Persisting

**Problem**:
Decrypted files not available in subsequent jobs.

**Solutions**:

1. **Check artifact configuration**:
   ```yaml
   # GitHub Actions
   - uses: actions/upload-artifact@v3
     with:
       name: secrets
       path: |
         .env.prod
         android/app/keystore.jks
   ```

2. **Verify artifact download**:
   ```yaml
   - uses: actions/download-artifact@v3
     with:
       name: secrets
   ```

3. **Check file paths match**:
   ```yaml
   # Upload path must match actual file location
   path: .env.production  # Not .env.prod if file is named differently
   ```

4. **Use workspace persistence** (GitLab):
   ```yaml
   artifacts:
     paths:
       - .env.prod
     expire_in: 1 hour
   ```

## Permission Issues

### Cannot Read Config File

**Problem**:
```bash
Error: failed to read config: permission denied
```

**Solutions**:

1. **Check file permissions**:
   ```bash
   ls -la secureflow.yaml
   chmod 644 secureflow.yaml
   ```

2. **Check file ownership**:
   ```bash
   ls -la secureflow.yaml
   chown $USER:$USER secureflow.yaml
   ```

3. **Run with appropriate user**:
   ```bash
   sudo -u correctuser secureflow encrypt
   ```

### Cannot Write Encrypted Files

**Problem**:
```bash
Error: failed to write encrypted file: permission denied
```

**Solutions**:

1. **Check output directory permissions**:
   ```bash
   ls -ld enc_keys
   chmod 755 enc_keys
   ```

2. **Create directory with correct permissions**:
   ```bash
   mkdir -p enc_keys
   chmod 755 enc_keys
   ```

3. **Check disk space**:
   ```bash
   df -h .
   ```

### Cannot Read Input Files

**Problem**:
```bash
Error: failed to read input file: permission denied
```

**Solutions**:

1. **Check input file permissions**:
   ```bash
   ls -la .env.prod
   chmod 644 .env.prod
   ```

2. **Verify file ownership**:
   ```bash
   chown $USER:$USER .env.prod
   ```

## Performance Issues

### Encryption Takes Too Long

**Problem**:
Encryption process is very slow.

**Solutions**:

1. **Check file sizes**:
   ```bash
   ls -lh .env.prod android/app/keystore.jks
   ```
   Large files (>100MB) take longer to encrypt

2. **Monitor system resources**:
   ```bash
   top
   # Or
   htop
   ```

3. **Encrypt files individually**:
   ```bash
   # Split large configs
   secureflow encrypt --config secureflow-part1.yaml
   secureflow encrypt --config secureflow-part2.yaml
   ```

4. **Check for I/O bottlenecks**:
   ```bash
   iostat -x 1
   ```

### CI/CD Pipeline Timeout

**Problem**:
Pipeline times out during decryption.

**Solutions**:

1. **Increase timeout**:
   ```yaml
   # GitHub Actions
   - name: Decrypt
     timeout-minutes: 10
     run: secureflow decrypt --password "$PASSWORD" --non-interactive
   ```

2. **Decrypt only necessary files**:
   ```yaml
   # Use different configs for different jobs
   secureflow decrypt --config secureflow-minimal.yaml
   ```

3. **Use faster runners**:
   - GitHub: Use ubuntu-latest instead of custom runners
   - GitLab: Use runners with better specs

4. **Cache decrypted files** (if safe):
   ```yaml
   # Only if appropriate for your use case
   - uses: actions/cache@v3
     with:
       path: .env.prod
       key: secrets-${{ hashFiles('enc_keys/*.encrypted') }}
   ```

## Getting Help

If you've tried the solutions above and still have issues:

### Collect Information

1. **SecureFlow version**:
   ```bash
   secureflow --version
   ```

2. **System information**:
   ```bash
   uname -a
   ```

3. **Configuration** (sanitized):
   ```bash
   cat secureflow.yaml
   ```

4. **Error messages**:
   ```bash
   secureflow encrypt 2>&1 | tee error.log
   ```

### Report Issues

1. **GitHub Issues**: https://github.com/MayR-Labs/secureflow-go/issues
2. **Include**:
   - SecureFlow version
   - Operating system
   - Full error message
   - Steps to reproduce
   - Configuration (remove sensitive info)

### Community Support

- Check existing issues: https://github.com/MayR-Labs/secureflow-go/issues
- Read documentation: https://github.com/MayR-Labs/secureflow-go
- Contact: MayR Labs support

## Common Error Messages Reference

| Error Message | Likely Cause | Solution |
|--------------|--------------|----------|
| `command not found` | SecureFlow not in PATH | Add to PATH or use full path |
| `permission denied` | Insufficient permissions | Use sudo or install to user directory |
| `config file not found` | Wrong directory or missing file | Run `secureflow init` or specify path |
| `wrong password` | Incorrect password | Check password hint in report.txt |
| `file not found` | Input file missing | Verify file exists or remove from config |
| `yaml: line X` | YAML syntax error | Fix YAML syntax at specified line |
| `cipher: message authentication failed` | Corrupted encrypted file | Re-encrypt from source or restore backup |
| `failed to create directory` | Permission or disk space | Check permissions and disk space |

## See Also

- [Configuration Guide](./configuration.md) - Configuration best practices
- [CI/CD Usage Guide](./cicd-usage.md) - CI/CD integration
- [Security Guide](./security.md) - Security best practices
