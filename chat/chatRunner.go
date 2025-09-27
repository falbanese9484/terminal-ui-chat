package chat

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/falbanese9484/terminal-chat/logger"
)

const ApiURL = "http://localhost:11434/api/generate"

type ChatBus struct {
	// This Bus is going to be used to feed messages to the TUI event loop
	Done    chan bool
	Content chan *ChatResponse
	// NOTE: ChatResponse will need to be able to take in different response specs.
	Error   chan error
	Context []int
	logger  *logger.Logger
}

/*
I need to start thinking about how I'm going to control different Model specs. I know that that
opencode uses a model registry to give users selection.

My system will most likely be largely based on OpenRouter, with an option to use local models.

TODO: Data Model How models will be organized and selectable. Struct Layer / Interface
*/

type ChatRequest struct {
	// Initial structure for the chat request.
	Model   string `json:"model"`
	Prompt  string `json:"prompt"`
	Stream  bool   `json:"stream"`
	Context []int  `json:"context,omitempty"`
}

type ChatResponse struct {
	// What we get back from the LLM Api
	Response string `json:"response"`
	Context  []int  `json:"context"`
	Done     bool   `json:"done"`
}

func NewChatBus(logger *logger.Logger) *ChatBus {
	return &ChatBus{
		Done:    make(chan bool),
		Content: make(chan *ChatResponse),
		Error:   make(chan error),
		logger:  logger,
	}
}

func (cb *ChatBus) Start(byteReader chan *ChatResponse) {
	for {
		select {
		case response := <-cb.Content:
			byteReader <- response
		case err := <-cb.Error:
			cb.logger.Error("failed to read incoming chat response", "error", err)
			return
		case <-cb.Done:
			byteReader <- &ChatResponse{Done: true, Context: cb.Context}
		}
	}
}

func (cb *ChatBus) RunChat(request *ChatRequest) {
	data, err := json.Marshal(request)
	cb.logger.Debug(fmt.Sprintf("%v", request))
	if err != nil {
		cb.Error <- err
		return
	}
	dataReader := bytes.NewReader(data)
	req, err := http.NewRequest("POST", ApiURL, dataReader)
	if err != nil {
		cb.Error <- err
		return
	}
	client := http.Client{}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		cb.Error <- err
		return
	}
	defer res.Body.Close()
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		line := scanner.Text()
		var chunk ChatResponse
		if err := json.Unmarshal([]byte(line), &chunk); err != nil {
			continue
		}

		if chunk.Response != "" {
			cb.Content <- &chunk
		}

		if chunk.Done {
			cb.Context = chunk.Context
			cb.Done <- true
			return
		}
	}
}
