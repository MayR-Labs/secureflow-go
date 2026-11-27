// Package crypto implements SecureFlow's file encryption.
//
// New files use the v2 format: AES-256-GCM (authenticated encryption) with a
// key derived via PBKDF2-HMAC-SHA256 (600,000 iterations by default).
//
//	offset  size  field
//	0       6     magic "SFENC" + 0x00
//	6       1     format version (2)
//	7       4     PBKDF2 iterations (uint32, big endian)
//	11      16    salt
//	27      12    GCM nonce
//	39      ...   ciphertext || 16-byte GCM tag
//
// The whole header (bytes 0-38) is passed to GCM as additional authenticated
// data, so tampering with any header field causes decryption to fail.
//
// Files in the legacy v1 format (OpenSSL-compatible "Salted__" AES-256-CBC,
// PBKDF2 10,000 iterations) are still decrypted transparently so existing
// repositories keep working; `secureflow refresh` re-encrypts them as v2.
package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

// Format identifies the on-disk encryption format of a file
type Format string

const (
	// FormatV2 is AES-256-GCM with PBKDF2-SHA256 (current)
	FormatV2 Format = "v2-aes-256-gcm"
	// FormatLegacy is OpenSSL-compatible AES-256-CBC with PBKDF2-SHA256 (read-only)
	FormatLegacy Format = "v1-aes-256-cbc"
)

// DefaultIterations is the PBKDF2 iteration count for new files (OWASP 2023+
// recommendation for PBKDF2-HMAC-SHA256).
const DefaultIterations = 600_000

const (
	v2Version     = 2
	v2SaltSize    = 16
	v2NonceSize   = 12
	v2MagicStr    = "SFENC\x00"
	v2HeaderSize  = len(v2MagicStr) + 1 + 4 + v2SaltSize + v2NonceSize
	keySize       = 32 // AES-256
	minIterations = 100_000
	maxIterations = 10_000_000 // guard against hostile files forcing huge work

	// Legacy (v1) OpenSSL format
	saltedPrefix   = "Salted__"
	saltSize       = 8
	ivSize         = 16
	legacyIter     = 10000
	legacyMinBytes = len(saltedPrefix) + saltSize + aes.BlockSize
)

var v2Magic = []byte(v2MagicStr)

// iterations is the PBKDF2 count used for encryption; tests lower it for speed.
var iterations = DefaultIterations

// ErrDecrypt is returned when decryption fails because the password is wrong
// or the file was modified.
var ErrDecrypt = errors.New("decryption failed: wrong password or the file was modified")

// DetectFormat reports which format an encrypted blob uses
func DetectFormat(data []byte) (Format, error) {
	switch {
	case bytes.HasPrefix(data, v2Magic):
		return FormatV2, nil
	case bytes.HasPrefix(data, []byte(saltedPrefix)):
		return FormatLegacy, nil
	default:
		return "", fmt.Errorf("invalid encrypted file format: not a SecureFlow file")
	}
}

// Encrypt encrypts plaintext with password using the v2 format
func Encrypt(plaintext []byte, password string) ([]byte, error) {
	header := make([]byte, v2HeaderSize)
	off := copy(header, v2Magic)
	header[off] = v2Version
	off++
	binary.BigEndian.PutUint32(header[off:], uint32(iterations))
	off += 4
	salt := header[off : off+v2SaltSize]
	nonce := header[off+v2SaltSize:]
	if _, err := io.ReadFull(rand.Reader, header[off:]); err != nil {
		return nil, fmt.Errorf("failed to generate salt/nonce: %w", err)
	}

	gcm, err := newGCM(password, salt, iterations)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(header, nonce, plaintext, header), nil
}

