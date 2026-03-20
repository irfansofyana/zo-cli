package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFrom_MissingFile(t *testing.T) {
	cfg, err := LoadFrom("/nonexistent/path/config.json")
	require.NoError(t, err)
	assert.Equal(t, "", cfg.APIKey)
	assert.Equal(t, "", cfg.BaseURL)
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "zo-cli", "config.json")

	original := &Config{
		APIKey:  "zo_sk_test123",
		BaseURL: "https://custom.api.com",
	}

	err := SaveTo(original, path)
	require.NoError(t, err)

	loaded, err := LoadFrom(path)
	require.NoError(t, err)
	assert.Equal(t, original.APIKey, loaded.APIKey)
	assert.Equal(t, original.BaseURL, loaded.BaseURL)
}

func TestSaveCreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "deep", "config.json")

	err := SaveTo(&Config{APIKey: "test"}, path)
	require.NoError(t, err)

	_, err = os.Stat(path)
	assert.NoError(t, err)
}

func TestLoadFrom_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	err := os.WriteFile(path, []byte("not json"), 0600)
	require.NoError(t, err)

	_, err = LoadFrom(path)
	assert.Error(t, err)
}

func TestSave_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	err := SaveTo(&Config{APIKey: "secret"}, path)
	require.NoError(t, err)

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
}

func TestLoadFrom_EmptyAPIKey(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	err := os.WriteFile(path, []byte(`{"base_url":"https://example.com"}`), 0600)
	require.NoError(t, err)

	cfg, err := LoadFrom(path)
	require.NoError(t, err)
	assert.Equal(t, "", cfg.APIKey)
	assert.Equal(t, "https://example.com", cfg.BaseURL)
}
