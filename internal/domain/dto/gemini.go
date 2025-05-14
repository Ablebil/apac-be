package dto

type GeminiRequest struct {
	UsePreference bool   `json:"use_preference" default:"true"`
	Text          string `json:"text"`
}
