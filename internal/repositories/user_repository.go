package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kraftzpepe/auth-service/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// CreateUser inserts a new user into the database
func (repo *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, username, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := repo.DB.ExecContext(ctx, query, user.ID, user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	return err
}

// GetUserByEmail retrieves a user by their email
func (repo *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	row := repo.DB.QueryRowContext(ctx, query, email)

	user := &models.User{}
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No user found
		}
		return nil, err
	}

	return user, nil
}

// GetUserByUUID retrieves a user by their UUID
func (repo *UserRepository) GetUserByUUID(ctx context.Context, uuid string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	row := repo.DB.QueryRowContext(ctx, query, uuid)

	user := &models.User{}
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No user found
		}
		return nil, err
	}

	return user, nil
}

// GetUserByUsername retrieves a user by their username
func (repo *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at
		FROM users
		WHERE username = $1
	`
	row := repo.DB.QueryRowContext(ctx, query, username)

	user := &models.User{}
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No user found
		}
		return nil, err
	}

	return user, nil
}
