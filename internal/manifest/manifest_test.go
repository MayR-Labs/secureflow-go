package manifest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadMissingReturnsEmpty(t *testing.T) {
	m, err := Load(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if m.PasswordHint != "" || len(m.Files) != 0 {
		t.Errorf("expected empty manifest, got %+v", m)
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	m := &Manifest{PasswordHint: "the usual", Note: "prod"}
	m.Upsert(File{Input: "b.txt", Encrypted: "b.txt.encrypted", Size: 3, EncryptedAt: time.Now()})
	m.Upsert(File{Input: "a.txt", Encrypted: "a.txt.encrypted"})
	m.Upsert(File{Input: "b.txt", Encrypted: "b.txt.encrypted", Size: 9})
	if err := m.Save(dir); err != nil {
		t.Fatal(err)
	}

	raw, _ := os.ReadFile(Path(dir))
	var generic map[string]any
	if err := json.Unmarshal(raw, &generic); err != nil {
		t.Fatalf("manifest is not valid JSON: %v", err)
	}

	got, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got.PasswordHint != "the usual" || got.Note != "prod" {
		t.Errorf("hint/note not preserved: %+v", got)
	}
	if len(got.Files) != 2 || got.Files[0].Input != "a.txt" {
		t.Fatalf("expected 2 sorted files, got %+v", got.Files)
	}
	if f, _ := got.Get("b.txt"); f.Size != 9 {
		t.Errorf("Upsert did not replace entry: %+v", f)
	}
	if ReadHint(dir) != "the usual" {
		t.Error("ReadHint mismatch")
	}

	got.Prune([]string{"a.txt"})
	if len(got.Files) != 1 {
		t.Errorf("Prune left %d files", len(got.Files))
	}
}

func TestLegacyReportFallbackAndCleanup(t *testing.T) {
	dir := t.TempDir()
	legacy := filepath.Join(dir, LegacyReportName)
	report := "Encryption Report\n=================\n\nNote: old note\nPassword Hint: For the wise only\n"
	if err := os.WriteFile(legacy, []byte(report), 0644); err != nil {
		t.Fatal(err)
	}

	if got := ReadHint(dir); got != "For the wise only" {
		t.Errorf("legacy hint = %q", got)
	}
	m, _ := Load(dir)
	if m.Note != "old note" {
		t.Errorf("legacy note = %q", m.Note)
	}

	if err := m.Save(dir); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(legacy); !os.IsNotExist(err) {
		t.Error("report.txt should be removed once manifest.json is written")
	}

	if err := os.WriteFile(legacy, []byte("Password Hint: N/A\n"), 0644); err != nil {
		t.Fatal(err)
	}
	os.Remove(Path(dir))
	if got := ReadHint(dir); got != "" {
		t.Errorf("N/A hint should be empty, got %q", got)
	}
}

func TestRejectsNewerSchema(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(Path(dir), []byte(`{"schema_version": 99}`), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(dir); err == nil {
		t.Error("expected error for newer schema")
	}
}
