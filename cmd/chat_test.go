package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/irfansofyana/zo-cli/internal/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatLoop_PromptPrefix(t *testing.T) {
	mock := &mockClient{
		askFunc: func(ctx context.Context, req api.AskRequest) (*api.AskResponse, error) {
			output, _ := json.Marshal("hello back")
			return &api.AskResponse{Output: output}, nil
		},
	}

	input := "hi\nexit\n"
	out := new(bytes.Buffer)
	err := chatLoop(context.Background(), mock, "", "", strings.NewReader(input), out)
	require.NoError(t, err)
	assert.Contains(t, out.String(), "zo-cli>")
}

func TestChatLoop_BasicConversation(t *testing.T) {
	callCount := 0
	mock := &mockClient{
		askFunc: func(ctx context.Context, req api.AskRequest) (*api.AskResponse, error) {
			callCount++
			output, _ := json.Marshal("Response " + req.Input)
			convID := "conv-1"
			return &api.AskResponse{
				Output:         output,
				ConversationID: convID,
			}, nil
		},
	}

	input := "hello\nworld\nexit\n"
	out := new(bytes.Buffer)
	err := chatLoop(context.Background(), mock, "", "", strings.NewReader(input), out)
	require.NoError(t, err)

	assert.Equal(t, 2, callCount)
	assert.Contains(t, out.String(), "Response hello")
	assert.Contains(t, out.String(), "Response world")
	assert.Contains(t, out.String(), "Goodbye!")
}

func TestChatLoop_ConversationIDPersists(t *testing.T) {
	var receivedConvIDs []string
	mock := &mockClient{
		askFunc: func(ctx context.Context, req api.AskRequest) (*api.AskResponse, error) {
			receivedConvIDs = append(receivedConvIDs, req.ConversationID)
			output, _ := json.Marshal("ok")
			return &api.AskResponse{
				Output:         output,
				ConversationID: "conv-abc",
			}, nil
		},
	}

	input := "first\nsecond\nthird\nexit\n"
	out := new(bytes.Buffer)
	err := chatLoop(context.Background(), mock, "", "", strings.NewReader(input), out)
	require.NoError(t, err)

	// First call should have no conversation ID
	assert.Equal(t, "", receivedConvIDs[0])
	// Subsequent calls should have the conversation ID from previous response
	assert.Equal(t, "conv-abc", receivedConvIDs[1])
	assert.Equal(t, "conv-abc", receivedConvIDs[2])
}

func TestChatLoop_SkipsEmptyInput(t *testing.T) {
	callCount := 0
	mock := &mockClient{
		askFunc: func(ctx context.Context, req api.AskRequest) (*api.AskResponse, error) {
			callCount++
			output, _ := json.Marshal("ok")
			return &api.AskResponse{Output: output}, nil
		},
	}

	input := "\n\nhello\n\nexit\n"
	out := new(bytes.Buffer)
	err := chatLoop(context.Background(), mock, "", "", strings.NewReader(input), out)
	require.NoError(t, err)
	assert.Equal(t, 1, callCount)
}

func TestChatLoop_PassesModelAndPersona(t *testing.T) {
	mock := &mockClient{
		askFunc: func(ctx context.Context, req api.AskRequest) (*api.AskResponse, error) {
			assert.Equal(t, "custom-model", req.ModelName)
			assert.Equal(t, "custom-persona", req.PersonaID)
			output, _ := json.Marshal("ok")
			return &api.AskResponse{Output: output}, nil
		},
	}

	input := "test\nexit\n"
	out := new(bytes.Buffer)
	err := chatLoop(context.Background(), mock, "custom-model", "custom-persona", strings.NewReader(input), out)
	require.NoError(t, err)
}

func TestChatLoop_HandlesAPIError(t *testing.T) {
	callCount := 0
	mock := &mockClient{
		askFunc: func(ctx context.Context, req api.AskRequest) (*api.AskResponse, error) {
			callCount++
			if callCount == 1 {
				return nil, &api.APIError{StatusCode: 500, Message: "server error"}
			}
			output, _ := json.Marshal("ok")
			return &api.AskResponse{Output: output}, nil
		},
	}

	input := "first\nsecond\nexit\n"
	out := new(bytes.Buffer)
	err := chatLoop(context.Background(), mock, "", "", strings.NewReader(input), out)
	require.NoError(t, err)

	// Should continue after error
	assert.Equal(t, 2, callCount)
	assert.Contains(t, out.String(), "Error: server error")
	assert.Contains(t, out.String(), "ok")
}

func TestChatLoop_QuitCommand(t *testing.T) {
	mock := &mockClient{}
	input := "quit\n"
	out := new(bytes.Buffer)
	err := chatLoop(context.Background(), mock, "", "", strings.NewReader(input), out)
	require.NoError(t, err)
	assert.Contains(t, out.String(), "Goodbye!")
}

func TestChatLoop_EOF(t *testing.T) {
	callCount := 0
	mock := &mockClient{
		askFunc: func(ctx context.Context, req api.AskRequest) (*api.AskResponse, error) {
			callCount++
			output, _ := json.Marshal("ok")
			return &api.AskResponse{Output: output}, nil
		},
	}

	// No exit command, just EOF
	input := "hello\n"
	out := new(bytes.Buffer)
	err := chatLoop(context.Background(), mock, "", "", strings.NewReader(input), out)
	require.NoError(t, err)
	assert.Equal(t, 1, callCount)
}
