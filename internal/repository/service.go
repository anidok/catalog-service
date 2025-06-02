package repository

import (
	"context"
	"encoding/json"

	"catalog-service/internal/logger"
	"catalog-service/internal/models"
	"catalog-service/internal/opensearch"
)

const (
	ServiceIndexName = "services"
)

type ServiceRepositoryImpl struct {
	opensearch.Client
}

func NewServiceRepository(client opensearch.Client) (ServiceRepository, error) {
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

func (r *ServiceRepositoryImpl) Search(ctx context.Context, query string, page, limit int) ([]*models.Service, int, error) {
	log := logger.NewContextLogger(ctx, "ServiceRepositoryImpl/Search")
	from := (page - 1) * limit
	log.Debugf("searching for query='%s', page=%d, limit=%d, from=%d", query, page, limit, from)
	searchBody := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"name", "description"},
				"type":   "best_fields",
			},
		},
		"from": from,
		"size": limit,
	}

	hits, total, err := r.Client.Search(ctx, ServiceIndexName, searchBody)
	if err != nil {
		log.Errorf(err, "failed to execute search: %v", err)
		return nil, 0, err
	}

	services := make([]*models.Service, 0, len(hits))
	for _, hit := range hits {
		var svc models.Service
		b, _ := json.Marshal(hit)
		if err := json.Unmarshal(b, &svc); err == nil {
			services = append(services, &svc)
		} else {
			log.Errorf(err, "failed to unmarshal search hit: %v", err)
		}
	}
	log.Infof("search completed: found %d results", total)
	return services, total, nil
}
