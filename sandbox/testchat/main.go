package main

import (
	"os"

	"github.com/falbanese9484/terminal-chat/chat"
)

func main() {
	args := os.Args
	var prompt string
	if len(args) > 1 {
		prompt = args[1]
	} else {
		prompt = "Hey there! Please give me a recursive function in rust for the fibonacci sequence"
	}
	bus := chat.NewChatBus()
	request := chat.ChatRequest{
		Model:  "llama3.2",
		Prompt: prompt,
		Stream: true,
	}

	go bus.Start()

	bus.RunChat(&request)
	select {}
}
