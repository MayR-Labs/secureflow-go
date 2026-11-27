// Package gitutil holds the small git helpers used by `secureflow add`.
package gitutil

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const header = "# SecureFlow: plaintext secrets (encrypted copies live in output_dir)"

func gitAvailable() bool {
	if _, err := exec.LookPath("git"); err != nil {
		return false
	}
	return exec.Command("git", "rev-parse", "--is-inside-work-tree").Run() == nil
}

// IsTracked reports whether path is already committed/staged in git
func IsTracked(path string) bool {
	if !gitAvailable() {
		return false
	}
	return exec.Command("git", "ls-files", "--error-unmatch", "--", path).Run() == nil
}

// IsIgnored reports whether path is already ignored. It asks git when
// available and otherwise looks for an exact entry in .gitignore.
func IsIgnored(gitignore, path string) bool {
	if gitAvailable() {
		// --no-index also checks paths that are already tracked
		return exec.Command("git", "check-ignore", "-q", "--no-index", "--", path).Run() == nil
	}
	lines, _ := readLines(gitignore)
	clean := filepath.ToSlash(filepath.Clean(path))
	for _, l := range lines {
		if l == clean || l == "/"+clean {
			return true
		}
	}
	return false
}

// AddToGitignore appends root-anchored entries for paths not already ignored.
// It returns the entries it added.
func AddToGitignore(gitignore string, paths []string) ([]string, error) {
	var add []string
	for _, p := range paths {
		clean := filepath.ToSlash(filepath.Clean(p))
		if strings.HasPrefix(clean, "../") || filepath.IsAbs(p) {
			continue // outside the repository; can't be ignored from here
		}
		if IsIgnored(gitignore, clean) {
			continue
		}
		add = append(add, "/"+clean)
	}
	if len(add) == 0 {
		return nil, nil
	}

	existing, err := os.ReadFile(gitignore)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	var b strings.Builder
	if len(existing) > 0 && !strings.HasSuffix(string(existing), "\n") {
		b.WriteString("\n")
	}
	if !strings.Contains(string(existing), header) {
		if len(existing) > 0 {
			b.WriteString("\n")
		}
		b.WriteString(header + "\n")
	}
	for _, a := range add {
		b.WriteString(a + "\n")
	}

	f, err := os.OpenFile(gitignore, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if _, err := f.WriteString(b.String()); err != nil {
		return nil, err
	}
	return add, nil
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var lines []string
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, strings.TrimSpace(s.Text()))
	}
	return lines, s.Err()
}
