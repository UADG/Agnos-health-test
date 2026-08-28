package usecase_test

import (
	"context"
	"errors"
	"testing"

	"agnos-assignment/domain"
	"agnos-assignment/usecase"
	"golang.org/x/crypto/bcrypt"
)

type mockHospitalRepo struct {
	mockGetByCode func(ctx context.Context, code string) (*domain.Hospital, error)
}
func (m *mockHospitalRepo) GetByCode(ctx context.Context, code string) (*domain.Hospital, error) {
	return m.mockGetByCode(ctx, code)
}

type mockStaffRepo struct {
	mockCreate        func(ctx context.Context, staff *domain.Staff) error
	mockGetByUsername func(ctx context.Context, username string) (*domain.Staff, error)
}
func (m *mockStaffRepo) Create(ctx context.Context, staff *domain.Staff) error {
	return m.mockCreate(ctx, staff)
}
func (m *mockStaffRepo) GetByUsername(ctx context.Context, username string) (*domain.Staff, error) {
	return m.mockGetByUsername(ctx, username)
}


func TestCreateStaff_Success(t *testing.T) {
	mockHR := &mockHospitalRepo{
		mockGetByCode: func(ctx context.Context, code string) (*domain.Hospital, error) {
			return &domain.Hospital{ID: "hosp-123", Code: "HOSPITAL_A"}, nil
		},
	}

	mockSR := &mockStaffRepo{
		mockCreate: func(ctx context.Context, staff *domain.Staff) error {
			return nil 
		},
	}

	su := usecase.NewStaffUsecase(mockSR, mockHR, "secret")

	err := su.CreateStaff(context.Background(), "test_user", "password123", "HOSPITAL_A")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestCreateStaff_HospitalNotFound(t *testing.T) {
	mockHR := &mockHospitalRepo{
		mockGetByCode: func(ctx context.Context, code string) (*domain.Hospital, error) {
			return nil, errors.New("hospital not found")
		},
	}

	mockSR := &mockStaffRepo{} 

	su := usecase.NewStaffUsecase(mockSR, mockHR, "secret")

	err := su.CreateStaff(context.Background(), "test_user", "password123", "WRONG_HOSPITAL")

	if err == nil {
		t.Errorf("expected error 'hospital not found', got nil")
	} else if err.Error() != "hospital not found" {
		t.Errorf("expected error 'hospital not found', got %v", err)
	}
}


func TestLogin_Success(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	mockSR := &mockStaffRepo{
		mockGetByUsername: func(ctx context.Context, username string) (*domain.Staff, error) {
			return &domain.Staff{
				ID:           "staff-123",
				Username:     "admin_a",
				PasswordHash: string(hashedPassword),
				HospitalID:   "hosp-123",
			}, nil
		},
	}

	mockHR := &mockHospitalRepo{
		mockGetByCode: func(ctx context.Context, code string) (*domain.Hospital, error) {
			return &domain.Hospital{ID: "hosp-123", Code: "HOSPITAL_A"}, nil
		},
	}

	su := usecase.NewStaffUsecase(mockSR, mockHR, "my-secret-key")

	token, err := su.Login(context.Background(), "admin_a", "password123", "HOSPITAL_A")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if token == "" {
		t.Errorf("expected token string, got empty")
	}
}

func TestLogin_InvalidUsername(t *testing.T) {
	mockSR := &mockStaffRepo{
		mockGetByUsername: func(ctx context.Context, username string) (*domain.Staff, error) {
			return nil, errors.New("staff not found")
		},
	}
	mockHR := &mockHospitalRepo{}

	su := usecase.NewStaffUsecase(mockSR, mockHR, "secret")

	_, err := su.Login(context.Background(), "wrong_user", "password123", "HOSPITAL_A")

	if err == nil || err.Error() != "invalid credentials" {
		t.Errorf("expected error 'invalid credentials', got %v", err)
	}
}

func TestLogin_HospitalMismatch(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	mockSR := &mockStaffRepo{
		mockGetByUsername: func(ctx context.Context, username string) (*domain.Staff, error) {
			return &domain.Staff{
				HospitalID:   "hosp-123", 
				PasswordHash: string(hashedPassword),
			}, nil
		},
	}

	mockHR := &mockHospitalRepo{
		mockGetByCode: func(ctx context.Context, code string) (*domain.Hospital, error) {
			return &domain.Hospital{ID: "hosp-999", Code: "HOSPITAL_B"}, nil
		},
	}

	su := usecase.NewStaffUsecase(mockSR, mockHR, "secret")

	_, err := su.Login(context.Background(), "admin_a", "password123", "HOSPITAL_B")

	if err == nil || err.Error() != "invalid credentials" {
		t.Errorf("expected error 'invalid credentials' due to hospital mismatch, got %v", err)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	mockSR := &mockStaffRepo{
		mockGetByUsername: func(ctx context.Context, username string) (*domain.Staff, error) {
			return &domain.Staff{
				HospitalID:   "hosp-123",
				PasswordHash: string(hashedPassword),
			}, nil
		},
	}

	mockHR := &mockHospitalRepo{
		mockGetByCode: func(ctx context.Context, code string) (*domain.Hospital, error) {
			return &domain.Hospital{ID: "hosp-123", Code: "HOSPITAL_A"}, nil
		},
	}

	su := usecase.NewStaffUsecase(mockSR, mockHR, "secret")

	_, err := su.Login(context.Background(), "admin_a", "wrong_password", "HOSPITAL_A")

	if err == nil || err.Error() != "invalid credentials" {
		t.Errorf("expected error 'invalid credentials' due to wrong password, got %v", err)
	}
}