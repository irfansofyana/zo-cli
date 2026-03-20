package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAsk_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/zo/ask", r.URL.Path)
		assert.Equal(t, "Bearer test-key", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		body, _ := io.ReadAll(r.Body)
		var req AskRequest
		require.NoError(t, json.Unmarshal(body, &req))
		assert.Equal(t, "hello", req.Input)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"output":          "Hi there!",
			"conversation_id": "conv-123",
		})
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL, "test-key")
	resp, err := client.Ask(context.Background(), AskRequest{Input: "hello"})
	require.NoError(t, err)

	var output string
	require.NoError(t, json.Unmarshal(resp.Output, &output))
	assert.Equal(t, "Hi there!", output)
	assert.Equal(t, "conv-123", resp.ConversationID)
}

func TestAsk_WithOptionalFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req AskRequest
		require.NoError(t, json.Unmarshal(body, &req))
		assert.Equal(t, "test message", req.Input)
		assert.Equal(t, "conv-456", req.ConversationID)
		assert.Equal(t, "anthropic:claude-sonnet-4", req.ModelName)
		assert.Equal(t, "persona-1", req.PersonaID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"output": "response",
		})
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL, "test-key")
	_, err := client.Ask(context.Background(), AskRequest{
		Input:          "test message",
		ConversationID: "conv-456",
		ModelName:      "anthropic:claude-sonnet-4",
		PersonaID:      "persona-1",
	})
	require.NoError(t, err)
}

func TestAsk_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid api key"})
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL, "bad-key")
	_, err := client.Ask(context.Background(), AskRequest{Input: "hello"})
	require.Error(t, err)

	apiErr, ok := err.(*APIError)
	require.True(t, ok)
	assert.Equal(t, http.StatusUnauthorized, apiErr.StatusCode)
	assert.Equal(t, "invalid api key", apiErr.Message)
}

func TestListModels_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/models/available", r.URL.Path)
		assert.Equal(t, "Bearer test-key", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"models": []map[string]interface{}{
				{"model_name": "anthropic:claude-sonnet-4", "label": "Claude Sonnet", "vendor": "Anthropic", "is_byok": false},
				{"model_name": "openai:gpt-4o", "label": "GPT-4o", "vendor": "OpenAI", "is_byok": false},
			},
			"featured_models_are_free": true,
		})
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL, "test-key")
	resp, err := client.ListModels(context.Background())
	require.NoError(t, err)
	assert.Len(t, resp.Models, 2)
	assert.Equal(t, "anthropic:claude-sonnet-4", resp.Models[0].ModelName)
	assert.Equal(t, "Claude Sonnet", resp.Models[0].Label)
	assert.True(t, resp.FeaturedModelsAreFree)
}

func TestListPersonas_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/personas/available", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"personas": []map[string]interface{}{
				{"id": "p1", "name": "Default", "prompt": "You are helpful."},
				{"id": "p2", "name": "Coder", "prompt": "You are a coder."},
			},
		})
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL, "test-key")
	resp, err := client.ListPersonas(context.Background())
	require.NoError(t, err)
	assert.Len(t, resp.Personas, 2)
	assert.Equal(t, "p1", resp.Personas[0].ID)
	assert.Equal(t, "Default", resp.Personas[0].Name)
}

func TestNewHTTPClient_DefaultBaseURL(t *testing.T) {
	client := NewHTTPClient("", "key")
	assert.Equal(t, DefaultBaseURL, client.baseURL)
}
