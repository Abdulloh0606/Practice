package service

import (
	"context"
	"errors"
	"fmt"
	"minitrello/pkg/errs"
	"minitrello/pkg/models"

	"minitrello/pkg/auth"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) RegisterUser(ctx context.Context, user *models.UserInput) error {
	//хэширование
	cost := 10
	pHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), cost)
	if err != nil {
		return fmt.Errorf("hash password err: %w", err)
	}

	newUser := &models.User{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: string(pHash),
	}
	err = s.repo.CreateUser(ctx, newUser)
	if err != nil {
		fmt.Printf("error: %v", err)
		if errors.Is(err, errs.ErrEmailAlreadyExists) {
			return errs.ErrEmailAlreadyExists
		}
		return err
	}
	return nil
}
func (s *Service) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

func (s *Service) GetByID(ctx context.Context, id int) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

func (s *Service) LoginUser(ctx context.Context, email, password string) (string, error) {

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errs.ErrInvalidPassword
	}

	token, err := auth.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) UpdateUserName(ctx context.Context, userID int, newName string) error {

	return s.repo.UpdateUserName(ctx, userID, newName)

}

func (s *Service) DeleteUser(ctx context.Context, id int) error {
	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("delete user error: %w", err)
	}
	return nil
}
