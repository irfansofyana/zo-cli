package cmd

import (
	"fmt"
	"os"

	"github.com/irfansofyana/zo-cli/api"
	"github.com/irfansofyana/zo-cli/config"
	"github.com/spf13/cobra"
)

var (
	apiKeyFlag string

	// clientFactory builds the API client. Replaced in tests.
	clientFactory = defaultClientFactory
)

var rootCmd = &cobra.Command{
	Use:           "zo",
	Short:         "Zo CLI - interact with the Zo Computer API",
	Long:          "A command-line tool for chatting with Zo, listing models, and managing personas.",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&apiKeyFlag, "api-key", "", "API key (overrides ZO_API_KEY env and config file)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func getClient() (api.ZoClient, error) {
	return clientFactory()
}

func defaultClientFactory() (api.ZoClient, error) {
	key := resolveAPIKey()
	if key == "" {
		return nil, fmt.Errorf("no API key configured; run 'zo config set-key' or set ZO_API_KEY")
	}

	cfg, _ := config.Load()
	baseURL := ""
	if cfg != nil {
		baseURL = cfg.BaseURL
	}

	return api.NewHTTPClient(baseURL, key), nil
}

func resolveAPIKey() string {
	if apiKeyFlag != "" {
		return apiKeyFlag
	}
	if key := os.Getenv("ZO_API_KEY"); key != "" {
		return key
	}
	cfg, err := config.Load()
	if err != nil {
		return ""
	}
	return cfg.APIKey
}
