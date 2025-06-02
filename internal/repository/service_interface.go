package repository

import (
	"catalog-service/internal/models"
	"context"
)

type ServiceRepository interface {
	Create(ctx context.Context, service *models.Service) error
	Search(ctx context.Context, query string, page, limit int) ([]*models.Service, int, error)
	FindByID(ctx context.Context, id string) (*models.Service, error)
	Delete(ctx context.Context, id string) error
}
