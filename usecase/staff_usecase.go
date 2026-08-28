package usecase

import (
	"context"
	"errors"
	"time"

	"agnos-assignment/domain"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type staffUsecase struct {
	staffRepo    domain.StaffRepository
	hospitalRepo domain.HospitalRepository
	jwtSecret    string
}

func NewStaffUsecase(sr domain.StaffRepository, hr domain.HospitalRepository, secret string) domain.StaffUsecase {
	return &staffUsecase{
		staffRepo:    sr,
		hospitalRepo: hr,
		jwtSecret:    secret,
	}
}

func (u *staffUsecase) CreateStaff(ctx context.Context, username, password, hospitalCode string) error {
	hospital, err := u.hospitalRepo.GetByCode(ctx, hospitalCode)
	if err != nil {
		return errors.New("hospital not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	staff := &domain.Staff{
		Username:     username,
		PasswordHash: string(hashedPassword),
		HospitalID:   hospital.ID,
	}

	return u.staffRepo.Create(ctx, staff)
}

func (u *staffUsecase) Login(ctx context.Context, username, password, hospitalCode string) (string, error) {
	staff, err := u.staffRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	hospital, err := u.hospitalRepo.GetByCode(ctx, hospitalCode)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if staff.HospitalID != hospital.ID {
		return "", errors.New("invalid credentials") 
	}

	err = bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"user_id":     staff.ID,
		"hospital_id": staff.HospitalID,
		"exp":         time.Now().Add(time.Hour * 24).Unix(), 
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.jwtSecret))
}