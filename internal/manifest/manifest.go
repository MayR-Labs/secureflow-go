// Package manifest stores metadata about encrypted files (password hint, note,
// per-file details) as JSON in <output_dir>/manifest.json.
//
// It never stores anything derived from plaintext content (such as a hash),
// since that would let someone confirm guesses about a secret offline.
package manifest

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/MayR-Labs/secureflow-go/internal/utils"
)

const (
	// FileName is the manifest file inside output_dir
	FileName = "manifest.json"
	// LegacyReportName is the plain-text report written by SecureFlow <= 1.2
	LegacyReportName = "report.txt"
	// SchemaVersion is bumped on incompatible manifest changes
	SchemaVersion = 1
	// DefaultNote is used when no note is given
	DefaultNote = "Encrypted secrets for CI/CD"
)

// Manifest describes the encrypted files in an output directory
type Manifest struct {
	SchemaVersion int       `json:"schema_version"`
	Note          string    `json:"note,omitempty"`
	PasswordHint  string    `json:"password_hint,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Files         []File    `json:"files"`
}

// File describes one encrypted file
type File struct {
	Input           string    `json:"input"`
	Encrypted       string    `json:"encrypted"`
	Format          string    `json:"format"`
	Size            int64     `json:"size_bytes"`
	Lines           int       `json:"lines"`
	SourceModified  time.Time `json:"source_modified,omitempty"`
	EncryptedAt     time.Time `json:"encrypted_at"`
	EncryptedSHA256 string    `json:"encrypted_sha256"`
}

// Path returns the manifest path for an output directory
func Path(outputDir string) string {
	return filepath.Join(outputDir, FileName)
}

// Load reads the manifest in outputDir. If none exists it falls back to the
// legacy report.txt (hint and note only), and otherwise returns an empty manifest.
func Load(outputDir string) (*Manifest, error) {
	data, err := os.ReadFile(Path(outputDir))
	if errors.Is(err, os.ErrNotExist) {
		return loadLegacy(outputDir), nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", Path(outputDir), err)
	}
	if m.SchemaVersion > SchemaVersion {
		return nil, fmt.Errorf("%s was written by a newer SecureFlow (schema %d); please upgrade", Path(outputDir), m.SchemaVersion)
	}
	return &m, nil
}

// ReadHint returns the stored password hint, or "" if there is none
func ReadHint(outputDir string) string {
	m, err := Load(outputDir)
	if err != nil {
		return ""
	}
	return m.PasswordHint
}

// Upsert adds or replaces the entry for f.Input
func (m *Manifest) Upsert(f File) {
	for i := range m.Files {
		if m.Files[i].Input == f.Input {
			m.Files[i] = f
			return
		}
	}
	m.Files = append(m.Files, f)
}

// Get returns the entry for input, if present
func (m *Manifest) Get(input string) (File, bool) {
	for _, f := range m.Files {
		if f.Input == input {
			return f, true
		}
	}
	return File{}, false
}

// Prune drops entries whose input is not in keep
func (m *Manifest) Prune(keep []string) {
	set := make(map[string]bool, len(keep))
	for _, k := range keep {
		set[k] = true
	}
	files := m.Files[:0]
	for _, f := range m.Files {
		if set[f.Input] {
			files = append(files, f)
		}
	}
	m.Files = files
}

// Save writes the manifest to outputDir and removes a superseded report.txt
func (m *Manifest) Save(outputDir string) error {
	now := time.Now().UTC().Truncate(time.Second)
	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}
	m.UpdatedAt = now
	m.SchemaVersion = SchemaVersion
	if m.Note == "" {
		m.Note = DefaultNote
	}
	if m.Files == nil {
		m.Files = []File{}
	}
	sort.SliceStable(m.Files, func(i, j int) bool { return m.Files[i].Input < m.Files[j].Input })

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode manifest: %w", err)
	}
	data = append(data, '\n')

	if err := utils.WriteFileAtomic(Path(outputDir), data, 0644); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	// manifest.json replaces the old text report
	_ = os.Remove(filepath.Join(outputDir, LegacyReportName))
	return nil
}

// loadLegacy builds a manifest from a SecureFlow <= 1.2 report.txt
func loadLegacy(outputDir string) *Manifest {
	m := &Manifest{SchemaVersion: SchemaVersion}

	f, err := os.Open(filepath.Join(outputDir, LegacyReportName))
	if err != nil {
		return m
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case strings.HasPrefix(line, "Password Hint:"):
			if hint := strings.TrimSpace(strings.TrimPrefix(line, "Password Hint:")); hint != "N/A" {
				m.PasswordHint = hint
			}
		case strings.HasPrefix(line, "Note:"):
			m.Note = strings.TrimSpace(strings.TrimPrefix(line, "Note:"))
		}
	}
	return m
}
