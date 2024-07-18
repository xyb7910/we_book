package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"we_book/internal/domain"
	"we_book/internal/repository"
)

var (
	ErrUserDuplicateEmail  = repository.ErrUserDuplicateEmail
	ErrUserInvalidPassword = errors.New("invalid password or email")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	// 考虑把加密放到service层
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}
