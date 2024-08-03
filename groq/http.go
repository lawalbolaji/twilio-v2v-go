package groq

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type CompletionPayload struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
}

type CompletionResponse struct {
	Id      string   `json:"id"`
	Created int      `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RestClient struct {
	apiKey string
}

func (client *RestClient) GetCompletion(payload *CompletionPayload) (CompletionResponse, error) {
	payloadAsJson, _ := json.Marshal(payload)
	payloadReader := bytes.NewBuffer(payloadAsJson)
	req, _ := http.NewRequest(http.MethodPost, "https://api.groq.com/openai/v1/chat/completions", payloadReader)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.apiKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return CompletionResponse{}, errors.New(err.Error())
	}

	var cResp CompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&cResp); err != nil {
		return CompletionResponse{}, errors.New(err.Error())
	}

	return cResp, nil
}

func newRestClient(apiKey string) *RestClient {
	return &RestClient{
		apiKey: apiKey,
	}
}
