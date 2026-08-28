package postgres

import (
	"context"
	"errors"

	"agnos-assignment/domain"
	"gorm.io/gorm"
)

type hospitalRepository struct {
	db *gorm.DB
}

func NewHospitalRepository(db *gorm.DB) domain.HospitalRepository {
	return &hospitalRepository{db: db}
}

func (r *hospitalRepository) GetByCode(ctx context.Context, code string) (*domain.Hospital, error) {
	var hospital domain.Hospital

	err := r.db.WithContext(ctx).Where("code = ?", code).First(&hospital).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("hospital not found")
		}
		return nil, err
	}

	return &hospital, nil
}