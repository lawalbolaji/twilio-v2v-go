package groq

import (
	"bytes"
	"encoding/json"
	"fmt"
	chttp "twilio-v2v/pkg/http"
)

type Groq struct {
	APIKey  string
	BaseUrl string
	/* add other configuration option here */
}

type CompletionPayload struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompletionResponse struct {
	Id      string   `json:"id"`
	Created int      `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

func (c *CompletionResponse) UnmarshalJSON(d []byte) error {
	type C CompletionResponse
	if err := json.Unmarshal(d, (*C)(c)); err != nil {
		return err
	}

	if len(c.Choices) == 0 {
		return fmt.Errorf("expected to find CompletionResponse Payload, found: '%s' instead", string(d))
	}
	return nil
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

func (g *Groq) GetCompletion(query string) (string, error) {
	j, err := json.Marshal(g.createPayload(query))
	if err != nil {
		return "", err
	}

	c := &chttp.Client{}
	path := "/chat/completions"
	_, res, err := c.Post(g.BaseUrl+path, bytes.NewBuffer(j), map[string]string{"Authorization": fmt.Sprintf("Bearer %s", g.APIKey)})
	if err != nil {
		return "", err
	}

	var cr CompletionResponse
	err = json.Unmarshal(res, &cr)
	if err != nil {
		return "", err
	}

	return cr.Choices[0].Message.Content, nil
}

func (g *Groq) createPayload(query string) *CompletionPayload {
	return &CompletionPayload{
		Messages: []Message{
			{
				Role: "system",
				Content: `you are a helpful assistant that helps with
						  general knowledge questions. you should
						  prioritize the shortest answers that convey
						  the point. your answers must never exceed
						  30 words under any circumstances.`,
			},
			{
				Role:    "user",
				Content: query,
			},
		},
		Model: "llama-3.1-8b-instant",
	}
}
