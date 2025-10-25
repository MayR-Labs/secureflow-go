package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

const (
	// OpenSSL magic header
	saltedPrefix = "Salted__"
	saltSize     = 8
	keySize      = 32 // AES-256
	ivSize       = 16 // AES block size
	pbkdf2Iter   = 10000
)

// deriveKeyAndIV derives a key and IV from password and salt using PBKDF2
// This matches OpenSSL's key derivation when using -pbkdf2
func deriveKeyAndIV(password, salt []byte) (key, iv []byte) {
	// OpenSSL uses PBKDF2 with SHA256 for -pbkdf2 flag
	keyIV := pbkdf2.Key(password, salt, pbkdf2Iter, keySize+ivSize, sha256.New)
	key = keyIV[:keySize]
	iv = keyIV[keySize:]
	return
}

// EncryptFile encrypts a file using AES-256-CBC in OpenSSL-compatible format
func EncryptFile(inputPath, outputPath, password string) error {
	// Read input file
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Generate random salt
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive key and IV
	key, iv := deriveKeyAndIV([]byte(password), salt)

	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	// Apply PKCS7 padding
	plaintext = pkcs7Pad(plaintext, aes.BlockSize)

	// Encrypt
	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	// Create output: "Salted__" + salt + ciphertext
	output := make([]byte, len(saltedPrefix)+saltSize+len(ciphertext))
	copy(output, []byte(saltedPrefix))
	copy(output[len(saltedPrefix):], salt)
	copy(output[len(saltedPrefix)+saltSize:], ciphertext)

	// Write to file
	if err := os.WriteFile(outputPath, output, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

// DecryptFile decrypts a file that was encrypted using OpenSSL-compatible format
func DecryptFile(inputPath, outputPath, password string) error {
	// Read encrypted file
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read encrypted file: %w", err)
	}

	// Check for "Salted__" prefix
	if len(data) < len(saltedPrefix)+saltSize {
		return fmt.Errorf("invalid encrypted file format")
	}

	prefix := string(data[:len(saltedPrefix)])
	if prefix != saltedPrefix {
		return fmt.Errorf("invalid encrypted file format: missing 'Salted__' prefix")
	}

	// Extract salt and ciphertext
	salt := data[len(saltedPrefix) : len(saltedPrefix)+saltSize]
	ciphertext := data[len(saltedPrefix)+saltSize:]

	// Derive key and IV
	key, iv := deriveKeyAndIV([]byte(password), salt)

	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}

	// Check ciphertext length
	if len(ciphertext)%aes.BlockSize != 0 {
		return fmt.Errorf("ciphertext is not a multiple of block size")
	}

	// Decrypt
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	// Remove PKCS7 padding
	plaintext, err = pkcs7Unpad(plaintext)
	if err != nil {
		return fmt.Errorf("decryption failed (wrong password?): %w", err)
	}

	// Write to file
	if err := os.WriteFile(outputPath, plaintext, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

// pkcs7Pad applies PKCS7 padding to the data
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := make([]byte, padding)
	for i := range padText {
		padText[i] = byte(padding)
	}
	return append(data, padText...)
}

// pkcs7Unpad removes PKCS7 padding from the data
func pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("invalid padding: empty data")
	}

	padding := int(data[length-1])
	if padding > length || padding == 0 {
		return nil, fmt.Errorf("invalid padding")
	}

	// Check padding bytes
	for i := 0; i < padding; i++ {
		if data[length-1-i] != byte(padding) {
			return nil, fmt.Errorf("invalid padding")
		}
	}

	return data[:length-padding], nil
}
