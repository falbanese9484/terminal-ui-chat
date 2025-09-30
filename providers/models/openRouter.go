package models

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/falbanese9484/terminal-chat/logger"
	"github.com/falbanese9484/terminal-chat/types"
	_ "github.com/joho/godotenv/autoload"
)

const ORApiURL = "https://openrouter.ai/api/v1/chat/completions"

type OpenRouter struct {
	URL     string
	ApiKey  string
	Model   string
	Context []OpenRouterMessage
	logger  *logger.Logger
}

func NewOpenRouter(logger *logger.Logger, model string) (*OpenRouter, error) {
	APIKEY := os.Getenv("OPENROUTER_API_KEY")
	if APIKEY == "" {
		return nil, errors.New("OPENROUTER_API_KEY is required")
	}

	return &OpenRouter{
		URL:    ORApiURL,
		ApiKey: APIKEY,
		Model:  model,
		logger: logger,
	}, nil
}

func (or *OpenRouter) GenerateRequest(prompt string) *types.ChatRequest {
	return &types.ChatRequest{
		Model:  or.Model,
		Prompt: prompt,
		Stream: true,
	}
}

type OpenRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenRouterRequest struct {
	Model    string              `json:"model"`
	Messages []OpenRouterMessage `json:"messages"`
	Stream   bool                `json:"stream"`
}

type OpenRouterStreamResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Content string `json:"content,omitempty"`
			Role    string `json:"role,omitempty"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
}

func (or *OpenRouter) Chat(conn *types.BusConnector) {
	messages := append(or.Context, OpenRouterMessage{Role: "user", Content: conn.Request.Prompt})
	request := OpenRouterRequest{
		Model:    or.Model,
		Messages: messages,
		Stream:   true,
	}
	rawReq, err := json.Marshal(&request)
	if err != nil {
		conn.ErrorChan <- err
		return
	}
	dataReader := bytes.NewReader(rawReq)
	hReq, err := http.NewRequestWithContext(conn.Ctx, "POST", or.URL, dataReader)
	if err != nil {
		conn.ErrorChan <- err
		return
	}
	client := http.Client{}
	hReq.Header.Add("Content-Type", "application/json")
	hReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", or.ApiKey))
	res, err := client.Do(hReq)
	if err != nil {
		conn.ErrorChan <- err
		return
	}
	defer res.Body.Close()
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 6 && line[:6] == "data: " {
			data := line[6:]
			if data == "[DONE]" {
				conn.DoneChannel <- true
				return
			}
			var response OpenRouterStreamResponse
			if err := json.Unmarshal([]byte(data), &response); err != nil {
				conn.ErrorChan <- err
				return
			}
			textResponse := response.Choices[0].Delta.Content
			returnRes := &types.ChatResponse{Response: textResponse}
			conn.ResponseChan <- returnRes
		}
	}
}
