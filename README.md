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
go build -o zo-cli .
```

## Configuration

Set your API key using any of these methods (highest priority first):

1. **Flag**: `--api-key zo_sk_...`
2. **Environment variable**: `export ZO_API_KEY=zo_sk_...`
3. **Config file**: `zo-cli config set-key zo_sk_...`

The config file is stored at `~/.config/zo-cli/config.json`.

```bash
# Set key interactively
zo-cli config set-key

# Set key directly
zo-cli config set-key zo_sk_your_key_here

# View current config
zo-cli config show
```

## Usage

### Ask a single question

```bash
zo-cli ask "What is the meaning of life?"
```

With options:

```bash
# Use a specific model
zo-cli ask --model "anthropic:claude-sonnet-4" "Explain quantum computing"

# Continue a previous conversation
zo-cli ask --conversation-id "conv-abc123" "Tell me more"

# Use a specific persona
zo-cli ask --persona "coder" "Write a fizzbuzz in Go"

# Structured output via JSON schema file
zo-cli ask --output-format schema.json "List 3 colors"
```

The conversation ID is printed to stderr, so you can capture it for scripting:

```bash
CONV_ID=$(zo-cli ask "Hello" 2>&1 1>/dev/null | grep conversation_id | cut -d' ' -f2)
zo-cli ask --conversation-id "$CONV_ID" "Follow up question"
```

### Interactive chat

```bash
zo-cli chat
```

This starts a REPL that automatically maintains conversation context. Type `exit` or `quit` to end.

```bash
# Chat with a specific model
zo-cli chat --model "anthropic:claude-sonnet-4"

# Chat with a specific persona
zo-cli chat --persona "coder"
```

### List models

```bash
zo-cli models list
```

### List personas

```bash
zo-cli personas list
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
│   ├── ask.go           # zo-cli ask
│   ├── chat.go          # zo-cli chat (interactive REPL)
│   ├── models.go        # zo-cli models list
│   ├── personas.go      # zo-cli personas list
│   └── config.go        # zo-cli config set-key / show
└── config/
    └── config.go        # Config file load/save
```

## Development

```bash
# Run tests
go test ./...

# Build
go build -o zo-cli .

# Vet
go vet ./...
```

## Backlog

- **Streaming support (`--stream`)**: Add SSE streaming for `zo-cli ask`.

## License

MIT
