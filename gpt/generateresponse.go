package gpt

import (
	"context"
	"errors"

	"github.com/lushenle/mmchatgpt/config"
	gogpt "github.com/sashabaranov/go-gpt3"
)

func GenerateResponse(message string) (string, error) {
	apiKey := config.GetOpenAIAPIKey()
	ctx := context.Background()

	// Create the OpenAI API client with your API key.
	client := gogpt.NewClient(apiKey)

	// Configure the parameters for the GPT-3 completion task.
	completionReq := gogpt.CompletionRequest{
		Prompt:           message,
		MaxTokens:        300,
		Temperature:      0.7,
		Model:            gogpt.GPT3TextDavinci003,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
		TopP:             1,
	}

	// Submit the request to the OpenAI API and wait for the response.
	completionResp, err := client.CreateCompletion(ctx, completionReq)
	if err != nil {
		return "", err
	}

	// Extract the response text from the completion task response.
	if len(completionResp.Choices) > 0 {
		response := completionResp.Choices[0].Text
		return response, nil
	} else {
		return "", errors.New("received empty response from OpenAI API")
	}
}
