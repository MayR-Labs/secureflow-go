package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestColorize(t *testing.T) {
	tests := []struct {
		name     string
		color    string
		text     string
		expected string
	}{
		{
			name:     "Red text",
			color:    ColorRed,
			text:     "Error",
			expected: ColorRed + "Error" + ColorReset,
		},
		{
			name:     "Green text",
			color:    ColorGreen,
			text:     "Success",
			expected: ColorGreen + "Success" + ColorReset,
		},
		{
			name:     "Empty text",
			color:    ColorBlue,
			text:     "",
			expected: ColorBlue + ColorReset,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Colorize(tt.color, tt.text)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Test existing file
	existingFile := filepath.Join(tmpDir, "existing.txt")
	if err := os.WriteFile(existingFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	if !FileExists(existingFile) {
		t.Error("FileExists returned false for existing file")
	}
	
	// Test non-existent file
	nonExistentFile := filepath.Join(tmpDir, "nonexistent.txt")
	if FileExists(nonExistentFile) {
		t.Error("FileExists returned true for non-existent file")
	}
	
	// Test directory
	if !FileExists(tmpDir) {
		t.Error("FileExists returned false for existing directory")
	}
}

func TestEnsureDir(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Test creating new directory
	newDir := filepath.Join(tmpDir, "newdir")
	err := EnsureDir(newDir)
	if err != nil {
		t.Fatalf("EnsureDir failed: %v", err)
	}
	
	// Verify directory was created
	stat, err := os.Stat(newDir)
	if err != nil {
		t.Fatalf("Directory was not created: %v", err)
	}
	
	if !stat.IsDir() {
		t.Error("Created path is not a directory")
	}
	
	// Test creating nested directories
	nestedDir := filepath.Join(tmpDir, "level1", "level2", "level3")
	err = EnsureDir(nestedDir)
	if err != nil {
		t.Fatalf("EnsureDir failed for nested path: %v", err)
	}
	
	if !FileExists(nestedDir) {
		t.Error("Nested directory was not created")
	}
	
	// Test calling EnsureDir on existing directory (should not error)
	err = EnsureDir(newDir)
	if err != nil {
		t.Errorf("EnsureDir failed on existing directory: %v", err)
	}
}

func TestGetFileInfo(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create test file with known content
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "Line 1\nLine 2\nLine 3\nLine 4\n"
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// Get file info
	info, err := GetFileInfo(testFile)
	if err != nil {
		t.Fatalf("GetFileInfo failed: %v", err)
	}
	
	// Verify path
	if info.Path != testFile {
		t.Errorf("Expected path %s, got %s", testFile, info.Path)
	}
	
	// Verify size
	if info.Size != int64(len(content)) {
		t.Errorf("Expected size %d, got %d", len(content), info.Size)
	}
	
	// Verify line count
	expectedLines := 4
	if info.Lines != expectedLines {
		t.Errorf("Expected %d lines, got %d", expectedLines, info.Lines)
	}
	
	// Verify LastModified is recent (within last minute)
	timeDiff := time.Since(info.LastModified)
	if timeDiff > time.Minute {
		t.Errorf("LastModified is too old: %v", info.LastModified)
	}
}

func TestGetFileInfoEmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create empty file
	emptyFile := filepath.Join(tmpDir, "empty.txt")
	if err := os.WriteFile(emptyFile, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}
	
	info, err := GetFileInfo(emptyFile)
	if err != nil {
		t.Fatalf("GetFileInfo failed on empty file: %v", err)
	}
	
	if info.Size != 0 {
		t.Errorf("Expected size 0, got %d", info.Size)
	}
	
	if info.Lines != 0 {
		t.Errorf("Expected 0 lines, got %d", info.Lines)
	}
}

func TestGetFileInfoNonExistent(t *testing.T) {
	_, err := GetFileInfo("/nonexistent/file.txt")
	if err == nil {
		t.Fatal("Expected error for non-existent file")
	}
}

func TestGetFileInfoSingleLineNoNewline(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create file with single line and no trailing newline
	testFile := filepath.Join(tmpDir, "single.txt")
	content := "Single line without newline"
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	info, err := GetFileInfo(testFile)
	if err != nil {
		t.Fatalf("GetFileInfo failed: %v", err)
	}
	
	// Should count as 1 line
	if info.Lines != 1 {
		t.Errorf("Expected 1 line, got %d", info.Lines)
	}
}

func TestColorConstants(t *testing.T) {
	// Verify color constants are not empty
	colors := map[string]string{
		"ColorReset":  ColorReset,
		"ColorRed":    ColorRed,
		"ColorGreen":  ColorGreen,
		"ColorYellow": ColorYellow,
		"ColorBlue":   ColorBlue,
	}
	
	for name, color := range colors {
		if color == "" {
			t.Errorf("%s is empty", name)
		}
		
		// Verify it's an ANSI escape sequence
		if !strings.HasPrefix(color, "\033[") {
			t.Errorf("%s doesn't start with ANSI escape: %q", name, color)
		}
	}
}
