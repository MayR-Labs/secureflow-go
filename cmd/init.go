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
You can modify this file to match your project structure.

Available templates:
  - default (React Native/mobile app)
  - reactnative (React Native specific)
  - flutter (Flutter mobile app)
  - web (Web application)
  - docker (Docker deployment)
  - k8s (Kubernetes deployment)
  - microservices (Microservices architecture)

Use --template flag to specify a template, or leave blank for interactive selection.`,
	RunE: runInit,
}

var templateName string

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&templateName, "template", "", "Config template to use (default, reactnative, flutter, web, docker, k8s, microservices)")
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

	// Determine which template to use
	var cfg *config.Config
	
	if templateName != "" {
		// Template provided via flag
		cfg = config.TemplateConfig(templateName)
		fmt.Printf("%s 📝 Using %s template\n", utils.ColorBlue, templateName)
	} else if !nonInteractive {
		// Interactive template selection
		fmt.Printf("%s 🎨 Select a configuration template:\n\n", utils.ColorBlue)
		fmt.Println("  1. Default (React Native/Mobile App)")
		fmt.Println("  2. React Native")
		fmt.Println("  3. Flutter")
		fmt.Println("  4. Web Application")
		fmt.Println("  5. Docker Deployment")
		fmt.Println("  6. Kubernetes (K8s)")
		fmt.Println("  7. Microservices")
		fmt.Println()
		
		choice, err := utils.ReadLine("Enter your choice (1-7) [1]: ")
		if err != nil {
			return err
		}
		
		// Default to 1 if empty
		if choice == "" {
			choice = "1"
		}
		
		switch choice {
		case "1":
			cfg = config.DefaultConfig()
			templateName = "default"
		case "2":
			cfg = config.ReactNativeConfig()
			templateName = "reactnative"
		case "3":
			cfg = config.FlutterConfig()
			templateName = "flutter"
		case "4":
			cfg = config.WebConfig()
			templateName = "web"
		case "5":
			cfg = config.DockerConfig()
			templateName = "docker"
		case "6":
			cfg = config.K8sConfig()
			templateName = "k8s"
		case "7":
			cfg = config.MicroservicesConfig()
			templateName = "microservices"
		default:
			cfg = config.DefaultConfig()
			templateName = "default"
		}
		
		fmt.Printf("%s 📝 Using %s template\n\n", utils.ColorBlue, templateName)
	} else {
		// Non-interactive mode, use default
		cfg = config.DefaultConfig()
	}

	// Save to file
	if err := cfg.Save(cfgFile); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	fmt.Printf("%s ✅ Created %s\n", utils.ColorGreen, cfgFile)
	fmt.Println("\nYou can now edit this file to match your project structure.")
	fmt.Println("Then run: secureflow encrypt")

	return nil
}
