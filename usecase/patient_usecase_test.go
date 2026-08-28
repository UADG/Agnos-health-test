package usecase_test

import (
	"context"
	"testing"

	"agnos-assignment/domain"
	"agnos-assignment/usecase"
)

// ==========================================
// 1. สร้าง Mock สำหรับ Patient Repo และ HIS Client
// ==========================================

type mockPatientRepo struct {
	mockSearch func(ctx context.Context, hospitalID string, criteria domain.PatientSearchCriteria) ([]*domain.Patient, error)
	mockCreate func(ctx context.Context, patient *domain.Patient) error
}

func (m *mockPatientRepo) Search(ctx context.Context, hospitalID string, criteria domain.PatientSearchCriteria) ([]*domain.Patient, error) {
	return m.mockSearch(ctx, hospitalID, criteria)
}
func (m *mockPatientRepo) Create(ctx context.Context, patient *domain.Patient) error {
	return m.mockCreate(ctx, patient)
}

type mockHISClient struct {
	mockFetch func(ctx context.Context, id string) (*domain.Patient, error)
}

func (m *mockHISClient) FetchPatientByID(ctx context.Context, id string) (*domain.Patient, error) {
	return m.mockFetch(ctx, id)
}

// ==========================================
// 2. Test Cases สำหรับ SearchPatients
// ==========================================

// เคสที่ 1: ค้นหาใน Database เราเจอเลย (ไม่ต้องเรียก HIS API)
func TestSearchPatients_FoundInDB(t *testing.T) {
	mockRepo := &mockPatientRepo{
		mockSearch: func(ctx context.Context, hospitalID string, criteria domain.PatientSearchCriteria) ([]*domain.Patient, error) {
			// จำลองว่าเจอคนไข้ 1 คน
			return []*domain.Patient{{ID: "1", FirstNameTH: "สมชาย"}}, nil
		},
	}
	mockHIS := &mockHISClient{
		mockFetch: func(ctx context.Context, id string) (*domain.Patient, error) {
			t.Errorf("HIS API should not be called if patient is found in local DB")
			return nil, nil
		},
	}

	pu := usecase.NewPatientUsecase(mockRepo, mockHIS)
	patients, err := pu.SearchPatients(context.Background(), "HOSP-1", domain.PatientSearchCriteria{FirstName: "สมชาย"})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(patients) != 1 {
		t.Errorf("expected 1 patient, got %d", len(patients))
	}
}

// เคสที่ 2: ค้นหาใน DB ไม่เจอ และส่ง NationalID มา -> ต้องไปเรียก HIS API แล้วเซฟลง DB
func TestSearchPatients_NotFoundInDB_FetchFromHIS(t *testing.T) {
	mockRepo := &mockPatientRepo{
		mockSearch: func(ctx context.Context, hospitalID string, criteria domain.PatientSearchCriteria) ([]*domain.Patient, error) {
			return []*domain.Patient{}, nil // คืนค่า Array ว่าง (ไม่เจอใน DB)
		},
		mockCreate: func(ctx context.Context, patient *domain.Patient) error {
			return nil // จำลองว่าเซฟลง DB สำเร็จ
		},
	}
	
	mockHIS := &mockHISClient{
		mockFetch: func(ctx context.Context, id string) (*domain.Patient, error) {
			// จำลองว่า HIS API เจอคนไข้
			return &domain.Patient{ID: "2", NationalID: "1234567890123"}, nil
		},
	}

	pu := usecase.NewPatientUsecase(mockRepo, mockHIS)
	patients, err := pu.SearchPatients(context.Background(), "HOSP-1", domain.PatientSearchCriteria{NationalID: "1234567890123"})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(patients) != 1 {
		t.Errorf("expected 1 patient to be returned from HIS, got %d", len(patients))
	}
}

// เคสที่ 3: ค้นหาด้วยชื่อ (ไม่มี ID) ใน DB ไม่เจอ -> ต้องไม่เรียก HIS API (คืนค่าว่างเลย)
func TestSearchPatients_NotFoundInDB_NoIDProvided(t *testing.T) {
	mockRepo := &mockPatientRepo{
		mockSearch: func(ctx context.Context, hospitalID string, criteria domain.PatientSearchCriteria) ([]*domain.Patient, error) {
			return []*domain.Patient{}, nil 
		},
	}
	
	mockHIS := &mockHISClient{
		mockFetch: func(ctx context.Context, id string) (*domain.Patient, error) {
			t.Errorf("HIS API should not be called because no NationalID/PassportID was provided")
			return nil, nil
		},
	}

	pu := usecase.NewPatientUsecase(mockRepo, mockHIS)
	patients, err := pu.SearchPatients(context.Background(), "HOSP-1", domain.PatientSearchCriteria{FirstName: "สมหญิง"})

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(patients) != 0 {
		t.Errorf("expected 0 patients, got %d", len(patients))
	}
}