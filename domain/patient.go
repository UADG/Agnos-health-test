package domain

import (
	"context"
	"time"
)

type Patient struct {
	ID           string    `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	HospitalID   string    `json:"hospital_id"`
	PatientHN    string    `json:"patient_hn"`
	NationalID   string    `json:"national_id"`
	PassportID   string    `json:"passport_id"`
	FirstNameTH  string    `json:"first_name_th"`
	MiddleNameTH string    `json:"middle_name_th"`
	LastNameTH   string    `json:"last_name_th"`
	FirstNameEN  string    `json:"first_name_en"`
	MiddleNameEN string    `json:"middle_name_en"`
	LastNameEN   string    `json:"last_name_en"`
	DateOfBirth  string    `json:"date_of_birth"`
	PhoneNumber  string    `json:"phone_number"`
	Email        string    `json:"email"`
	Gender       string    `json:"gender"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PatientSearchCriteria struct {
	NationalID  string
	PassportID  string
	FirstName   string
	LastName    string
	DateOfBirth string
	PhoneNumber string
	Email       string
}

type PatientRepository interface {
	Search(ctx context.Context, hospitalID string, criteria PatientSearchCriteria) ([]*Patient, error)
	Create(ctx context.Context, patient *Patient) error
}

type HISClient interface {
	FetchPatientByID(ctx context.Context, id string) (*Patient, error) 
}

type PatientUsecase interface {
	SearchPatients(ctx context.Context, hospitalID string, criteria PatientSearchCriteria) ([]*Patient, error)
}