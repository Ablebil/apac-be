package dto

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required,len=6,numeric"`
}

type ChoosePreferenceRequest struct {
	Email       string   `json:"email" validate:"required,email"`
	Preferences []string `json:"preferences" validate:"required"`
}

type LoginRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=6"`
	RememberMe bool   `json:"remember_me"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type GoogleCallbackRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state" validate:"required"`
	Error string `json:"error"`
}

type GoogleProfileResponse struct {
	ID       string `json:"google_id" validate:"required"`
	Email    string `json:"email" validate:"required, email"`
	Username string `json:"username" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Verified bool   `json:"verified" validate:"required"`
}
