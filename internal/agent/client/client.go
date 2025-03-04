package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/horhhe/disperser-calculator/internal/models"
)

type AgentClient struct {
	baseURL string
}

func NewAgentClient(baseURL string) *AgentClient {
	return &AgentClient{baseURL: baseURL}
}

func (ac *AgentClient) GetTask() (models.Task, error) {
	url := fmt.Sprintf("%s/internal/task", ac.baseURL)
	resp, err := http.Get(url)
	if err != nil {
		return models.Task{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return models.Task{}, errors.New("no tasks available")
	}
	if resp.StatusCode != http.StatusOK {
		return models.Task{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var body struct {
		Task models.Task `json:"task"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return models.Task{}, err
	}
	return body.Task, nil
}

func (ac *AgentClient) PostTaskResult(result models.TaskResultRequest) error {
	url := fmt.Sprintf("%s/internal/task", ac.baseURL)
	payload, err := json.Marshal(result)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
