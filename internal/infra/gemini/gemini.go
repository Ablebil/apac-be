package gemini

import (
	"apac/internal/domain/env"
	"context"
	"encoding/json"

	"google.golang.org/genai"
)

type GeminiItf interface {
	Prompt(string) (map[string]interface{}, error)
}

type Gemini struct {
	client *genai.Client
	config *genai.GenerateContentConfig
	model  string
}

var responseSchema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"trip_name": {
			Type: genai.TypeString,
		},
		"days": {
			Type: genai.TypeArray,
			Items: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"day": {
						Type: genai.TypeInteger,
					},
					"location": {
						Type: genai.TypeString,
					},
					"activities": {
						Type: genai.TypeArray,
						Items: &genai.Schema{
							Type: genai.TypeString,
						},
					},
					"meals": {
						Type: genai.TypeArray,
						Items: &genai.Schema{
							Type: genai.TypeString,
						},
					},
					"notes": {
						Type: genai.TypeString,
					},
				},
				Required: []string{
					"day",
					"location",
					"activities",
					"meals",
				},
			},
		},
		"notes": {
			Type: genai.TypeString,
		},
	},
	Required: []string{
		"trip_name",
		"days",
	},
	PropertyOrdering: []string{
		"trip_name",
		"days",
		"notes",
	},
}

func NewGemini(env *env.Env) (GeminiItf, error) {
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  env.GeminiAPIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema:   responseSchema,
	}

	return &Gemini{
		client: client,
		config: config,
		model:  env.GeminiModel,
	}, nil
}

func (g *Gemini) Prompt(prompt string) (map[string]interface{}, error) {
	result, err := g.client.Models.GenerateContent(
		context.Background(),
		g.model,
		genai.Text(prompt),
		g.config,
	)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	err = json.Unmarshal([]byte(result.Text()), &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
