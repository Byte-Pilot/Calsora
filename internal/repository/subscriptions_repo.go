package repository

import (
	"Calsora/internal/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type SubscriptionsRepository interface {
	Create(userID int) (*models.Subscriptions, error)
	GetSubDataByID(userID int) (*models.Subscriptions, error)
}

type subscriptionsRepository struct {
	db *pgxpool.Pool
}

func NewSubscriptionsRepository(db *pgxpool.Pool) *subscriptionsRepository {
	return &subscriptionsRepository{db: db}
}

func (r *subscriptionsRepository) Create(userID int) (*models.Subscriptions, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var userSub models.Subscriptions
	query := `INSERT INTO subscriptions (user_id) VALUES ($1) RETURNING user_id, plan, expires_at`
	err := r.db.QueryRow(ctx, query, userID).Scan(&userSub.UserID, &userSub.Plan, &userSub.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &userSub, nil
}

func (r *subscriptionsRepository) GetSubDataByID(userID int) (*models.Subscriptions, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var userSub models.Subscriptions
	query := `SELECT id, user_id, plan, status, started_at, expires_at FROM subscriptions WHERE user_id = $1`
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&userSub.ID, &userSub.UserID, &userSub.Plan, &userSub.Status, &userSub.StartedAT, &userSub.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &userSub, nil
}
