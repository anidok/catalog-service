package repository

import (
	"catalog-service/internal/models"
	"context"
)

type ServiceRepository interface {
	Create(ctx context.Context, service *models.Service) error
}
