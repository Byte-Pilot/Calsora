package repository

import (
	"Calsora/internal/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type UserRepository interface {
	Create(user *models.User) error
	GetById(id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	ChangePass(id int, newPass string) error
	DeleteId(id int) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *userRepository {
	return &userRepository{db: db}
}
func (r *userRepository) Create(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO users (email, password, bday) VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, user.Email, user.Password, user.Bday).Scan(&user.ID, &user.CreatedAt)
}

func (r *userRepository) GetById(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User
	query := `SELECT id, email, password, bday, created_at FROM users WHERE id=$1`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.Bday, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User
	query := `SELECT id, email, password, bday, created_at FROM users WHERE email=$1`
	err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.Bday, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) ChangePass(id int, newPass string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE users SET password=$1 WHERE id=$2`
	_, err := r.db.Exec(ctx, query, newPass, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) DeleteId(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM users WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
