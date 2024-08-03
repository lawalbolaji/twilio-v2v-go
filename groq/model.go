package groq

import (
	"errors"
	"log"
)

type Groq struct {
	client RestClient
}

func (groq *Groq) GetCompletion(query string) (string, error) {
	payload := &CompletionPayload{
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

	cResp, err := groq.client.GetCompletion(payload)
	if err != nil {
		return "", errors.New(err.Error())
	}

	return cResp.Choices[0].Message.Content, nil
}

func NewGroqClient(apiKey string) *Groq {
	if apiKey == "" {
		log.Fatal("Missing API key for Groq")
	}
	return &Groq{
		client: *newRestClient(apiKey),
	}
}
