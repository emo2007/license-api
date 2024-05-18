package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func CreateRequest[T any](r *http.Request) (*T, error) {
	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("error read request body. %w", err)
	}

	var request T

	if err := json.Unmarshal(data, &request); err != nil {
		return nil, fmt.Errorf("error unmarshal join request. %w", err)
	}

	return &request, nil
}
