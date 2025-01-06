package repositories

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kraftzpepe/auth-service/internal/models"
)

type PasswordResetToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type PasswordResetTokenRepository struct {
	DB *sql.DB
}

func NewPasswordResetTokenRepository(db *sql.DB) *PasswordResetTokenRepository {
	return &PasswordResetTokenRepository{DB: db}
}

// CreateToken creates a new password reset token for a user
func (r *PasswordResetTokenRepository) CreateToken(userID uuid.UUID, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO password_reset_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.DB.Exec(query, userID, token, expiresAt)
	return err
}

func (repo *PasswordResetTokenRepository) SaveToken(userID uuid.UUID, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO password_reset_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := repo.DB.Exec(query, userID, token, expiresAt)
	return err
}

// FindToken retrieves a password reset token from the database
func (repo *PasswordResetTokenRepository) FindToken(token string) (*models.PasswordResetToken, error) {
	query := `
		SELECT user_id, token, expires_at
		FROM password_reset_tokens
		WHERE token = $1
	`
	row := repo.DB.QueryRow(query, token)

	var resetToken models.PasswordResetToken
	if err := row.Scan(&resetToken.UserID, &resetToken.Token, &resetToken.ExpiresAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No token found
		}
		return nil, err
	}

	return &resetToken, nil
}

// DeleteToken removes a password reset token from the database
func (r *PasswordResetTokenRepository) DeleteToken(token string) error {
	query := `DELETE FROM password_reset_tokens WHERE token = $1`
	_, err := r.DB.Exec(query, token)
	return err
}
