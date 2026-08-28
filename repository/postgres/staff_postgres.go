package postgres

import (
	"context"
	"errors"

	"agnos-assignment/domain"
	"gorm.io/gorm"
)

type staffRepository struct {
	db *gorm.DB
}

func NewStaffRepository(db *gorm.DB) domain.StaffRepository {
	return &staffRepository{db: db}
}

func (r *staffRepository) Create(ctx context.Context, staff *domain.Staff) error {
	return r.db.WithContext(ctx).Create(staff).Error
}

func (r *staffRepository) GetByUsername(ctx context.Context, username string) (*domain.Staff, error) {
	var staff domain.Staff
	
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&staff).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("staff not found")
		}
		return nil, err
	}
	
	return &staff, nil
}