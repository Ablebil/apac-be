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

// var responseSchemaDefault = &genai.Schema{
// 	Type: "object",
// 	Properties: map[string]*genai.Schema{
// 		"id":          {Type: "string", Description: "Unique identifier for the trip"},
// 		"title":       {Type: "string", Description: "Title of the trip"},
// 		"destination": {Type: "string", Description: "Destination country/city"},
// 		"startDate":   {Type: "string", Format: "date-time", Description: "Start date of the trip"},
// 		"endDate":     {Type: "string", Format: "date-time", Description: "End date of the trip"},
// 		"duration":    {Type: "integer", Description: "Duration of the trip in days"},
// 		"travelers":   {Type: "integer", Description: "Number of travelers"},
// 		"budget":      {Type: "string", Description: "Estimated budget for the trip"},
// 		"summary":     {Type: "string", Description: "Brief summary of the trip"},
// 		"days": {
// 			Type: "array",
// 			Items: &genai.Schema{
// 				Type: "object",
// 				Properties: map[string]*genai.Schema{
// 					"day":         {Type: "integer", Description: "Day number of the trip"},
// 					"date":        {Type: "string", Format: "date-time", Description: "Date of the day"},
// 					"title":       {Type: "string", Description: "Title for the day's activities"},
// 					"description": {Type: "string", Description: "Description of the day's activities"},
// 					"activities": {
// 						Type: "array",
// 						Items: &genai.Schema{
// 							Type: "object",
// 							Properties: map[string]*genai.Schema{
// 								"time":        {Type: "string", Description: "Time of the activity"},
// 								"title":       {Type: "string", Description: "Title of the activity"},
// 								"description": {Type: "string", Description: "Description of the activity"},
// 								"location":    {Type: "string", Description: "Location of the activity"},
// 								"address":     {Type: "string", Description: "Address of the location"},
// 								"cost":        {Type: "string", Description: "Cost of the activity"},
// 								"tags": {
// 									Type: "array",
// 									Items: &genai.Schema{
// 										Type: "string",
// 									},
// 									Description: "Tags for the activity",
// 								},
// 							},
// 							Required: []string{"time", "title", "description"},
// 						},
// 					},
// 					"accommodation": {
// 						Type: "object",
// 						Properties: map[string]*genai.Schema{
// 							"name":     {Type: "string", Description: "Name of the accommodation"},
// 							"address":  {Type: "string", Description: "Address of the accommodation"},
// 							"checkIn":  {Type: "string", Description: "Check-in time"},
// 							"checkOut": {Type: "string", Description: "Check-out time"},
// 							"cost":     {Type: "string", Description: "Cost per night"},
// 						},
// 						Required: []string{"name", "address"},
// 					},
// 					"meals": {
// 						Type: "object",
// 						Properties: map[string]*genai.Schema{
// 							"breakfast": {
// 								Type: "object",
// 								Properties: map[string]*genai.Schema{
// 									"time":        {Type: "string", Description: "Time of the meal"},
// 									"title":       {Type: "string", Description: "Title/name of the meal"},
// 									"description": {Type: "string", Description: "Description of the meal"},
// 									"location":    {Type: "string", Description: "Location of the meal"},
// 									"address":     {Type: "string", Description: "Address of the location"},
// 									"cost":        {Type: "string", Description: "Cost of the meal"},
// 								},
// 							},
// 							"lunch": {
// 								Type: "object",
// 								Properties: map[string]*genai.Schema{
// 									"time":        {Type: "string", Description: "Time of the meal"},
// 									"title":       {Type: "string", Description: "Title/name of the meal"},
// 									"description": {Type: "string", Description: "Description of the meal"},
// 									"location":    {Type: "string", Description: "Location of the meal"},
// 									"address":     {Type: "string", Description: "Address of the location"},
// 									"cost":        {Type: "string", Description: "Cost of the meal"},
// 								},
// 							},
// 							"dinner": {
// 								Type: "object",
// 								Properties: map[string]*genai.Schema{
// 									"time":        {Type: "string", Description: "Time of the meal"},
// 									"title":       {Type: "string", Description: "Title/name of the meal"},
// 									"description": {Type: "string", Description: "Description of the meal"},
// 									"location":    {Type: "string", Description: "Location of the meal"},
// 									"address":     {Type: "string", Description: "Address of the location"},
// 									"cost":        {Type: "string", Description: "Cost of the meal"},
// 								},
// 							},
// 						},
// 					},
// 					"transportation": {
// 						Type: "object",
// 						Properties: map[string]*genai.Schema{
// 							"mode":          {Type: "string", Description: "Mode of transportation"},
// 							"details":       {Type: "string", Description: "Details about the transportation"},
// 							"departureTime": {Type: "string", Description: "Departure time"},
// 							"arrivalTime":   {Type: "string", Description: "Arrival time"},
// 							"cost":          {Type: "string", Description: "Cost of transportation"},
// 						},
// 					},
// 					"notes": {Type: "string", Description: "Additional notes for the day"},
// 				},
// 				Required: []string{"day", "date", "title", "description"},
// 			},
// 		},
// 		"totalCost": {Type: "string", Description: "Total cost of the trip"},
// 		"createdAt": {Type: "string", Format: "date-time", Description: "When the itinerary was created"},
// 		"updatedAt": {Type: "string", Format: "date-time", Description: "When the itinerary was last updated"},
// 	},
// 	Required: []string{"id", "title", "destination", "startDate", "endDate", "duration", "days"},
// }

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
