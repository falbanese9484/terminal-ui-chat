package models

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/falbanese9484/terminal-chat/logger"
	"github.com/falbanese9484/terminal-chat/types"
)

const ApiURL = "http://localhost:11434/api/generate"

type OllamaProvider struct {
	Url     string
	logger  *logger.Logger
	context []int
	model   string
}

// NewOllamaProvider creates a new OllamaProvider configured to use ApiURL.
// The provided logger is attached and the internal context is initialized as an empty slice.
func NewOllamaProvider(logger *logger.Logger, model string) *OllamaProvider {
	return &OllamaProvider{
		Url:     ApiURL,
		logger:  logger,
		context: []int{},
		model:   model,
	}
}

func (op *OllamaProvider) GenerateRequest(prompt string) *types.ChatRequest {
	return &types.ChatRequest{
		Model:  op.model,
		Prompt: prompt,
		Stream: true,
	}
}

func (op *OllamaProvider) Chat(connector *types.BusConnector) {
	if len(op.context) > 0 {
		connector.Request.Context = op.context
	}
	data, err := json.Marshal(connector.Request)
	op.logger.Debug(fmt.Sprintf("%v", connector.Request))
	if err != nil {
		connector.ErrorChan <- err
		return
	}
	dataReader := bytes.NewReader(data)
	req, err := http.NewRequestWithContext(connector.Ctx, "POST", op.Url, dataReader)
	if err != nil {
		connector.ErrorChan <- err
		return
	}
	client := http.Client{}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		connector.ErrorChan <- err
		return
	}
	defer res.Body.Close()
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		line := scanner.Text()
		var chunk types.ChatResponse
		if err := json.Unmarshal([]byte(line), &chunk); err != nil {
			continue
		}

		if chunk.Response != "" {
			connector.ResponseChan <- &chunk
		}

		if chunk.Done {
			op.context = chunk.Context
			connector.DoneChannel <- true
			return
		}
	}
}
