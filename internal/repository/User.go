package repository

import (
	"context"
	"errors"
	"minitrello/pkg/errs"
	"minitrello/pkg/models"
	"strings"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	query := `insert into users(name, email, password_hash)
	 VALUES($1, $2, $3);`

	_, err := r.postgres.Exec(ctx, query, user.Name, user.Email, user.PasswordHash)
	if err != nil {
		if strings.Contains(err.Error(), "users_email_key") {
			return errs.ErrEmailAlreadyExists
		}
		return err
	}

	return nil
}
func (r *Repository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
        select id, name, email, password_hash, role, created_at  
		from users 
		where email = $1
    `

	var user models.User
	err := r.postgres.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `
        select id, name, email, password_hash, role, created_at  
		from users
		where id = $1
    `

	var user models.User
	err := r.postgres.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *Repository) DeleteUser(ctx context.Context, id int) error {
	query := `
		delete from users where id = $1
	`
	result, err := r.postgres.Exec(ctx, query, id)
	if err != nil {
		if result.RowsAffected() == 0 {
			return errs.ErrUserNotFound
		}
		return err
	}
	return nil
}

func (r *Repository) UpdateUserName(ctx context.Context, userID int, newName string) error {
	query := `
		update users
		set name = $1
		where id = $2
	`
	_, err := r.postgres.Exec(ctx, query, newName, userID)
	if err != nil {
		return err
	}

	return nil
}
