# Agent Development Guide

## Build/Run Commands
- **Build**: `go build -o bash-butler ./cmd/bash-butler/`
- **Run**: `go run ./cmd/bash-butler/main.go [model-name]`
- **Install**: `./build.sh` (installs to ~/.bash-butler/bin)
- **Test**: No test framework detected - add tests using Go's testing package

## Environment Variables
- `LOG_FILE_PATH=./logs/` - Required for logging
- `OPENROUTER_API_KEY=<key>` - For OpenRouter provider
- `DEBUG=1` - Enable verbose logging

## Code Style Guidelines
- **Module**: `github.com/falbanese9484/terminal-chat`
- **Go Version**: 1.24.0+ required
- **Imports**: Standard library first, third-party, then local packages with alias `uiModels` for ui/models
- **Naming**: CamelCase for exported, camelCase for unexported, interface names without "I" prefix
- **Types**: Group type definitions using `type (...)` block syntax
- **Error Handling**: Return errors explicitly, use `fmt.Errorf` with `%w` for wrapping
- **Logging**: Use structured logging with slog via custom Logger type
- **Comments**: Document exported functions/types, minimal inline comments

## Architecture
- **UI**: Bubble Tea TUI framework with Glamour markdown rendering
- **Providers**: Interface-based design for Ollama/OpenRouter LLM providers
- **Patterns**: Channel-based async communication, component composition