package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const DefaultBaseURL = "https://api.zo.computer"

// ZoClient is the interface all commands depend on.
type ZoClient interface {
	Ask(ctx context.Context, req AskRequest) (*AskResponse, error)
	AskStream(ctx context.Context, req AskRequest, handler func(chunk string) error) error
	ListModels(ctx context.Context) (*ModelsResponse, error)
	ListPersonas(ctx context.Context) (*PersonasResponse, error)
}

// HTTPClient implements ZoClient using net/http.
type HTTPClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewHTTPClient(baseURL, apiKey string) *HTTPClient {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	return &HTTPClient{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

func (c *HTTPClient) Ask(ctx context.Context, req AskRequest) (*AskResponse, error) {
	req.Stream = false
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/zo/ask", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	c.setHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if json.Unmarshal(respBody, &errResp) == nil && errResp.Error != "" {
			return nil, &APIError{StatusCode: resp.StatusCode, Message: errResp.Error}
		}
		return nil, &APIError{StatusCode: resp.StatusCode, Message: fmt.Sprintf("unexpected status %d: %s", resp.StatusCode, string(respBody))}
	}

	var askResp AskResponse
	if err := json.Unmarshal(respBody, &askResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &askResp, nil
}

func (c *HTTPClient) AskStream(ctx context.Context, req AskRequest, handler func(chunk string) error) error {
	req.Stream = true
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/zo/ask", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	c.setHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		var errResp ErrorResponse
		if json.Unmarshal(respBody, &errResp) == nil && errResp.Error != "" {
			return &APIError{StatusCode: resp.StatusCode, Message: errResp.Error}
		}
		return &APIError{StatusCode: resp.StatusCode, Message: fmt.Sprintf("unexpected status %d: %s", resp.StatusCode, string(respBody))}
	}

	return ReadSSE(resp.Body, handler)
}

func (c *HTTPClient) ListModels(ctx context.Context) (*ModelsResponse, error) {
	var result ModelsResponse
	if err := c.doGet(ctx, "/models/available", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *HTTPClient) ListPersonas(ctx context.Context) (*PersonasResponse, error) {
	var result PersonasResponse
	if err := c.doGet(ctx, "/personas/available", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *HTTPClient) doGet(ctx context.Context, path string, out interface{}) error {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	c.setHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if json.Unmarshal(respBody, &errResp) == nil && errResp.Error != "" {
			return &APIError{StatusCode: resp.StatusCode, Message: errResp.Error}
		}
		return &APIError{StatusCode: resp.StatusCode, Message: fmt.Sprintf("unexpected status %d: %s", resp.StatusCode, string(respBody))}
	}

	if err := json.Unmarshal(respBody, out); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}
	return nil
}

func (c *HTTPClient) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
}
