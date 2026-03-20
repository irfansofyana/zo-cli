package cmd

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCmd_HasSubcommands(t *testing.T) {
	names := make([]string, 0)
	for _, c := range rootCmd.Commands() {
		names = append(names, c.Name())
	}
	assert.Contains(t, names, "ask")
	assert.Contains(t, names, "chat")
	assert.Contains(t, names, "models")
	assert.Contains(t, names, "personas")
	assert.Contains(t, names, "config")
	assert.Contains(t, names, "help")
}

func TestRootCmd_HasAPIKeyFlag(t *testing.T) {
	flag := rootCmd.PersistentFlags().Lookup("api-key")
	assert.NotNil(t, flag)
}

func TestRootCmd_UseName(t *testing.T) {
	assert.Equal(t, "zo-cli", rootCmd.Use)
}

func TestRequireAPIKey_ErrorMentionsZoCli(t *testing.T) {
	orig := apiKeyFlag
	apiKeyFlag = ""
	defer func() { apiKeyFlag = orig }()

	t.Setenv("ZO_API_KEY", "")
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(t.TempDir(), "empty"))

	err := requireAPIKey()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "zo-cli config set-key")
}
