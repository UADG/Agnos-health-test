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
	// 1. ตรวจสอบว่ามีโรงพยาบาลนี้อยู่จริงหรือไม่
	hospital, err := u.hospitalRepo.GetByCode(ctx, hospitalCode)
	if err != nil {
		return errors.New("hospital not found")
	}

	// 2. Hash Password เพื่อความปลอดภัย
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 3. สร้าง Object Staff และสั่งบันทึกลงฐานข้อมูล
	staff := &domain.Staff{
		Username:     username,
		PasswordHash: string(hashedPassword),
		HospitalID:   hospital.ID,
	}

	return u.staffRepo.Create(ctx, staff)
}

func (u *staffUsecase) Login(ctx context.Context, username, password, hospitalCode string) (string, error) {
	// 1. ดึงข้อมูลพนักงานจาก Username
	staff, err := u.staffRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// 2. ดึงข้อมูลโรงพยาบาลที่พยายามล็อกอิน
	hospital, err := u.hospitalRepo.GetByCode(ctx, hospitalCode)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// 3. ตรวจสอบว่าพนักงานคนนี้สังกัดโรงพยาบาลนี้จริงๆ (ป้องกันการ Cross-login)
	if staff.HospitalID != hospital.ID {
		return "", errors.New("invalid credentials") // ใช้ข้อความกว้างๆ เพื่อไม่ให้แฮกเกอร์รู้ข้อมูล
	}

	// 4. ตรวจสอบ Password
	err = bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// 5. ล็อกอินสำเร็จ สร้าง JWT Token
	claims := jwt.MapClaims{
		"user_id":     staff.ID,
		"hospital_id": staff.HospitalID,
		"exp":         time.Now().Add(time.Hour * 24).Unix(), // Token หมดอายุใน 24 ชั่วโมง
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.jwtSecret))
}