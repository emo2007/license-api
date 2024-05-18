package rest

import (
	"errors"
	"net/http"
)

var (
	ErrorBadPathParams = errors.New("bad path params")
)

type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func buildApiError(code int, message string) apiError {
	return apiError{
		Code:    code,
		Message: message,
	}
}

func mapError(err error) apiError {
	switch {
	// server error
	case errors.Is(err, ErrorBadPathParams):
		return buildApiError(http.StatusBadRequest, "Invalid Path Params")

	default:
		return buildApiError(http.StatusInternalServerError, "Internal Server Error")
	}
}
