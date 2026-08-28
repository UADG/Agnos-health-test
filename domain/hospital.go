package domain

import (
	"time"
	"context"
)

type Hospital struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"` 
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type HospitalRepository interface {
	GetByCode(ctx context.Context,code string) (*Hospital, error)
}