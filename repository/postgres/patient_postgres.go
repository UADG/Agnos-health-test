package postgres

import (
	"context"

	"agnos-assignment/domain"
	"gorm.io/gorm"
)

type patientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) domain.PatientRepository {
	return &patientRepository{db: db}
}

func (r *patientRepository) Search(ctx context.Context, hospitalID string, criteria domain.PatientSearchCriteria) ([]*domain.Patient, error) {
	var patients []*domain.Patient

	query := r.db.WithContext(ctx).Where("hospital_id = ?", hospitalID)

	if criteria.NationalID != "" {
		query = query.Where("national_id = ?", criteria.NationalID)
	}
	if criteria.PassportID != "" {
		query = query.Where("passport_id = ?", criteria.PassportID)
	}
	if criteria.FirstName != "" {
		query = query.Where("first_name_th = ? OR first_name_en = ?", criteria.FirstName, criteria.FirstName)
	}
	if criteria.LastName != "" {
		query = query.Where("last_name_th = ? OR last_name_en = ?", criteria.LastName, criteria.LastName)
	}
	if criteria.DateOfBirth != "" {
		query = query.Where("date_of_birth = ?", criteria.DateOfBirth)
	}
	if criteria.PhoneNumber != "" {
		query = query.Where("phone_number = ?", criteria.PhoneNumber)
	}
	if criteria.Email != "" {
		query = query.Where("email = ?", criteria.Email)
	}

	err := query.Find(&patients).Error
	if err != nil {
		return nil, err
	}

	return patients, nil
}

func (r *patientRepository) Create(ctx context.Context, patient *domain.Patient) error {
	return r.db.WithContext(ctx).Create(patient).Error
}