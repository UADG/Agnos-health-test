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

// NewHospitalRepository คืนค่าเป็น Interface domain.HospitalRepository
func NewHospitalRepository(db *gorm.DB) domain.HospitalRepository {
	return &hospitalRepository{db: db}
}

// GetByCode ดึงข้อมูลโรงพยาบาลจากตาราง hospitals โดยใช้คอลัมน์ code
func (r *hospitalRepository) GetByCode(ctx context.Context, code string) (*domain.Hospital, error) {
	var hospital domain.Hospital

	// ค้นหา record เดียวที่ตรงกับ code ที่ส่งมา
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&hospital).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("hospital not found")
		}
		return nil, err
	}

	return &hospital, nil
}