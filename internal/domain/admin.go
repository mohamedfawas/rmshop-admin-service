package domain

import (
	"context"
)

type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt string
}

type AdminRepository interface {
	GetUserDetails(ctx context.Context, userID string) (*User, error)
}

type AdminService interface {
	GetUserDetails(ctx context.Context, userID string) (*User, error)
}
