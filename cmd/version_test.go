package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCmd_HasVersion(t *testing.T) {
	assert.NotEmpty(t, rootCmd.Version)
}

func TestVersion_NeverEmpty(t *testing.T) {
	assert.NotEmpty(t, Version)
}
