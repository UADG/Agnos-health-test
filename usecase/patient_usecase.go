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
	// 1. ค้นหาผู้ป่วยในฐานข้อมูลของเรา (Postgres) ก่อน
	patients, err := u.patientRepo.Search(ctx, hospitalID, criteria)
	if err != nil {
		return nil, err
	}

	// 2. Middleware Logic: หากระบุเลขบัตร (NationalID หรือ PassportID) แล้วไม่เจอในฐานข้อมูลของเรา
	// เราอาจจะต้องไปดึงจากระบบ HIS ของโรงพยาบาล A
	if len(patients) == 0 && (criteria.NationalID != "" || criteria.PassportID != "") {
		searchID := criteria.NationalID
		if criteria.PassportID != "" {
			searchID = criteria.PassportID
		}

		// ยิง API ภายนอกไปดึงข้อมูล
		hisPatient, err := u.hisClient.FetchPatientByID(ctx, searchID)
		if err == nil && hisPatient != nil {
			// บังคับว่าผู้ป่วยที่ดึงมาต้องอยู่ใน Hospital ID ของ Staff ที่ล็อกอินอยู่เท่านั้น
			hisPatient.HospitalID = hospitalID 
			
			// เซฟข้อมูลที่เพิ่งได้มาลงฐานข้อมูลของเรา (เสมือนการทำ Sync/Cache)
			_ = u.patientRepo.Create(ctx, hisPatient) 
			
			// นำข้อมูลที่ดึงมาได้ใส่กลับไปเป็นผลลัพธ์
			patients = append(patients, hisPatient)
		}
	}

	return patients, nil
}