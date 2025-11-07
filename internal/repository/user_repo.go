package repository

import (
	"Calsora/internal/models"
	"time"
)

type UserRepositoryInterface interface {
	Create(user *models.User) error
	GetById(id uint) (*models.User, error)
}

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepositury(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}
func (r *UserRepository) Create(user *model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO user (email, password, bday, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRow(ctx, query, user.Email, user.Password, user.Bday, user.CreatedAt).Scan(&user.ID)
}
