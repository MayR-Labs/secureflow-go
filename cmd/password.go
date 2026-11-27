package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/MayR-Labs/secureflow-go/internal/utils"
)

const (
	// PasswordEnvVar lets CI pass the password without putting it on the command line
	PasswordEnvVar = "SECUREFLOW_PASSWORD"
	// NewPasswordEnvVar supplies the new password for `secureflow refresh`
	NewPasswordEnvVar = "SECUREFLOW_NEW_PASSWORD"
)

// passwordFromFlagOrEnv returns the password from --password or SECUREFLOW_PASSWORD
func passwordFromFlagOrEnv() string {
	if password != "" {
		return password
	}
	return os.Getenv(PasswordEnvVar)
}

// promptForPassword resolves the password from flag/env or, in interactive mode,
// prompts for it. If hint is non-empty it is shown before prompting.
func promptForPassword(prompt, hint string) (string, error) {
	if pwd := passwordFromFlagOrEnv(); pwd != "" {
		return pwd, nil
	}
	return readPasswordInteractive(prompt, hint)
}

// readPasswordInteractive always prompts (unless --non-interactive)
func readPasswordInteractive(prompt, hint string) (string, error) {
	if nonInteractive {
		return "", fmt.Errorf("password required in non-interactive mode (use --password or set %s)", PasswordEnvVar)
	}

	printHint(hint)

	pwd, err := utils.ReadPassword(utils.Colorize(utils.ColorBlue, prompt))
	if errors.Is(err, utils.ErrNoInput) {
		return "", fmt.Errorf("no password provided: stdin is not a terminal (use --password or set %s)", PasswordEnvVar)
	}
	if err != nil {
		return "", err
	}
	if pwd == "" {
		return "", fmt.Errorf("password cannot be empty")
	}
	return pwd, nil
}

// confirmPassword asks for pwd again and fails if it doesn't match
func confirmPassword(pwd string) error {
	confirm, err := utils.ReadPassword(utils.Colorize(utils.ColorBlue, "🔐 Confirm password: "))
	if err != nil {
		return err
	}
	if confirm != pwd {
		return fmt.Errorf("passwords do not match")
	}
	return nil
}

// printHint prints the password hint, if any
func printHint(hint string) {
	if hint == "" {
		return
	}
	fmt.Fprintln(os.Stderr, utils.Colorize(utils.ColorYellow, "💡 Password hint: "+hint))
}

// promptHint asks for a password hint, offering current as the default.
// Blank keeps current, "-" clears it.
func promptHint(current string) (string, error) {
	prompt := "🔑 (Optional) Enter a password hint (leave blank to skip): "
	if current != "" {
		prompt = fmt.Sprintf("🔑 (Optional) Enter a password hint (leave blank to keep %q, '-' to clear): ", current)
	}
	input, err := utils.ReadLine(utils.Colorize(utils.ColorBlue, prompt))
	if err != nil {
		return "", err
	}
	switch input {
	case "":
		return current, nil
	case "-":
		return "", nil
	default:
		return input, nil
	}
}
