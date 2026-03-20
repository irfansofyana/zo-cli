package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/irfansofyana/zo-cli/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Zo CLI configuration",
}

var configSetKeyCmd = &cobra.Command{
	Use:   "set-key [key]",
	Short: "Set the API key",
	Long:  "Set the API key. Pass as argument or enter interactively.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var key string
		if len(args) > 0 {
			key = args[0]
		} else {
			fmt.Print("Enter API key: ")
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				key = strings.TrimSpace(scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("read input: %w", err)
			}
		}

		if key == "" {
			return fmt.Errorf("API key cannot be empty")
		}

		cfg, err := config.Load()
		if err != nil {
			cfg = &config.Config{}
		}
		cfg.APIKey = key

		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("save config: %w", err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "API key saved to %s\n", config.DefaultPath())
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		out := cmd.OutOrStdout()
		fmt.Fprintf(out, "Config file: %s\n", config.DefaultPath())

		if cfg.APIKey != "" {
			masked := maskKey(cfg.APIKey)
			fmt.Fprintf(out, "API key:     %s\n", masked)
		} else {
			fmt.Fprintf(out, "API key:     (not set)\n")
		}

		if cfg.BaseURL != "" {
			fmt.Fprintf(out, "Base URL:    %s\n", cfg.BaseURL)
		} else {
			fmt.Fprintf(out, "Base URL:    %s (default)\n", "https://api.zo.computer")
		}

		return nil
	},
}

func maskKey(key string) string {
	if len(key) <= 4 {
		return "****"
	}
	return strings.Repeat("*", len(key)-4) + key[len(key)-4:]
}

func init() {
	configCmd.AddCommand(configSetKeyCmd)
	configCmd.AddCommand(configShowCmd)
	rootCmd.AddCommand(configCmd)
}
