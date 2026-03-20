package cmd

import (
	"context"
	"os"

	"github.com/irfansofyana/zo-cli/internal/api"
)

// mockClient implements api.ZoClient for testing.
type mockClient struct {
	askFunc    func(ctx context.Context, req api.AskRequest) (*api.AskResponse, error)
	modelsFn   func(ctx context.Context) (*api.ModelsResponse, error)
	personasFn func(ctx context.Context) (*api.PersonasResponse, error)
}

func (m *mockClient) Ask(ctx context.Context, req api.AskRequest) (*api.AskResponse, error) {
	if m.askFunc != nil {
		return m.askFunc(ctx, req)
	}
	return nil, nil
}

func (m *mockClient) ListModels(ctx context.Context) (*api.ModelsResponse, error) {
	if m.modelsFn != nil {
		return m.modelsFn(ctx)
	}
	return nil, nil
}

func (m *mockClient) ListPersonas(ctx context.Context) (*api.PersonasResponse, error) {
	if m.personasFn != nil {
		return m.personasFn(ctx)
	}
	return nil, nil
}

// setMockClient replaces the client factory for testing and returns a cleanup function.
// It also sets ZO_API_KEY so that requireAPIKey() passes.
func setMockClient(mock *mockClient) func() {
	old := clientFactory
	clientFactory = func() (api.ZoClient, error) {
		return mock, nil
	}
	prevKey := os.Getenv("ZO_API_KEY")
	os.Setenv("ZO_API_KEY", "test-key")
	return func() {
		clientFactory = old
		os.Setenv("ZO_API_KEY", prevKey)
	}
}
