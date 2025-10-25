package cmd

import (
	"fmt"

	"github.com/MayR-Labs/secureflow-go/internal/config"
	"github.com/MayR-Labs/secureflow-go/internal/utils"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new secureflow.yaml configuration file",
	Long: `Creates a default secureflow.yaml file in the current directory.
You can modify this file to match your project structure.`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	// Check if config already exists
	if utils.FileExists(cfgFile) {
		fmt.Printf("%s ⚠️  Config file already exists: %s\n", utils.ColorYellow, cfgFile)

		if !nonInteractive {
			response, err := utils.ReadLine("Overwrite? (y/N): ")
			if err != nil {
				return err
			}
			if response != "y" && response != "Y" {
				fmt.Println("Aborted.")
				return nil
			}
		} else {
			return fmt.Errorf("config file already exists")
		}
	}

	// Create default config
	cfg := config.DefaultConfig()

	// Save to file
	if err := cfg.Save(cfgFile); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	fmt.Printf("%s ✅ Created %s\n", utils.ColorGreen, cfgFile)
	fmt.Println("\nYou can now edit this file to match your project structure.")
	fmt.Println("Then run: secureflow encrypt")

	return nil
}
