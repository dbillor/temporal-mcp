package model

import (
    "encoding/json"
    "net/http"
)

// ErrorResponse is a generic MCP error envelope.
type ErrorResponse struct {
    Error ErrorDetail `json:"error"`
}

// ErrorDetail describes the error.
type ErrorDetail struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

// NewError creates an ErrorResponse with the given code and message.
func NewError(code, msg string) *ErrorResponse {
    return &ErrorResponse{Error: ErrorDetail{Code: code, Message: msg}}
}

// Write writes the error to the http.ResponseWriter.
func (e *ErrorResponse) Write(w http.ResponseWriter, status int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(e)
}
