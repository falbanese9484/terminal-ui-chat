package chat

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

const ApiURL = "http://localhost:11434/api/generate"

type ChatBus struct {
	// This Bus is going to be used to feed messages to the TUI event loop
	Done    chan bool
	Content chan *ChatResponse
	Error   chan error
}

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

func NewChatBus() *ChatBus {
	return &ChatBus{
		Done:    make(chan bool),
		Content: make(chan *ChatResponse),
		Error:   make(chan error),
	}
}

func (cb *ChatBus) Start(byteReader chan *ChatResponse) {
	for {
		select {
		case response := <-cb.Content:
			byteReader <- response
		case err := <-cb.Error:
			log.Fatalf("%v", err)
			return
		case <-cb.Done:
			return
		}
	}
}

func (cb *ChatBus) RunChat(request *ChatRequest) {
	data, err := json.Marshal(request)
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
			cb.Done <- true
			return
		}
	}
}
