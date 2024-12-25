package repositories

import (
	"context"
	"database/sql"

	"github.com/kraftzpepe/auth-service/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (repo *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, username, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := repo.DB.ExecContext(ctx, query, user.ID, user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	return err
}
