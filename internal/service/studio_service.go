package service

import (
	"context"

	"github.com/google/uuid"

	"shutterspace/internal/domain"
	"shutterspace/internal/repository"
)

type StudioService interface {
	ListStudios(ctx context.Context, filter repository.StudioFilter) ([]domain.Studio, int64, error)
	GetStudioByID(ctx context.Context, id uuid.UUID) (*domain.Studio, error)
	GetStudioBySlug(ctx context.Context, slug string) (*domain.Studio, error)
	GetMyStudios(ctx context.Context, adminID uuid.UUID) ([]domain.Studio, error)
	GetAllStudioTypes(ctx context.Context) ([]domain.StudioType, error)
}

type studioServiceImpl struct {
	studioRepo     repository.StudioRepository
	studioTypeRepo repository.StudioTypeRepository
}

func NewStudioService(studioRepo repository.StudioRepository, typeRepo repository.StudioTypeRepository) StudioService {
	return &studioServiceImpl{
		studioRepo:     studioRepo,
		studioTypeRepo: typeRepo,
	}
}

func (s *studioServiceImpl) ListStudios(ctx context.Context, filter repository.StudioFilter) ([]domain.Studio, int64, error) {
	return s.studioRepo.FindAll(ctx, filter)
}

func (s *studioServiceImpl) GetStudioByID(ctx context.Context, id uuid.UUID) (*domain.Studio, error) {
	return s.studioRepo.FindByID(ctx, id)
}

func (s *studioServiceImpl) GetStudioBySlug(ctx context.Context, slug string) (*domain.Studio, error) {
	return s.studioRepo.FindBySlug(ctx, slug)
}

func (s *studioServiceImpl) GetMyStudios(ctx context.Context, adminID uuid.UUID) ([]domain.Studio, error) {
	return s.studioRepo.FindByManagedBy(ctx, adminID)
}

func (s *studioServiceImpl) GetAllStudioTypes(ctx context.Context) ([]domain.StudioType, error) {
	return s.studioTypeRepo.FindAll(ctx)
}
