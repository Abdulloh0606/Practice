package repository

import (
	"context"
	"fmt"
)

func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users(name, email, password_hash)
	 VALUES($1, $2, $3);`

	_, err := r.postgres.Exec(ctx, query, user.Name, user.Email, user.PasswordHash)
	if err != nil {
		return fmt.Errorf("create user err %v", err)
	}
	fmt.Printf("Account created succesfully!")
	return nil
}
