package domain

import (
	"context"
	"time"
)

type Staff struct {
	ID           string    `json:"id" gorm:"primaryKey;default:gen_random_uuid()"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	HospitalID   string    `json:"hospital_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type StaffRepository interface {
	Create(ctx context.Context, staff *Staff) error
	GetByUsername(ctx context.Context, username string) (*Staff, error)
}

type StaffUsecase interface {
	CreateStaff(ctx context.Context, username, password, hospitalCode string) error
	Login(ctx context.Context, username, password, hospitalCode string) (string, error) 
}
func (Staff) TableName() string {
	return "staff"
}