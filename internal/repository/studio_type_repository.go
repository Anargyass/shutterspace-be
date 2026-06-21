package repository

import (
	"context"

	"gorm.io/gorm"

	"shutterspace/internal/domain"
)

// StudioTypeRepository adalah interface untuk mengakses studio_types
type StudioTypeRepository interface {
	FindAll(ctx context.Context) ([]domain.StudioType, error)
}

type studioTypeRepository struct {
	db *gorm.DB
}

func NewStudioTypeRepository(db *gorm.DB) StudioTypeRepository {
	return &studioTypeRepository{db: db}
}

func (r *studioTypeRepository) FindAll(ctx context.Context) ([]domain.StudioType, error) {
	var types []domain.StudioType
	err := r.db.WithContext(ctx).Order("id ASC").Find(&types).Error
	return types, err
}
