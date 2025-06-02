package usecase

import (
	"catalog-service/internal/constants"
	"catalog-service/internal/dto"
	"catalog-service/internal/repository"
	"context"
)

type ServiceUsecase interface {
	Search(ctx context.Context, query string, page, limit int) ([]*dto.ServiceDTO, int, error)
	FindByID(ctx context.Context, id string) (*dto.ServiceDTO, error)
}

type serviceUsecase struct {
	repo repository.ServiceRepository
}

func NewServiceUsecase(repo repository.ServiceRepository) ServiceUsecase {
	return &serviceUsecase{repo: repo}
}

func (u *serviceUsecase) Search(ctx context.Context, query string, page, limit int) ([]*dto.ServiceDTO, int, error) {
	services, total, err := u.repo.Search(ctx, query, page, limit)
	if err != nil {
		return nil, 0, err
	}
	dtos := make([]*dto.ServiceDTO, 0, len(services))
	for _, svc := range services {
		dtos = append(dtos, &dto.ServiceDTO{
			ID:          svc.ID,
			Name:        svc.Name,
			Description: svc.Description,
			Versions:    svc.Versions,
			CreatedAt:   svc.CreatedAt.Format(constants.Iso8601Format),
			UpdatedAt:   svc.UpdatedAt.Format(constants.Iso8601Format),
		})
	}
	return dtos, total, nil
}

func (u *serviceUsecase) FindByID(ctx context.Context, id string) (*dto.ServiceDTO, error) {
	svc, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.ServiceDTO{
		ID:          svc.ID,
		Name:        svc.Name,
		Description: svc.Description,
		Versions:    svc.Versions,
		CreatedAt:   svc.CreatedAt.Format(constants.Iso8601Format),
		UpdatedAt:   svc.UpdatedAt.Format(constants.Iso8601Format),
	}, nil
}
