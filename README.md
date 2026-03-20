# zo-cli

A command-line tool for interacting with the [Zo Computer API](https://docs.zocomputer.com/api).

## Installation

```bash
go install github.com/irfansofyana/zo-cli@latest
```

Or build from source:

```bash
git clone https://github.com/irfansofyana/zo-cli.git
cd zo-cli
go build -o zo .
```

## Configuration

Set your API key using any of these methods (highest priority first):

1. **Flag**: `--api-key zo_sk_...`
2. **Environment variable**: `export ZO_API_KEY=zo_sk_...`
3. **Config file**: `zo config set-key zo_sk_...`

The config file is stored at `~/.config/zo-cli/config.json`.

```bash
# Set key interactively
zo config set-key

# Set key directly
zo config set-key zo_sk_your_key_here

# View current config
zo config show
```

## Usage

### Ask a single question

```bash
zo ask "What is the meaning of life?"
```

With options:

```bash
# Use a specific model
zo ask --model "anthropic:claude-sonnet-4" "Explain quantum computing"

# Continue a previous conversation
zo ask --conversation-id "conv-abc123" "Tell me more"

# Use a specific persona
zo ask --persona "coder" "Write a fizzbuzz in Go"

# Structured output via JSON schema file
zo ask --output-format schema.json "List 3 colors"
```

The conversation ID is printed to stderr, so you can capture it for scripting:

```bash
CONV_ID=$(zo ask "Hello" 2>&1 1>/dev/null | grep conversation_id | cut -d' ' -f2)
zo ask --conversation-id "$CONV_ID" "Follow up question"
```

### Interactive chat

```bash
zo chat
```

This starts a REPL that automatically maintains conversation context. Type `exit` or `quit` to end.

```bash
# Chat with a specific model
zo chat --model "anthropic:claude-sonnet-4"

# Chat with a specific persona
zo chat --persona "coder"
```

### List models

```bash
zo models list
```

### List personas

```bash
zo personas list
```

## Project Structure

```
zo-cli/
├── main.go              # Entry point
├── api/
│   ├── client.go        # ZoClient interface and HTTP implementation
│   └── types.go         # Request/response types
├── cmd/
│   ├── root.go          # Root command, global flags, API key resolution
│   ├── ask.go           # zo ask
│   ├── chat.go          # zo chat (interactive REPL)
│   ├── models.go        # zo models list
│   ├── personas.go      # zo personas list
│   └── config.go        # zo config set-key / show
└── config/
    └── config.go        # Config file load/save
```

## Development

```bash
# Run tests
go test ./...

# Build
go build -o zo .

# Vet
go vet ./...
```

## Backlog

- **Streaming support (`--stream`)**: Add SSE streaming for `zo ask`. The Zo API's streaming response uses typed SSE events (e.g. `FrontendModelResponse`, `End`) with JSON `data:` payloads containing a `content` field. A proper implementation needs to: parse `event:` lines to distinguish content from control events, decode JSON from `data:` lines and extract the `.content` text, handle end-of-stream events (which carry metadata like `conversation_id`), and suppress non-content payloads from terminal output.

## License

MIT