// Decrypt decrypts data in either the v2 or legacy format
func Decrypt(data []byte, password string) ([]byte, error) {
	format, err := DetectFormat(data)
	if err != nil {
		return nil, err
	}
	if format == FormatLegacy {
		return decryptLegacy(data, password)
	}

	if len(data) < v2HeaderSize+16 {
		return nil, fmt.Errorf("invalid encrypted file format: file is truncated")
	}
	off := len(v2Magic)
	if v := data[off]; v != v2Version {
		return nil, fmt.Errorf("unsupported encrypted file version %d (upgrade secureflow)", v)
	}
	off++
	iter := int(binary.BigEndian.Uint32(data[off:]))
	off += 4
	if iter < minIterations || iter > maxIterations {
		return nil, fmt.Errorf("invalid encrypted file: unsupported iteration count %d", iter)
	}
	salt := data[off : off+v2SaltSize]
	nonce := data[off+v2SaltSize : v2HeaderSize]
	header := data[:v2HeaderSize]

	gcm, err := newGCM(password, salt, iter)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, data[v2HeaderSize:], header)
	if err != nil {
		return nil, ErrDecrypt
	}
	return plaintext, nil
}

func newGCM(password string, salt []byte, iter int) (cipher.AEAD, error) {
	key := pbkdf2.Key([]byte(password), salt, iter, keySize, sha256.New)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}
	return gcm, nil
}

// EncryptFile encrypts inputPath to outputPath using the v2 format
func EncryptFile(inputPath, outputPath, password string) error {
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}
	ciphertext, err := Encrypt(plaintext, password)
	if err != nil {
		return err
	}
	if err := os.WriteFile(outputPath, ciphertext, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}
	return nil
}

// DecryptFile decrypts inputPath (v2 or legacy) to outputPath with 0600 permissions
func DecryptFile(inputPath, outputPath, password string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read encrypted file: %w", err)
	}
	plaintext, err := Decrypt(data, password)
	if err != nil {
		return err
	}
	if err := os.WriteFile(outputPath, plaintext, 0600); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}
	return nil
}

// ---- Legacy v1 (OpenSSL-compatible AES-256-CBC), decrypt only ----

// deriveKeyAndIV matches `openssl enc -aes-256-cbc -pbkdf2` key derivation
func deriveKeyAndIV(password, salt []byte) (key, iv []byte) {
	keyIV := pbkdf2.Key(password, salt, legacyIter, keySize+ivSize, sha256.New)
	return keyIV[:keySize], keyIV[keySize:]
}

func decryptLegacy(data []byte, password string) ([]byte, error) {
	if len(data) < legacyMinBytes {
		return nil, fmt.Errorf("invalid encrypted file format: file is truncated")
	}
	salt := data[len(saltedPrefix) : len(saltedPrefix)+saltSize]
	ciphertext := data[len(saltedPrefix)+saltSize:]
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of block size")
	}

	key, iv := deriveKeyAndIV([]byte(password), salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	plaintext := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plaintext, ciphertext)

	plaintext, err = pkcs7Unpad(plaintext)
	if err != nil {
		return nil, ErrDecrypt
	}
	return plaintext, nil
}

// encryptLegacy produces a v1 file; only used by tests to check backward compatibility
func encryptLegacy(plaintext []byte, password string) ([]byte, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	key, iv := deriveKeyAndIV([]byte(password), salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	padded := pkcs7Pad(plaintext, aes.BlockSize)
	out := make([]byte, len(saltedPrefix)+saltSize+len(padded))
	copy(out, saltedPrefix)
	copy(out[len(saltedPrefix):], salt)
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(out[len(saltedPrefix)+saltSize:], padded)
	return out, nil
}

// pkcs7Pad applies PKCS7 padding to the data
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	return append(data, bytes.Repeat([]byte{byte(padding)}, padding)...)
}

// pkcs7Unpad removes PKCS7 padding from the data
func pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("invalid padding: empty data")
	}
	padding := int(data[length-1])
	if padding > length || padding == 0 || padding > aes.BlockSize {
		return nil, fmt.Errorf("invalid padding")
	}
	for i := 0; i < padding; i++ {
		if data[length-1-i] != byte(padding) {
			return nil, fmt.Errorf("invalid padding")
		}
	}
	return data[:length-padding], nil
}
