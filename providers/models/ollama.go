package models

import (
	"bufio"
	"bytes"
	"context"
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
}

// NewOllamaProvider creates a new OllamaProvider configured to use ApiURL.
// The provided logger is attached and the internal context is initialized as an empty slice.
func NewOllamaProvider(logger *logger.Logger) *OllamaProvider {
	return &OllamaProvider{
		Url:     ApiURL,
		logger:  logger,
		context: []int{},
	}
}

func (op *OllamaProvider) GenerateRequest(prompt string) *types.ChatRequest {
	return &types.ChatRequest{
		Model:  "llama3.2",
		Prompt: prompt,
		Stream: true,
	}
}

func (op *OllamaProvider) Chat(ctx context.Context,
	request *types.ChatRequest,
	responseChannel chan *types.ChatResponse,
	errorChannel chan error, doneChannel chan bool,
) {
	if len(op.context) > 0 {
		request.Context = op.context
	}
	data, err := json.Marshal(request)
	op.logger.Debug(fmt.Sprintf("%v", request))
	if err != nil {
		errorChannel <- err
		return
	}
	dataReader := bytes.NewReader(data)
	req, err := http.NewRequestWithContext(ctx, "POST", op.Url, dataReader)
	if err != nil {
		errorChannel <- err
		return
	}
	client := http.Client{}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		errorChannel <- err
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
			responseChannel <- &chunk
		}

		if chunk.Done {
			op.context = chunk.Context
			doneChannel <- true
			return
		}
	}
}
