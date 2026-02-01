package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type AuthRepository interface {
	SaveRefreshToken(userID int, refresh string, exp time.Time) error
	GetRefreshToken(token string) (int, time.Time, error)
	DeleteRefreshToken(token string) error
	DeleteAllRefreshTokens(userID int) error
}

type authRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *authRepository {
	return &authRepository{db: db}
}

func (r *authRepository) SaveRefreshToken(userID int, refresh string, exp time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, userID, refresh, exp)
	return err
}

func (r *authRepository) GetRefreshToken(token string) (userID int, expiresAt time.Time, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT user_id, expires_at FROM refresh_tokens WHERE token = $1`
	err = r.db.QueryRow(ctx, query, token).Scan(&userID, &expiresAt)
	if err != nil {
		return 0, time.Now(), err
	}
	return userID, expiresAt, nil
}

func (r *authRepository) DeleteRefreshToken(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM refresh_tokens WHERE token = $1`
	_, err := r.db.Exec(ctx, query, token)
	return err
}

func (r *authRepository) DeleteAllRefreshTokens(userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}
