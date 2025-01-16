package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mohamedfawas/rmshop-admin-service/internal/domain"
)

type adminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) domain.AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) GetUserDetails(ctx context.Context, userID string) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, email, created_at 
         FROM users 
         WHERE id = $1`,
		userID).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error getting user details: %v", err)
	}
	return &user, nil
}
