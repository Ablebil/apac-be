package gemini

import (
	"apac/internal/domain/env"
	"context"
	"encoding/json"
	"os"

	"google.golang.org/genai"
)

type GeminiItf interface {
	Prompt([]string, string) (map[string]interface{}, error)
}

type Gemini struct {
	client *genai.Client
	config *genai.GenerateContentConfig
	model  string
}

func NewGemini(env *env.Env) (GeminiItf, error) {
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  env.GeminiAPIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	responseSchema, err := GetResponseSchema()
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

func GetResponseSchema() (*genai.Schema, error) {
	content, err := os.ReadFile("./resource/schema.json")
	if err != nil {
		return nil, err
	}

	var schema genai.Schema
	if err := schema.UnmarshalJSON(content); err != nil {
		return nil, err
	}

	return &schema, nil
}

func (g *Gemini) Prompt(preferences []string, prompt string) (map[string]interface{}, error) {
	var prefPrompt string
	if preferences != nil {
		prefPrompt += "FOLLOW PREFERENCES: ("
		for _, pref := range preferences {
			prefPrompt += pref + ", "
		}
		prefPrompt += ")\n\n"
	}

	result, err := g.client.Models.GenerateContent(
		context.Background(),
		g.model,
		genai.Text(prefPrompt+"NO NULL VALUES, NO N/A VALUES\n\n"+"PROMPT: "+prompt),
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
