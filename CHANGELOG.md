# Changelog

## v1.3.0

### Added
- `secureflow refresh` (alias `rotate`): re-encrypt every file with a new password (or the same one), all-or-nothing, with plaintext kept in memory only.
- `secureflow add <file>...`: adds files to `secureflow.yaml` (keeping comments), appends them to `.gitignore`, and encrypts them with the existing password.
- `output` is now optional in `secureflow.yaml`; it's inferred from `input` (`android/app/keystore.jks` → `android-app-keystore.jks.encrypted`). Configs with duplicate encrypted names are rejected.
- The password hint is shown before the password prompt in `decrypt`, `test` and `refresh`, and again after a wrong password.
- `--hint`, `--note` and `--force` flags for `encrypt`; `SECUREFLOW_PASSWORD` / `SECUREFLOW_NEW_PASSWORD` environment variables.
- Composite GitHub Action (`uses: MayR-Labs/secureflow-go@v1.3.0`) with checksum and build-provenance verification.
- Release binaries carry GitHub build-provenance attestations (`gh attestation verify`).

### Changed
- **Encryption format**: new files use AES-256-GCM with PBKDF2-SHA256 at 600,000 iterations (was AES-256-CBC, 10,000 iterations). Legacy files still decrypt; run `secureflow refresh` to upgrade. New files are no longer readable with `openssl enc`.
- Metadata moved from `report.txt` to `manifest.json` (an existing `report.txt` is read once and then replaced).
- Decrypted files and `copy_to` targets are written with `0600` permissions, atomically.
- `encrypt` confirms the password interactively and refuses to mix passwords with existing encrypted files.
- `decrypt`/`test` decrypt everything in memory before writing, so a wrong password changes nothing.
- `test` mirrors directory structure under `test_output_dir` (files with the same name no longer collide).
- Prompts go to stderr; one shared stdin reader (piped answers no longer get lost); Ctrl+C at a password prompt restores terminal echo.
- `install.sh` and `install-local` verify SHA-256 checksums before installing.
- CI: actions pinned to commit SHAs, `persist-credentials: false`, govulncheck, Dependabot.
