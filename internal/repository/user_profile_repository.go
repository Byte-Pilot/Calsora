package repository

import (
	"Calsora/internal/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type UserProfileRepository interface {
	GetDailyIntake(profile *models.UserProfile, goal *models.UserGoal, dailyIntake *models.NutritionTarget) error
}

type userProfileRepository struct {
	db *pgxpool.Pool
}

func NewUserProfileRepository(db *pgxpool.Pool) *userProfileRepository {
	return &userProfileRepository{db: db}
}

func (r *userProfileRepository) GetDailyIntake(profile *models.UserProfile, goal *models.UserGoal, dailyIntake *models.NutritionTarget) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO user_profile(user_id, sex, height, weight, activity) 
					VALUES ($1, $2, $3, $4, $5)
                    ON CONFLICT (user_id)
					DO UPDATE SET sex=EXCLUDED.sex, height=EXCLUDED.height, weight=EXCLUDED.weight, activity=EXCLUDED.activity`
	_, err = tx.Exec(ctx, query, profile.UserID, profile.Sex, profile.Height, profile.Weight, profile.Activity)
	if err != nil {
		return err
	}

	query = `INSERT INTO user_goal (user_id, type, target_weight, weekly_rate) VALUES ($1, $2, $3, $4)`
	_, err = tx.Exec(ctx, query, goal.UserID, goal.Type, goal.TargetWeight, goal.WeeklyRate)
	if err != nil {
		return err
	}

	query = `INSERT INTO nutrition_target (user_id, cal, protein, carbs, fats) VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(ctx, query, dailyIntake.UserID, dailyIntake.Cal, dailyIntake.Protein, dailyIntake.Carbs, dailyIntake.Fats)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
