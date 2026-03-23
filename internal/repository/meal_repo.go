package repository

import (
	"Calsora/internal/models"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type MealRepository interface {
	CreateMeal(meal *models.Meals, items []*models.MealItems) (mealID int, createdAt time.Time, itemsNew []*models.MealItems, err error)
	EditMeal(meal *models.Meals, items []*models.MealItems) ([]*models.MealItems, error)
	UpdateMeal(meal *models.Meals, items []*models.MealItems) error
	GetMealData(userID, mealID int) (meal *models.Meals, items []*models.MealItems, err error)
	GetDailyNutritionStats(userID, days int) ([]*models.Stats, error)
	DeleteMeal(userID, mealID int) error
}

type mealRepository struct {
	db *pgxpool.Pool
}

func NewMealRepository(db *pgxpool.Pool) *mealRepository {
	return &mealRepository{db: db}
}

func (m *mealRepository) CreateMeal(meal *models.Meals, items []*models.MealItems) (mealID int, createdAt time.Time, itemsNew []*models.MealItems, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := m.db.Begin(ctx)
	if err != nil {
		return 0, time.Time{}, nil, err
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO meals (user_id, name) VALUES ($1, $2) RETURNING id, created_at`
	if err := tx.QueryRow(ctx, query, meal.UserID, meal.Name).Scan(&meal.ID, &meal.CreatedAt); err != nil {
		return 0, time.Time{}, nil, err
	}

	query = `INSERT INTO meal_items (meal_id, name, grams, cal, protein, carbs, fats, confidence) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, meal_id`
	for i, item := range items {
		if err := tx.QueryRow(ctx, query, meal.ID, item.Name, item.Grams, item.Cal, item.Protein, item.Carbs, item.Fats, item.Confidence).Scan(&items[i].ID, &items[i].MealID); err != nil {
			return 0, time.Time{}, nil, err
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return 0, time.Time{}, nil, err
	}
	return meal.ID, meal.CreatedAt, items, nil
}

func (m *mealRepository) EditMeal(meal *models.Meals, items []*models.MealItems) ([]*models.MealItems, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := m.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := `UPDATE meals SET name=$1 WHERE id=$2 AND user_id=$3`
	cmdTag, err := tx.Exec(ctx, query, meal.Name, meal.ID, meal.UserID)
	if err != nil {
		return nil, err
	}
	if cmdTag.RowsAffected() == 0 {
		return nil, errors.New("meal not found")
	}

	query = `DELETE FROM meal_items WHERE meal_id=$1`
	cmdTag, err = tx.Exec(ctx, query, meal.ID)
	if err != nil {
		return nil, err
	}

	query = `INSERT INTO meal_items (meal_id, name, grams, cal, protein, carbs, fats, confidence) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, meal_id, created_at`
	for i, item := range items {
		if err := tx.QueryRow(ctx, query, meal.ID, item.Name, item.Grams, item.Cal, item.Protein, item.Carbs, item.Fats, item.Confidence).Scan(&items[i].ID, &items[i].MealID, &items[i].CreatedAt); err != nil {
			return nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (m *mealRepository) UpdateMeal(meal *models.Meals, items []*models.MealItems) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := m.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `UPDATE meals SET name=$1 WHERE user_id=$2 AND id=$3`
	cmdTag, err := tx.Exec(ctx, query, meal.Name, meal.UserID, meal.ID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("meal not found")
	}

	query = `UPDATE meal_items SET name=$1, grams=$2, cal=$3, protein=$4, carbs=$5, fats=$6 WHERE id=$7 AND  meal_id=$8`
	for _, item := range items {
		if _, err := tx.Exec(ctx, query, item.Name, item.Grams, item.Cal, item.Protein, item.Carbs, item.Fats, item.ID, meal.ID); err != nil {
			return err
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (m *mealRepository) GetMealData(userID, mealID int) (*models.Meals, []*models.MealItems, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	meal := &models.Meals{}
	var items []*models.MealItems

	query := `SELECT name, created_at FROM meals WHERE user_id=$1 AND id=$2`
	err := m.db.QueryRow(ctx, query, userID, mealID).Scan(&meal.Name, &meal.CreatedAt)
	if err != nil {
		return nil, nil, err
	}

	query = `SELECT id, name, grams, cal, protein, carbs, fats, confidence, created_at FROM meal_items WHERE meal_id=$1`
	rows, err := m.db.Query(ctx, query, mealID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &models.MealItems{}

		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Grams,
			&item.Cal,
			&item.Protein,
			&item.Carbs,
			&item.Fats,
			&item.Confidence,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return meal, items, nil
}

func (m *mealRepository) GetDailyNutritionStats(userID, days int) ([]*models.Stats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var stats []*models.Stats
	query := `SELECT DATE(m.created_at) as day, 
       COALESCE(SUM(mi.cal), 0) as cal, 
       COALESCE(SUM(mi.protein), 0) as protein, 
       COALESCE(SUM(mi.carbs), 0) as carbs, 
       COALESCE(SUM(mi.fats), 0) as fats 
	   FROM meals m 
	   JOIN meal_items mi ON mi.meal_id = m.id 
	   WHERE m.user_id = $1 AND m.created_at >= CURRENT_DATE - ($2 * INTERVAL '1 day')
	   GROUP BY day 
	   ORDER BY day;`
	rows, err := m.db.Query(ctx, query, userID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dayStats models.Stats
		err := rows.Scan(&dayStats.Day, &dayStats.Cal, &dayStats.Protein, &dayStats.Carbs, &dayStats.Fats)
		if err != nil {
			return nil, err
		}
		stats = append(stats, &dayStats)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	log.Println(stats)
	return stats, nil
}

func (m *mealRepository) DeleteMeal(userID, mealID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `DELETE FROM meals WHERE user_id = $1 AND id = $2`
	cmdTag, err := m.db.Exec(ctx, query, userID, mealID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("meal not found")
	}
	return nil
}
