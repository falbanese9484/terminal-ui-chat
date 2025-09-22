package chat

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const ApiURL = "http://localhost:11434/api/generate"

type ChatBus struct {
	Done    chan bool
	Content chan *ChatResponse
	Error   chan error
}

type ChatRequest struct {
	Model   string `json:"model"`
	Prompt  string `json:"prompt"`
	Stream  bool   `json:"stream"`
	Context []int  `json:"context,omitempty"`
}

type ChatResponse struct {
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

func (cb *ChatBus) Start() {
	for {
		select {
		case response := <-cb.Content:
			fmt.Printf("%s", response.Response)
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
