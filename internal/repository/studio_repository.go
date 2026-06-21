package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"shutterspace/internal/domain"
)

type studioRepository struct {
	db *gorm.DB
}

func NewStudioRepository(db *gorm.DB) StudioRepository {
	return &studioRepository{db: db}
}

func (r *studioRepository) FindAll(ctx context.Context, filter StudioFilter) ([]domain.Studio, int64, error) {
	var studios []domain.Studio
	var total int64

	query := r.db.WithContext(ctx).
		Preload("StudioType").
		Where("is_active = true")

	if filter.TypeSlug != "" {
		query = query.Joins("JOIN studio_types st ON st.id = studios.studio_type_id").
			Where("st.slug = ?", filter.TypeSlug)
	}
	if filter.Area != "" {
		query = query.Where("area = ?", filter.Area)
	}
	if filter.MinPrice > 0 {
		query = query.Where("price_per_hour >= ?", filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		query = query.Where("price_per_hour <= ?", filter.MaxPrice)
	}

	query.Model(&domain.Studio{}).Count(&total)

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 50 {
		filter.Limit = 10
	}

	err := query.
		Order("rating DESC, review_count DESC").
		Offset((filter.Page - 1) * filter.Limit).
		Limit(filter.Limit).
		Find(&studios).Error

	return studios, total, err
}

func (r *studioRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Studio, error) {
	var studio domain.Studio
	err := r.db.WithContext(ctx).
		Preload("StudioType").
		Preload("AvailabilitySlots").
		First(&studio, "studios.id = ?", id).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return &studio, nil
}

func (r *studioRepository) FindBySlug(ctx context.Context, slug string) (*domain.Studio, error) {
	var studio domain.Studio
	err := r.db.WithContext(ctx).
		Preload("StudioType").
		Preload("AvailabilitySlots").
		Where("slug = ?", slug).
		First(&studio).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return &studio, nil
}

func (r *studioRepository) Create(ctx context.Context, studio *domain.Studio) error {
	return r.db.WithContext(ctx).Create(studio).Error
}

func (r *studioRepository) Update(ctx context.Context, studio *domain.Studio) error {
	return r.db.WithContext(ctx).Save(studio).Error
}

func (r *studioRepository) FindByManagedBy(ctx context.Context, userID uuid.UUID) ([]domain.Studio, error) {
	var studios []domain.Studio
	err := r.db.WithContext(ctx).
		Preload("StudioType").
		Where("managed_by = ?", userID).
		Find(&studios).Error
	return studios, err
}
