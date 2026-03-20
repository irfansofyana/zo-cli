package api

import "encoding/json"

// --- /zo/ask ---

type AskRequest struct {
	Input          string           `json:"input"`
	ConversationID string           `json:"conversation_id,omitempty"`
	ModelName      string           `json:"model_name,omitempty"`
	PersonaID      string           `json:"persona_id,omitempty"`
	OutputFormat   *json.RawMessage `json:"output_format,omitempty"`
	Stream         bool             `json:"stream"`
}

type AskResponse struct {
	Output         json.RawMessage `json:"output"`
	ConversationID string          `json:"conversation_id,omitempty"`
}

// --- /models/available ---

type Model struct {
	ModelName     string   `json:"model_name"`
	Label         string   `json:"label"`
	Vendor        string   `json:"vendor"`
	Type          *string  `json:"type,omitempty"`
	Priority      *float64 `json:"priority,omitempty"`
	ContextWindow *int     `json:"context_window,omitempty"`
	IsByok        bool     `json:"is_byok"`
}

type ModelsResponse struct {
	Models                []Model `json:"models"`
	FeaturedModelsAreFree bool    `json:"featured_models_are_free"`
}

// --- /personas/available ---

type Persona struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Prompt string  `json:"prompt"`
	Model  *string `json:"model,omitempty"`
	Image  *string `json:"image,omitempty"`
}

type PersonasResponse struct {
	Personas []Persona `json:"personas"`
}

// --- errors ---

type ErrorResponse struct {
	Error string `json:"error"`
}

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return e.Message
}
