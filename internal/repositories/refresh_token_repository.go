package repositories

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/kraftzpepe/auth-service/internal/models"
)

type RefreshTokenRepository struct {
	DB *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{DB: db}
}

func (repo *RefreshTokenRepository) SaveRefreshToken(userID uuid.UUID, token string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := repo.DB.Exec(query, userID, token, expiresAt)
	return err
}

func (repo *RefreshTokenRepository) FindRefreshToken(token string) (*models.RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM refresh_tokens
		WHERE token = $1
	`
	row := repo.DB.QueryRow(query, token)

	var rt models.RefreshToken
	err := row.Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &rt, nil
}

func (repo *RefreshTokenRepository) UpdateRefreshToken(userID uuid.UUID, token string, expiresAt time.Time) error {
	query := `
		UPDATE refresh_tokens
		SET token = $1, expires_at = $2
		WHERE user_id = $3
	`
	_, err := repo.DB.Exec(query, token, expiresAt, userID)
	return err
}
