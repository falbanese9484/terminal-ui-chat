package models

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/falbanese9484/terminal-chat/logger"
	"github.com/falbanese9484/terminal-chat/types"
)

const (
	ApiURL     = "http://localhost:11434/api/generate"
	RefreshURL = "http://localhost:11434/api/tags"
)

type OllamaProvider struct {
	Url            string
	logger         *logger.Logger
	context        []int
	model          string
	ModelRefresher *types.ModelRefresher
}

// NewOllamaProvider creates a new OllamaProvider configured to use ApiURL.
// The provided logger is attached and the internal context is initialized as an empty slice.
func NewOllamaProvider(logger *logger.Logger,
	model string,
	mf *types.ModelRefresher,
) *OllamaProvider {
	return &OllamaProvider{
		Url:            ApiURL,
		logger:         logger,
		context:        []int{},
		model:          model,
		ModelRefresher: mf,
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

type OllamaResponse struct {
	Models []OllamaModel `json:"models"`
}

type OllamaModel struct {
	Name  string `json:"name"`
	Model string `json:"model"`
}

func (op *OllamaProvider) RetrieveModels() ([]types.Model, error) {
	// TODO: Add context timeout
	if !op.ModelRefresher.IsStale() {
		return op.ModelRefresher.RetrieveModels(), nil
	}

	req, err := http.NewRequest("GET", RefreshURL, nil)
	if err != nil {
		return nil, err
	}
	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response OllamaResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	modelList := []types.Model{}
	for _, v := range response.Models {
		modelList = append(modelList, types.Model{Name: v.Name})
	}

	if err := op.ModelRefresher.StashModels(modelList); err != nil {
		op.logger.Error("failed to stash models", "error", err)
		return modelList, nil
	}

	return modelList, nil
}

func (op *OllamaProvider) SetModel(model string) {
	op.model = model
}
