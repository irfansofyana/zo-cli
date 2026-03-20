package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/irfansofyana/zo-cli/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAskCmd_Success(t *testing.T) {
	mock := &mockClient{
		askFunc: func(ctx context.Context, req api.AskRequest) (*api.AskResponse, error) {
			assert.Equal(t, "hello world", req.Input)
			output, _ := json.Marshal("Hi there!")
			return &api.AskResponse{
				Output:         output,
				ConversationID: "conv-123",
			}, nil
		},
	}
	cleanup := setMockClient(mock)
	defer cleanup()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"ask", "hello", "world"})
	err := rootCmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Hi there!")
}

func TestAskCmd_WithConversationID(t *testing.T) {
	mock := &mockClient{
		askFunc: func(ctx context.Context, req api.AskRequest) (*api.AskResponse, error) {
			assert.Equal(t, "conv-456", req.ConversationID)
			assert.Equal(t, "follow up", req.Input)
			output, _ := json.Marshal("Response")
			return &api.AskResponse{Output: output}, nil
		},
	}
	cleanup := setMockClient(mock)
	defer cleanup()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"ask", "--conversation-id", "conv-456", "follow", "up"})
	err := rootCmd.Execute()
	require.NoError(t, err)
}

func TestAskCmd_WithModel(t *testing.T) {
	mock := &mockClient{
		askFunc: func(ctx context.Context, req api.AskRequest) (*api.AskResponse, error) {
			assert.Equal(t, "anthropic:claude-sonnet-4", req.ModelName)
			output, _ := json.Marshal("ok")
			return &api.AskResponse{Output: output}, nil
		},
	}
	cleanup := setMockClient(mock)
	defer cleanup()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"ask", "--model", "anthropic:claude-sonnet-4", "test"})
	err := rootCmd.Execute()
	require.NoError(t, err)
}

func TestAskCmd_Error(t *testing.T) {
	mock := &mockClient{
		askFunc: func(ctx context.Context, req api.AskRequest) (*api.AskResponse, error) {
			return nil, &api.APIError{StatusCode: 401, Message: "unauthorized"}
		},
	}
	cleanup := setMockClient(mock)
	defer cleanup()

	rootCmd.SetArgs([]string{"ask", "hello"})
	err := rootCmd.Execute()
	assert.Error(t, err)
}

func TestAskCmd_StructuredOutput(t *testing.T) {
	mock := &mockClient{
		askFunc: func(ctx context.Context, req api.AskRequest) (*api.AskResponse, error) {
			output, _ := json.Marshal(map[string]string{"name": "Zo"})
			return &api.AskResponse{Output: output}, nil
		},
	}
	cleanup := setMockClient(mock)
	defer cleanup()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"ask", "test"})
	err := rootCmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Zo")
}
