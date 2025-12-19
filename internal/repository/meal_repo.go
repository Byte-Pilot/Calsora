package repository

import (
	"Calsora/internal/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type MealRepositoryInterface interface {
	CreateMeal(meal *models.Meals) error
	DeleteMeal(meal int) error
}

type MealRepository struct {
	db *pgxpool.Pool
}

func NewMealRepository(db *pgxpool.Pool) *MealRepository {
	return &MealRepository{db: db}
}

func (m *MealRepository) CreateMeal(meals *models.Meals) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `INSERT INTO meals (user_id, name, cal, protein, carbs, fats) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`
	if err := m.db.QueryRow(ctx, query, meals.UserID, meals.Name, meals.Cal, meals.Protein, meals.Carbs, meals.Fats).Scan(&meals.ID, &meals.CreatedAt); err != nil {
		return err
	}
	return nil
}

func (m *MealRepository) DeleteMeal(meal int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `DELETE FROM meals WHERE id = $1`
	_, err := m.db.Exec(ctx, query, meal)
	if err != nil {
		return err
	}
	return nil
}
