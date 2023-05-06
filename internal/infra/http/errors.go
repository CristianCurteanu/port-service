package http

type ApiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
