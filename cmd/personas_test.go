package cmd

import (
	"bytes"
	"context"
	"testing"

	"github.com/irfansofyana/zo-cli/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPersonasListCmd(t *testing.T) {
	mock := &mockClient{
		personasFn: func(ctx context.Context) (*api.PersonasResponse, error) {
			model := "anthropic:claude-sonnet-4"
			return &api.PersonasResponse{
				Personas: []api.Persona{
					{ID: "p1", Name: "Default", Prompt: "Be helpful", Model: &model},
					{ID: "p2", Name: "Coder", Prompt: "Be a coder"},
				},
			}, nil
		},
	}
	cleanup := setMockClient(mock)
	defer cleanup()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"personas", "list"})
	err := rootCmd.Execute()
	require.NoError(t, err)
}

func TestPersonasListCmd_Error(t *testing.T) {
	mock := &mockClient{
		personasFn: func(ctx context.Context) (*api.PersonasResponse, error) {
			return nil, &api.APIError{StatusCode: 500, Message: "server error"}
		},
	}
	cleanup := setMockClient(mock)
	defer cleanup()

	rootCmd.SetArgs([]string{"personas", "list"})
	err := rootCmd.Execute()
	assert.Error(t, err)
}
