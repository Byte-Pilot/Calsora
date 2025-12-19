package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type AuthRepositoryInterface interface {
	SaveRefreshToken(userID int, refresh string, exp time.Time) error
	GetRefreshToken(token string) (int, time.Time, error)
	DeleteRefreshToken(token string) error
}

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) SaveRefreshToken(userID int, refresh string, exp time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)"
	_, err := r.db.Exec(ctx, query, userID, refresh, exp)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) GetRefreshToken(token string) (userID int, expiresAt time.Time, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT user_id, expires_at FROM refresh_tokens WHERE token = $1"
	err = r.db.QueryRow(ctx, query, token).Scan(&userID, &expiresAt)
	if err != nil {
		return 0, time.Now(), err
	}
	return userID, expiresAt, nil
}

func (r *AuthRepository) DeleteRefreshToken(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "DELETE FROM refresh_tokens WHERE token = $1"
	_, err := r.db.Exec(ctx, query, token)
	if err != nil {
		return err
	}
	return nil
}
