package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionCmd_IsRegistered(t *testing.T) {
	names := make([]string, 0)
	for _, c := range rootCmd.Commands() {
		names = append(names, c.Name())
	}
	assert.Contains(t, names, "version")
}

func TestVersionCmd_PrintsVersion(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"version"})
	err := rootCmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), Version)
}

func TestRootCmd_HasVersion(t *testing.T) {
	assert.NotEmpty(t, rootCmd.Version)
}

func TestVersion_NeverEmpty(t *testing.T) {
	// Version must always be non-empty: either injected via ldflags,
	// read from build info (go install @version), or the "dev" fallback.
	assert.NotEmpty(t, Version)
}
