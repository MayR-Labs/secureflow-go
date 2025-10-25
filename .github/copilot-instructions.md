# GitHub Copilot Instructions for SecureFlow

## Project Overview

SecureFlow is a Go-based CLI tool for securely encrypting and decrypting sensitive files using AES-256-CBC encryption with OpenSSL-compatible standards. It's designed to replace fragile Bash scripts with a fast, cross-platform executable.

## Architecture

### Project Structure
```
secureflow-go/
├── cmd/              # CLI commands (Cobra)
│   ├── root.go       # CLI entrypoint and root command
│   ├── encrypt.go    # Encryption command
│   ├── decrypt.go    # Decryption command
│   ├── test.go       # Test decryption command
│   └── init.go       # Initialize config command
├── internal/
│   ├── crypto/       # Encryption/decryption logic
│   ├── config/       # secureflow.yaml handling
│   └── utils/        # File handling, error logging
├── main.go           # Application entry point
├── go.mod            # Go module dependencies
└── README.md
```

## Key Technologies

- **Language**: Go 1.17+
- **CLI Framework**: cobra (github.com/spf13/cobra)
- **Config Format**: YAML (gopkg.in/yaml.v3)
- **Encryption**: AES-256-CBC via crypto/aes, crypto/cipher, OpenSSL-compatible

## Coding Standards

### Go Best Practices
- Follow standard Go formatting (use `gofmt` or `goimports`)
- Use meaningful variable and function names
- Keep functions focused and single-purpose
- Handle errors explicitly; never ignore them
- Use context for cancellation where appropriate
- Write idiomatic Go code following [Effective Go](https://golang.org/doc/effective_go.html)

### Error Handling
- Return errors, don't panic (except in main.go for fatal errors)
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Log errors with consistent formatting
- Exit codes: 0 for success, 1 for errors

### Code Organization
- Keep `internal/` packages focused on core logic
- Use `cmd/` for CLI command definitions only
- Avoid circular dependencies
- Export only what's necessary (use lowercase for internal functions)

### Comments and Documentation
- Add package documentation comments
- Document exported functions and types
- Use clear, concise comments
- Avoid obvious comments

## Security Requirements

1. **Password Handling**
   - Never log or store passwords in plaintext
   - Clear password variables after use when possible
   - Use secure password input (no echo) for interactive mode

2. **Encryption Standards**
   - AES-256-CBC with PBKDF2 key derivation
   - OpenSSL-compatible format (Salted__ header)
   - Use crypto/rand for random generation

3. **File Operations**
   - Validate file paths before operations
   - Handle permission errors gracefully
   - Create directories with appropriate permissions (0755)
   - Set encrypted file permissions appropriately (0644)

## CLI Command Guidelines

### Command Structure
- Use Cobra for all commands
- Support both interactive and non-interactive modes
- Provide helpful error messages
- Include usage examples in help text

### Flags and Options
- `--config`: Custom config file path
- `--password`: Non-interactive password (for CI/CD)
- `--non-interactive`: Skip all prompts
- `--version`: Show version information

### User Experience
- Use colored output for better readability (optional, graceful degradation)
- Show progress for long operations
- Provide clear success/failure messages
- Include emoji in output for visual clarity (✅, ❌, 🔐, etc.)

## Configuration File (secureflow.yaml)

```yaml
output_dir: enc_keys           # Encrypted files directory
test_output_dir: test_dec_keys # Test decryption directory

files:
  - input: .env.prod                    # Source file
    output: .env.prod.encrypted         # Encrypted filename
  - input: android/app/keystore.jks
    output: keystore.jks.encrypted
```

## Testing

- Write unit tests for crypto functions
- Test error conditions
- Mock file I/O where appropriate
- Test both interactive and non-interactive modes
- Run tests with: `go test ./...`

## Build and Release

- Build command: `go build -o secureflow`
- Cross-compile for Linux, macOS, Windows
- Version management via git tags
- Binary should be statically linked when possible

## Dependencies

- Prefer standard library when possible
- Keep dependencies minimal
- Use go modules for dependency management
- Pin dependency versions in go.mod

## Common Commands

```bash
# Run locally
go run main.go [command]

# Build
go build -o secureflow

# Test
go test ./...

# Format code
gofmt -w .

# Run linter (if golangci-lint installed)
golangci-lint run
```

## Encryption/Decryption Implementation

The crypto package should implement OpenSSL-compatible encryption:
- Use PBKDF2 for key derivation from password
- Add "Salted__" prefix followed by 8-byte salt
- Use AES-256-CBC mode
- Match OpenSSL's `openssl enc -aes-256-cbc -salt -pbkdf2` behavior

## Report Generation

After encryption, generate a report.txt file with:
- Encryption note and password hint
- Timestamp
- For each file: name, encrypted name, size, line count, last modified

## Contribution Guidelines

When extending or modifying this project:
1. Maintain backward compatibility with encrypted files
2. Update README.md with new features
3. Add tests for new functionality
4. Follow existing code patterns and structure
5. Keep CLI interface consistent
