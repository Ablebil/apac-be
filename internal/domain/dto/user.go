package dto

import "mime/multipart"

type GetProfileResponse struct {
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	PhotoURL    string   `json:"photo_url"`
	Preferences []string `json:"preferences"`
}

type EditProfileRequest struct {
	Name            string                `form:"name"`
	CurrentPassword string                `form:"current_password" validate:"required_with=NewPassword,omitempty,min=6"`
	NewPassword     string                `form:"new_password" validate:"omitempty,min=6"`
	Photo           *multipart.FileHeader `form:"photo"`
}

type AddPreferenceRequest struct {
	Preferences []string `json:"preferences" validate:"required"`
}

type RemovePreferenceRequest struct {
	Preference string `json:"preference" validate:"required"`
}
