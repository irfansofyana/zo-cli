package cmd

import (
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
	assert.Contains(t, names, "version")
	assert.Contains(t, names, "help")
}

func TestRootCmd_HasAPIKeyFlag(t *testing.T) {
	flag := rootCmd.PersistentFlags().Lookup("api-key")
	assert.NotNil(t, flag)
}
