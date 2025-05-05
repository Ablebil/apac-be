package response

type Err struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Payload any    `json:"payload,omitempty"`
}

type Res struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Payload    any    `json:"payload,omitempty"`
}
