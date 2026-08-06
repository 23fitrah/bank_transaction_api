package model

type Response struct {
	Status       string            `json:"status"`
	ResponseCode string            `json:"response_code"`
	Message      string            `json:"message"`
	Errors       map[string]string `json:"errors,omitempty"`
	Data         interface{}       `json:"payload,omitempty"`
}
