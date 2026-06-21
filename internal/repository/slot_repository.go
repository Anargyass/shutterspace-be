package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"shutterspace/internal/domain"
)

type slotRepository struct {
	db *gorm.DB
}

func NewSlotRepository(db *gorm.DB) SlotRepository {
	return &slotRepository{db: db}
}

func (r *slotRepository) FindByStudioAndDay(ctx context.Context, studioID uuid.UUID, day string) (*domain.AvailabilitySlot, error) {
	var slot domain.AvailabilitySlot
	err := r.db.WithContext(ctx).
		Where("studio_id = ? AND day_of_week = ?", studioID, day).
		First(&slot).Error
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return &slot, nil
}

func (r *slotRepository) FindAllByStudio(ctx context.Context, studioID uuid.UUID) ([]domain.AvailabilitySlot, error) {
	var slots []domain.AvailabilitySlot
	err := r.db.WithContext(ctx).
		Where("studio_id = ?", studioID).
		Order("CASE day_of_week WHEN 'monday' THEN 1 WHEN 'tuesday' THEN 2 WHEN 'wednesday' THEN 3 WHEN 'thursday' THEN 4 WHEN 'friday' THEN 5 WHEN 'saturday' THEN 6 WHEN 'sunday' THEN 7 END").
		Find(&slots).Error
	return slots, err
}

func (r *slotRepository) UpsertSlots(ctx context.Context, slots []domain.AvailabilitySlot) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "studio_id"}, {Name: "day_of_week"}},
			DoUpdates: clause.AssignmentColumns([]string{"open_time", "close_time", "is_open"}),
		}).
		Create(&slots).Error
}
