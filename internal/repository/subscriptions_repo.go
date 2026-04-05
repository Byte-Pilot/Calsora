package repository

import (
	"Calsora/internal/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type SubscriptionsRepository interface {
	Create(userID int) (*models.Subscriptions, error)
	CreatePremium(userID int, plan string, exp time.Time) error
	GetSubDataByID(userID int) (*models.Subscriptions, error)
	MarkExpired(subID int) error
}

type subscriptionsRepository struct {
	db *pgxpool.Pool
}

func NewSubscriptionsRepository(db *pgxpool.Pool) *subscriptionsRepository {
	return &subscriptionsRepository{db: db}
}

func (r *subscriptionsRepository) Create(userID int) (*models.Subscriptions, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var userSub models.Subscriptions
	query := `INSERT INTO subscriptions (user_id) VALUES ($1) RETURNING user_id, plan, expires_at`
	err := r.db.QueryRow(ctx, query, userID).Scan(&userSub.UserID, &userSub.Plan, &userSub.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &userSub, nil
}

func (r *subscriptionsRepository) CreatePremium(userID int, plan string, exp time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO  subscriptions (user_id, plan, expires_at) VALUES ($1, $2, $3) `
	_, err := r.db.Exec(ctx, query, userID, plan, exp)
	return err
}

func (r *subscriptionsRepository) GetSubDataByID(userID int) (*models.Subscriptions, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var userSub models.Subscriptions
	query := `SELECT id, user_id, plan, status, started_at, expires_at 
			  FROM subscriptions WHERE user_id = $1 AND status = 'active'
			  ORDER BY expires_at DESC LIMIT 1`
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&userSub.ID, &userSub.UserID, &userSub.Plan, &userSub.Status, &userSub.StartedAT, &userSub.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &userSub, nil
}

func (r *subscriptionsRepository) MarkExpired(subID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE subscriptions SET status= 'expired' WHERE id=$1`
	_, err := r.db.Exec(ctx, query, subID)
	return err
}
