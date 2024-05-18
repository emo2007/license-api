package domain

import (
	"encoding/json"
	"fmt"
)

// Generic

type Collection[T any] struct {
	Items      []T        `json:"items,omitempty"`
	Pagination Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	NextCursor string `json:"next_cursor,omitempty"`
	TotalItems uint32 `json:"total_items,omitempty"`
}

func BuildRequest[T any](data []byte) (*T, error) {
	var req T

	if err := json.Unmarshal(data, &req); err != nil {
		return nil, fmt.Errorf("error unmarshal request. %w", err)
	}

	return &req, nil
}
