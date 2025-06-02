package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type Service struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Versions    []Version `json:"versions"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Version struct {
	VersionNumber string `json:"version_number"`
	Details       string `json:"details"`
}

func ParseService(data []byte) (*Service, error) {
	var svc Service
	if err := json.Unmarshal(data, &svc); err != nil {
		return nil, fmt.Errorf("failed to parse Service: %w", err)
	}
	return &svc, nil
}
