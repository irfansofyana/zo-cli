package cmd

import (
	"bytes"
	"context"
	"testing"

	"github.com/irfansofyana/zo-cli/internal/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModelsListCmd(t *testing.T) {
	mType := "free"
	mock := &mockClient{
		modelsFn: func(ctx context.Context) (*api.ModelsResponse, error) {
			return &api.ModelsResponse{
				Models: []api.Model{
					{ModelName: "anthropic:claude-sonnet-4", Label: "Claude Sonnet", Vendor: "Anthropic", Type: &mType, IsByok: false},
					{ModelName: "openai:gpt-4o", Label: "GPT-4o", Vendor: "OpenAI", Type: &mType, IsByok: false},
				},
				FeaturedModelsAreFree: true,
			}, nil
		},
	}
	cleanup := setMockClient(mock)
	defer cleanup()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"models", "list"})
	err := rootCmd.Execute()
	require.NoError(t, err)
}

func TestModelsListCmd_Error(t *testing.T) {
	mock := &mockClient{
		modelsFn: func(ctx context.Context) (*api.ModelsResponse, error) {
			return nil, &api.APIError{StatusCode: 401, Message: "unauthorized"}
		},
	}
	cleanup := setMockClient(mock)
	defer cleanup()

	rootCmd.SetArgs([]string{"models", "list"})
	err := rootCmd.Execute()
	assert.Error(t, err)
}
