package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"catalog-service/internal/logger"
	"catalog-service/internal/models"
	"catalog-service/internal/opensearch"

	"github.com/google/uuid"
)

const (
	ServiceIndexName   = "services"
	UpdatedAtSortField = "updated_at"
)

type ServiceRepositoryImpl struct {
	opensearch.Client
}

func NewServiceRepository(client opensearch.Client) (ServiceRepository, error) {
	return &ServiceRepositoryImpl{Client: client}, nil
}

func (r *ServiceRepositoryImpl) Create(ctx context.Context, service *models.Service) error {
	log := logger.NewContextLogger(ctx, "ServiceRepositoryImpl/Create")
	if err := r.prepareService(service); err != nil {
		log.Errorf(err, "failed to prepare service")
		return fmt.Errorf("failed to prepare service: %w", err)
	}

	log.Debug("inserting record in services index")
	return r.IndexDocument(ctx, service.ID, service, ServiceIndexName)
}

func (r *ServiceRepositoryImpl) Search(ctx context.Context, query string, page, limit int) ([]*models.Service, int, error) {
	log := logger.NewContextLogger(ctx, "ServiceRepositoryImpl/Search")
	from := (page - 1) * limit
	log.Debugf("searching for query='%s', page=%d, limit=%d, from=%d", query, page, limit, from)

	searchBody := buildSearchBody(query, from, limit)

	hits, total, err := r.Client.Search(ctx, ServiceIndexName, searchBody)
	if err != nil {
		log.Errorf(err, "failed to execute search")
		return nil, 0, fmt.Errorf("search query failed: %w", err)
	}

	services := make([]*models.Service, 0, len(hits))
	for _, hit := range hits {
		var svc models.Service
		bytes, _ := json.Marshal(hit)
		if err := json.Unmarshal(bytes, &svc); err != nil {
			log.Errorf(err, "failed to unmarshal search hit")
			continue
		}
		services = append(services, &svc)
	}

	log.Infof("search completed: found %d results", total)
	return services, total, nil
}

func (r *ServiceRepositoryImpl) prepareService(service *models.Service) error {
	if service == nil {
		return fmt.Errorf("service cannot be nil")
	}

	if service.ID == "" {
		service.ID = uuid.New().String()
	}

	now := time.Now().UTC()
	if service.CreatedAt.IsZero() {
		service.CreatedAt = now
	}
	service.UpdatedAt = now

	return nil
}

func buildSearchBody(query string, from, size int) map[string]interface{} {
	sortClause := []map[string]interface{}{
		{UpdatedAtSortField: map[string]interface{}{"order": "desc"}},
	}

	if query == "" {
		return map[string]interface{}{
			"query": map[string]interface{}{
				"match_all": map[string]interface{}{},
			},
			"from": from,
			"size": size,
			"sort": sortClause,
		}
	}

	return map[string]interface{}{
		"query": map[string]interface{}{
			"simple_query_string": map[string]interface{}{
				"query":            fmt.Sprintf("\"%s\"", query),
				"fields":           []string{"name", "description"},
				"default_operator": "and",
			},
		},
		"from": from,
		"size": size,
		"sort": sortClause,
	}
}
