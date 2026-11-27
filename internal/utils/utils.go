package utils

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[0;31m"
	ColorGreen  = "\033[0;32m"
	ColorYellow = "\033[1;33m"
	ColorBlue   = "\033[0;34m"
)

// Colorize adds color to text
func Colorize(color, text string) string {
	return color + text + ColorReset
}

// stdin is shared by every prompt. Creating a new bufio.Reader per prompt can
// swallow input that was read ahead into a discarded buffer (e.g. when answers
// are piped in), so all line reads go through this single reader.
var stdin = bufio.NewReader(os.Stdin)

// StdinIsTerminal reports whether stdin is an interactive terminal
func StdinIsTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

// ErrNoInput is returned when a prompt needs input but stdin is closed/empty
var ErrNoInput = errors.New("no input available on stdin")

// ReadPassword reads a password without echoing it. When stdin is not a
// terminal (piped input) it reads a single line instead. If the user presses
// Ctrl+C while typing, the terminal's echo is restored before exiting.
func ReadPassword(prompt string) (string, error) {
	fmt.Fprint(os.Stderr, prompt)

	if !StdinIsTerminal() {
		line, err := readLine()
		fmt.Fprintln(os.Stderr)
		return line, err
	}

	fd := int(os.Stdin.Fd())
	state, err := term.GetState(fd)
	if err != nil {
		return "", fmt.Errorf("failed to read terminal state: %w", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	done := make(chan struct{})
	defer func() {
		signal.Stop(sigCh)
		close(done)
	}()
	go func() {
		select {
		case <-sigCh:
			_ = term.Restore(fd, state)
			fmt.Fprintln(os.Stderr)
			os.Exit(130)
		case <-done:
		}
	}()

	bytePassword, err := term.ReadPassword(fd)
	fmt.Fprintln(os.Stderr) // newline after hidden input
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}
	return string(bytePassword), nil
}

// ReadLine prints prompt and reads a line of text from stdin
func ReadLine(prompt string) (string, error) {
	fmt.Fprint(os.Stderr, prompt)
	return readLine()
}

func readLine() (string, error) {
	text, err := stdin.ReadString('\n')
	if err != nil {
		// Accept a final line without a trailing newline
		if errors.Is(err, io.EOF) && text != "" {
			return strings.TrimSpace(text), nil
		}
		if errors.Is(err, io.EOF) {
			return "", ErrNoInput
		}
		return "", fmt.Errorf("failed to read input: %w", err)
	}
	return strings.TrimSpace(text), nil
}

// Confirm asks a yes/no question; anything but y/yes is "no"
func Confirm(prompt string) (bool, error) {
	answer, err := ReadLine(prompt)
	if err != nil {
		return false, err
	}
	answer = strings.ToLower(answer)
	return answer == "y" || answer == "yes", nil
}

// WriteFileAtomic writes data to a temp file in the same directory, sets perm,
// and renames it over path. Readers never see a half-written file, and perm is
// applied even when path already exists (os.WriteFile only applies perm when
// creating a new file).
func WriteFileAtomic(path string, data []byte, perm os.FileMode) error {
	tmp, err := StageFile(path, data, perm)
	if err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("failed to replace %s: %w", path, err)
	}
	return nil
}

// StageFile writes data to a temp file next to path and returns its name.
// The caller renames it into place (or removes it) later; this lets several
// files be prepared first and swapped in together.
func StageFile(path string, data []byte, perm os.FileMode) (string, error) {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".tmp-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file in %s: %w", dir, err)
	}
	name := tmp.Name()
	fail := func(err error) (string, error) {
		tmp.Close()
		os.Remove(name)
		return "", err
	}
	if err := tmp.Chmod(perm); err != nil && runtime.GOOS != "windows" {
		return fail(fmt.Errorf("failed to set permissions: %w", err))
	}
	if _, err := tmp.Write(data); err != nil {
		return fail(fmt.Errorf("failed to write %s: %w", path, err))
	}
	if err := tmp.Sync(); err != nil {
		return fail(fmt.Errorf("failed to flush %s: %w", path, err))
	}
	if err := tmp.Close(); err != nil {
		os.Remove(name)
		return "", fmt.Errorf("failed to close %s: %w", path, err)
	}
	return name, nil
}

// CountLines counts newline-terminated lines (plus a final unterminated one)
func CountLines(data []byte) int {
	if len(data) == 0 {
		return 0
	}
	n := bytes.Count(data, []byte{'\n'})
	if data[len(data)-1] != '\n' {
		n++
	}
	return n
}

// FileInfo represents metadata about a file
type FileInfo struct {
	Path         string
	Size         int64
	Lines        int
	LastModified time.Time
}

// GetFileInfo retrieves metadata about a file
func GetFileInfo(path string) (*FileInfo, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	// Count lines
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	lines := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return &FileInfo{
		Path:         path,
		Size:         stat.Size(),
		Lines:        lines,
		LastModified: stat.ModTime(),
	}, nil
}

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return nil
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// CopyFile copies src to dst atomically with the given permissions
func CopyFile(src, dst string, perm os.FileMode) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}
	return WriteFileAtomic(dst, data, perm)
}
