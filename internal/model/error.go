package model

type ErrorMessage struct {
	Message   string `json:"message"`
	ErrorCode string `json:"error_code"`
}
