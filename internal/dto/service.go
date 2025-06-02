package dto

import "catalog-service/internal/models"

type ServiceDTO struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Versions    []models.Version `json:"versions"`
	CreatedAt   string           `json:"created_at"`
	UpdatedAt   string           `json:"updated_at"`
}

type ServiceListData struct {
	Count    int           `json:"count"`
	Services []*ServiceDTO `json:"services"`
	Next     *string       `json:"next"`
}