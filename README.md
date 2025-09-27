# Terminal Chat

Terminal UI for chatting with LLM models.

## [Demo Video](https://s3.us-east-1.amazonaws.com/images.proaistudios.com/terminal_ui_chat.mov)

Purpose of this build is to familiarize myself with the bubbletea TUI framework.

Right now I have working streaming chat responses with a running local ollama model.

To run:
```bash
export LOG_FILE_PATH=./logs/
go run ./sandbox/ui-chat/main.go
```

The chat supports markdown using charmbracelets glamour library.
