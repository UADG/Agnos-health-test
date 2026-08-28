package usecase

import (
	"context"

	"agnos-assignment/domain"
)

type patientUsecase struct {
	patientRepo domain.PatientRepository
	hisClient   domain.HISClient
}

func NewPatientUsecase(pr domain.PatientRepository, hc domain.HISClient) domain.PatientUsecase {
	return &patientUsecase{
		patientRepo: pr,
		hisClient:   hc,
	}
}

func (u *patientUsecase) SearchPatients(ctx context.Context, hospitalID string, criteria domain.PatientSearchCriteria) ([]*domain.Patient, error) {
	patients, err := u.patientRepo.Search(ctx, hospitalID, criteria)
	if err != nil {
		return nil, err
	}

	if len(patients) == 0 && (criteria.NationalID != "" || criteria.PassportID != "") {
		searchID := criteria.NationalID
		if criteria.PassportID != "" {
			searchID = criteria.PassportID
		}

		hisPatient, err := u.hisClient.FetchPatientByID(ctx, searchID)
		if err == nil && hisPatient != nil {
			hisPatient.HospitalID = hospitalID 
			
			_ = u.patientRepo.Create(ctx, hisPatient) 
			
			patients = append(patients, hisPatient)
		}
	}

	return patients, nil
}