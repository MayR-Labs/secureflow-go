package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/MayR-Labs/secureflow-go/internal/config"
	"github.com/MayR-Labs/secureflow-go/internal/crypto"
	"github.com/MayR-Labs/secureflow-go/internal/manifest"
	"github.com/MayR-Labs/secureflow-go/internal/utils"
)

const (
	// Decrypted secrets are readable only by the current user
	decryptedPerm os.FileMode = 0600
	// Encrypted files are safe to share and are committed to git
	encryptedPerm os.FileMode = 0644
)

// encryptedPath returns where fm's encrypted file lives
func encryptedPath(cfg *config.Config, fm config.FileMapping) string {
	return filepath.Join(cfg.OutputDir, fm.EncryptedName())
}

// sealed is an encrypted file ready to be written, plus its manifest entry
type sealed struct {
	path  string
	data  []byte
	entry manifest.File
}

// seal encrypts plaintext for fm and builds its manifest entry
func seal(cfg *config.Config, fm config.FileMapping, plaintext []byte, pwd string, sourceModified time.Time) (*sealed, error) {
	ciphertext, err := crypto.Encrypt(plaintext, pwd)
	if err != nil {
		return nil, err
	}
	sum := sha256.Sum256(ciphertext)
	return &sealed{
		path: encryptedPath(cfg, fm),
		data: ciphertext,
		entry: manifest.File{
			Input:           fm.Input,
			Encrypted:       fm.EncryptedName(),
			Format:          string(crypto.FormatV2),
			Size:            int64(len(plaintext)),
			Lines:           utils.CountLines(plaintext),
			SourceModified:  sourceModified.UTC().Truncate(time.Second),
			EncryptedAt:     time.Now().UTC().Truncate(time.Second),
			EncryptedSHA256: hex.EncodeToString(sum[:]),
		},
	}, nil
}

// sealInput reads fm.Input from disk and seals it
func sealInput(cfg *config.Config, fm config.FileMapping, pwd string) (*sealed, error) {
	info, err := os.Stat(fm.Input)
	if err != nil {
		return nil, err
	}
	plaintext, err := os.ReadFile(fm.Input)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", fm.Input, err)
	}
	return seal(cfg, fm, plaintext, pwd, info.ModTime())
}

// writeSealed writes encrypted files all-or-nothing: every file is staged to a
// temp file first and only then renamed into place.
func writeSealed(files []*sealed) error {
	staged := make([]string, 0, len(files))
	cleanup := func() {
		for _, tmp := range staged {
			os.Remove(tmp)
		}
	}
	for _, f := range files {
		if err := utils.EnsureDir(filepath.Dir(f.path)); err != nil {
			cleanup()
			return err
		}
		tmp, err := utils.StageFile(f.path, f.data, encryptedPerm)
		if err != nil {
			cleanup()
			return err
		}
		staged = append(staged, tmp)
	}
	for i, f := range files {
		if err := os.Rename(staged[i], f.path); err != nil {
			cleanup()
			return fmt.Errorf("failed to write %s: %w", f.path, err)
		}
	}
	return nil
}

// openEncrypted decrypts fm's encrypted file in memory
func openEncrypted(cfg *config.Config, fm config.FileMapping, pwd string) ([]byte, crypto.Format, error) {
	data, err := os.ReadFile(encryptedPath(cfg, fm))
	if err != nil {
		return nil, "", err
	}
	format, err := crypto.DetectFormat(data)
	if err != nil {
		return nil, "", err
	}
	plaintext, err := crypto.Decrypt(data, pwd)
	return plaintext, format, err
}

// writePlaintext writes decrypted content to path with owner-only permissions
func writePlaintext(path string, plaintext []byte) error {
	if err := utils.EnsureDir(filepath.Dir(path)); err != nil {
		return err
	}
	return utils.WriteFileAtomic(path, plaintext, decryptedPerm)
}

// verifyPassword checks pwd against existing encrypted files (skipping
// inputs in skip). It returns nil if none exist yet.
func verifyPassword(cfg *config.Config, pwd string, skip map[string]bool) error {
	for _, fm := range cfg.Files {
		if skip[fm.Input] || !utils.FileExists(encryptedPath(cfg, fm)) {
			continue
		}
		_, _, err := openEncrypted(cfg, fm, pwd)
		if errors.Is(err, crypto.ErrDecrypt) {
			return fmt.Errorf("%w (checked %s); use 'secureflow refresh' to change the password", errPasswordMismatch, encryptedPath(cfg, fm))
		}
		if err != nil {
			return err
		}
		return nil // one successful check is enough
	}
	return nil
}

// errPasswordMismatch means a password doesn't open the existing encrypted files
var errPasswordMismatch = errors.New("password does not match the existing encrypted files")

// isWrongPassword reports whether err means the password was wrong
func isWrongPassword(err error) bool {
	return errors.Is(err, crypto.ErrDecrypt)
}
