package dto

type GoogleCallbackRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state" validate:"required"`
	Error string `json:"error"`
}

type GoogleProfileResponse struct {
	Email      string
	Username   string
	Name       string
	IsVerified bool
}
