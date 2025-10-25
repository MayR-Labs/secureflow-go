package crypto

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEncryptDecryptFile(t *testing.T) {
	// Create temp directory for test files
	tmpDir := t.TempDir()
	
	// Create test input file
	inputPath := filepath.Join(tmpDir, "test_input.txt")
	testContent := "This is a test file with sensitive data!\nLine 2\nLine 3"
	if err := os.WriteFile(inputPath, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test input file: %v", err)
	}
	
	// Define encrypted file path
	encryptedPath := filepath.Join(tmpDir, "test_input.txt.encrypted")
	
	// Define decrypted file path
	decryptedPath := filepath.Join(tmpDir, "test_input_decrypted.txt")
	
	// Test password
	password := "test_password_123"
	
	// Test encryption
	t.Run("Encrypt", func(t *testing.T) {
		err := EncryptFile(inputPath, encryptedPath, password)
		if err != nil {
			t.Fatalf("EncryptFile failed: %v", err)
		}
		
		// Verify encrypted file exists
		if _, err := os.Stat(encryptedPath); os.IsNotExist(err) {
			t.Fatal("Encrypted file was not created")
		}
		
		// Verify encrypted file has "Salted__" prefix
		data, err := os.ReadFile(encryptedPath)
		if err != nil {
			t.Fatalf("Failed to read encrypted file: %v", err)
		}
		
		if len(data) < len(saltedPrefix) {
			t.Fatal("Encrypted file is too short")
		}
		
		if string(data[:len(saltedPrefix)]) != saltedPrefix {
			t.Errorf("Expected 'Salted__' prefix, got: %s", string(data[:len(saltedPrefix)]))
		}
	})
	
	// Test decryption
	t.Run("Decrypt", func(t *testing.T) {
		err := DecryptFile(encryptedPath, decryptedPath, password)
		if err != nil {
			t.Fatalf("DecryptFile failed: %v", err)
		}
		
		// Verify decrypted content matches original
		decryptedContent, err := os.ReadFile(decryptedPath)
		if err != nil {
			t.Fatalf("Failed to read decrypted file: %v", err)
		}
		
		if string(decryptedContent) != testContent {
			t.Errorf("Decrypted content doesn't match original.\nExpected: %q\nGot: %q", testContent, string(decryptedContent))
		}
	})
	
	// Test decryption with wrong password
	t.Run("DecryptWrongPassword", func(t *testing.T) {
		wrongPasswordPath := filepath.Join(tmpDir, "test_wrong_password.txt")
		err := DecryptFile(encryptedPath, wrongPasswordPath, "wrong_password")
		if err == nil {
			t.Fatal("Expected error with wrong password, got nil")
		}
	})
}

func TestEncryptFileNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	inputPath := filepath.Join(tmpDir, "nonexistent.txt")
	outputPath := filepath.Join(tmpDir, "output.encrypted")
	
	err := EncryptFile(inputPath, outputPath, "password")
	if err == nil {
		t.Fatal("Expected error when encrypting non-existent file")
	}
}

func TestDecryptFileInvalidFormat(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create invalid encrypted file (no "Salted__" prefix)
	invalidPath := filepath.Join(tmpDir, "invalid.encrypted")
	if err := os.WriteFile(invalidPath, []byte("Invalid content"), 0644); err != nil {
		t.Fatalf("Failed to create invalid file: %v", err)
	}
	
	outputPath := filepath.Join(tmpDir, "output.txt")
	err := DecryptFile(invalidPath, outputPath, "password")
	if err == nil {
		t.Fatal("Expected error when decrypting invalid file")
	}
}

func TestPkcs7Padding(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		blockSize int
		expected  int // expected padded length
	}{
		{
			name:      "Empty data",
			data:      []byte{},
			blockSize: 16,
			expected:  16,
		},
		{
			name:      "One byte",
			data:      []byte{0x01},
			blockSize: 16,
			expected:  16,
		},
		{
			name:      "Full block",
			data:      make([]byte, 16),
			blockSize: 16,
			expected:  32, // Should add full padding block
		},
		{
			name:      "Partial block",
			data:      []byte{0x01, 0x02, 0x03},
			blockSize: 16,
			expected:  16,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			padded := pkcs7Pad(tt.data, tt.blockSize)
			if len(padded) != tt.expected {
				t.Errorf("Expected padded length %d, got %d", tt.expected, len(padded))
			}
			
			// Verify padding is correct
			if len(padded) > 0 {
				paddingValue := padded[len(padded)-1]
				if paddingValue < 1 || paddingValue > byte(tt.blockSize) {
					t.Errorf("Invalid padding value: %d", paddingValue)
				}
			}
			
			// Test unpadding
			unpadded, err := pkcs7Unpad(padded)
			if err != nil {
				t.Fatalf("Unpad failed: %v", err)
			}
			
			if string(unpadded) != string(tt.data) {
				t.Errorf("Unpadded data doesn't match original.\nExpected: %v\nGot: %v", tt.data, unpadded)
			}
		})
	}
}

func TestPkcs7UnpadInvalid(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "Empty data",
			data: []byte{},
		},
		{
			name: "Invalid padding value",
			data: []byte{0x01, 0x02, 0x03, 0x20}, // padding value 32 but only 4 bytes
		},
		{
			name: "Inconsistent padding",
			data: []byte{0x01, 0x02, 0x03, 0x02, 0x03}, // last byte says padding is 3, but not all padding bytes match
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := pkcs7Unpad(tt.data)
			if err == nil {
				t.Error("Expected error for invalid padding")
			}
		})
	}
}

func TestDeriveKeyAndIV(t *testing.T) {
	password := []byte("test_password")
	salt := []byte("12345678") // 8 bytes
	
	key, iv := deriveKeyAndIV(password, salt)
	
	// Verify key size (AES-256 requires 32 bytes)
	if len(key) != keySize {
		t.Errorf("Expected key size %d, got %d", keySize, len(key))
	}
	
	// Verify IV size (AES block size is 16 bytes)
	if len(iv) != ivSize {
		t.Errorf("Expected IV size %d, got %d", ivSize, len(iv))
	}
	
	// Verify deterministic output (same password and salt should give same key/IV)
	key2, iv2 := deriveKeyAndIV(password, salt)
	if string(key) != string(key2) {
		t.Error("Key derivation is not deterministic")
	}
	if string(iv) != string(iv2) {
		t.Error("IV derivation is not deterministic")
	}
	
	// Verify different salt gives different key
	salt2 := []byte("87654321")
	key3, iv3 := deriveKeyAndIV(password, salt2)
	if string(key) == string(key3) {
		t.Error("Different salts should produce different keys")
	}
	if string(iv) == string(iv3) {
		t.Error("Different salts should produce different IVs")
	}
}
