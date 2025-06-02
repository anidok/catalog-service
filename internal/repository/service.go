package repository

import (
	"context"

	"catalog-service/internal/logger"
	"catalog-service/internal/models"
	"catalog-service/internal/opensearch"
)

const (
	ServiceIndexName = "services"
)

type ServiceRepositoryImpl struct {
	*opensearch.Client
}

func NewServiceRepository(client *opensearch.Client) (ServiceRepository, error) {
	return &ServiceRepositoryImpl{
		Client: client,
	}, nil
}

func (r *ServiceRepositoryImpl) Create(ctx context.Context, service *models.Service) error {
	log := logger.NewContextLogger(ctx, "ServiceRepositoryImpl/Create")
	id := service.ID
	log.Debug("inserting record in services index")
	return r.IndexDocument(ctx, id, service, ServiceIndexName)
}

