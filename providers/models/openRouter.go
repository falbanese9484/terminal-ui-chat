package models

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/falbanese9484/terminal-chat/logger"
	"github.com/falbanese9484/terminal-chat/types"
	_ "github.com/joho/godotenv/autoload"
)

const (
	ORApiURL     = "https://openrouter.ai/api/v1/chat/completions"
	ORRefreshURL = "https://openrouter.ai/api/v1/models"
)

type OpenRouter struct {
	URL            string
	ApiKey         string
	ModelRefresher *types.ModelRefresher
	Model          string
	Context        []OpenRouterMessage
	logger         *logger.Logger
}

func NewOpenRouter(logger *logger.Logger, model string, mf *types.ModelRefresher) (*OpenRouter, error) {
	APIKEY := os.Getenv("OPENROUTER_API_KEY")
	if APIKEY == "" {
		return nil, errors.New("OPENROUTER_API_KEY is required")
	}

	return &OpenRouter{
		URL:            ORApiURL,
		ApiKey:         APIKEY,
		Model:          model,
		logger:         logger,
		ModelRefresher: mf,
	}, nil
}

func (or *OpenRouter) GenerateRequest(prompt string) *types.ChatRequest {
	// Generates the request the way that the Frontend UI expects.
	return &types.ChatRequest{
		Model:  or.Model,
		Prompt: prompt,
		Stream: true,
	}
}

type OpenRouterMessage struct {
	// Used to save context in the form of messages. User or Assistant
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (or *OpenRouter) buildScanner(conn *types.BusConnector,
	msgs []OpenRouterMessage,
) (*http.Response, error) {
	// Builds the *bufio.Scanner for the chat to iterate and read
	request := OpenRouterRequest{
		Model:    or.Model,
		Messages: append(or.Context, msgs...),
		Stream:   true,
	}
	rawReq, err := json.Marshal(&request)
	if err != nil {
		return nil, err
	}
	dataReader := bytes.NewReader(rawReq)
	hReq, err := http.NewRequestWithContext(conn.Ctx, "POST", or.URL, dataReader)
	if err != nil {
		return nil, err
	}
	// TODO: Need better management of timeouts overall
	client := http.Client{
		Timeout: 60 * time.Second,
	}
	hReq.Header.Add("Content-Type", "application/json")
	hReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", or.ApiKey))
	res, err := client.Do(hReq)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, fmt.Errorf("OpenRouter API returned status %d", res.StatusCode)
	}
	return res, nil
}

func (or *OpenRouter) Chat(conn *types.BusConnector) {
	// Handles Streaming LLM responses and forwarding to the ChatBus for the UI
	messages := []OpenRouterMessage{{Role: "user", Content: conn.Request.Prompt}}
	res, err := or.buildScanner(conn, messages)
	if err != nil {
		conn.ErrorChan <- err
		return
	}
	defer res.Body.Close()
	scanner := bufio.NewScanner(res.Body)
	var assistantResponse strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 6 && line[:6] == "data: " {
			data := line[6:]
			if data == "[DONE]" {
				messages = append(messages, OpenRouterMessage{
					Role:    "assistant",
					Content: assistantResponse.String(),
				})
				or.Context = append(or.Context, messages...)
				or.logger.Debug("finished chat stream: %v", messages)
				conn.DoneChannel <- true
				return
			}
			var response OpenRouterStreamResponse
			if err := json.Unmarshal([]byte(data), &response); err != nil {
				conn.ErrorChan <- err
				return
			}
			if len(response.Choices) == 0 {
				conn.DoneChannel <- true
				return
			}
			textResponse := response.Choices[0].Delta.Content
			assistantResponse.WriteString(textResponse)
			returnRes := &types.ChatResponse{Response: textResponse}
			conn.ResponseChan <- returnRes
		}
	}
	if err := scanner.Err(); err != nil {
		conn.ErrorChan <- err
		return
	}
}

func (or *OpenRouter) RetrieveModels() ([]types.Model, error) {
	// TODO: Add context timeout
	if !or.ModelRefresher.IsStale() {
		return or.ModelRefresher.RetrieveModels(), nil
	}
	req, err := http.NewRequest("GET", ORRefreshURL, nil) // TODO: Needs some kind of context
	if err != nil {
		return nil, err
	}
	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, fmt.Errorf("request to retrieve models failed with status: %d", res.StatusCode)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response OpenRouterModelsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	modelsList := []types.Model{}
	for _, v := range response.Data {
		newM := types.Model{
			Name: v.ID,
		}
		modelsList = append(modelsList, newM)
	}
	if err := or.ModelRefresher.StashModels(modelsList); err != nil {
		or.logger.Error("failed to cache models!", "error", err)
		return modelsList, nil
	}
	return modelsList, nil
}

func (or *OpenRouter) SetModel(model string) {
	or.Model = model
}
