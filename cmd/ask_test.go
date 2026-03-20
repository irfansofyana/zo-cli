package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/irfansofyana/zo-cli/internal/api"
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

func TestAskCmd_OutputFormatSendsSchema(t *testing.T) {
	schema := `{"type":"object","properties":{"name":{"type":"string"}}}`
	schemaFile := filepath.Join(t.TempDir(), "schema.json")
	require.NoError(t, os.WriteFile(schemaFile, []byte(schema), 0644))

	mock := &mockClient{
		askFunc: func(ctx context.Context, req api.AskRequest) (*api.AskResponse, error) {
			require.NotNil(t, req.OutputFormat, "OutputFormat should be set")
			var parsed map[string]interface{}
			require.NoError(t, json.Unmarshal(*req.OutputFormat, &parsed))
			assert.Equal(t, "object", parsed["type"])
			output, _ := json.Marshal(map[string]string{"name": "Zo"})
			return &api.AskResponse{Output: output}, nil
		},
	}
	cleanup := setMockClient(mock)
	defer cleanup()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"ask", "--output-format", schemaFile, "test"})
	err := rootCmd.Execute()
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Zo")
}

func TestAskCmd_OutputFormatInvalidFile(t *testing.T) {
	cleanup := setMockClient(&mockClient{})
	defer cleanup()

	rootCmd.SetArgs([]string{"ask", "--output-format", "/nonexistent/schema.json", "test"})
	err := rootCmd.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "read output format file")
}

func TestAskCmd_OutputFormatOmittedWhenNotSet(t *testing.T) {
	askOutputFormat = "" // reset from prior tests
	mock := &mockClient{
		askFunc: func(ctx context.Context, req api.AskRequest) (*api.AskResponse, error) {
			assert.Nil(t, req.OutputFormat, "OutputFormat should be nil when flag is not set")
			output, _ := json.Marshal("ok")
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
}
