package usecase

import (
	"catalog-service/internal/constants"
	"catalog-service/internal/dto"
	"catalog-service/internal/models"
	"catalog-service/internal/repository"
	"context"
)

type ServiceUsecase interface {
	Search(ctx context.Context, query string, page, limit int) ([]*dto.ServiceDTO, int, error)
	FindByID(ctx context.Context, id string) (*dto.ServiceDTO, error)
	Create(ctx context.Context, req *dto.ServiceDTO) (*dto.ServiceDTO, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, id string, req *dto.ServiceDTO) (*dto.ServiceDTO, error)
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

func (u *serviceUsecase) Create(ctx context.Context, req *dto.ServiceDTO) (*dto.ServiceDTO, error) {
	svc := &models.Service{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Versions:    req.Versions,
	}
	if err := u.repo.Create(ctx, svc); err != nil {
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

func (u *serviceUsecase) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}

func (u *serviceUsecase) Update(ctx context.Context, id string, req *dto.ServiceDTO) (*dto.ServiceDTO, error) {
	svc, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Description != "" {
		svc.Description = req.Description
	}
	if len(req.Versions) > 0 {
		svc.Versions = append(svc.Versions, req.Versions...)
	}
	if err := u.repo.Update(ctx, svc); err != nil {
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
