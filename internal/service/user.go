package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"we_book/internal/domain"
	"we_book/internal/repository"
)

var (
	ErrUserDuplicateEmail    = repository.ErrUserDuplicate
	ErrInvalidUserOrPassword = errors.New("invalid password or email")
)

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email, password string) (domain.User, error)
	Edit(ctx *gin.Context, info domain.User) error
	FindById(ctx context.Context, id int64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) SignUp(ctx context.Context, u domain.User) error {
	// 考虑把加密放到service层
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *userService) Login(ctx context.Context, email, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	// 进行密码校验
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, err
}

func (svc *userService) Edit(ctx *gin.Context, info domain.User) error {
	return svc.repo.Edit(ctx, info)
}

func (svc *userService) FindById(ctx context.Context, id int64) (domain.User, error) {
	return svc.repo.FindById(ctx, id)
}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	user, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		return user, err
	}
	u := domain.User{
		Phone: phone,
	}

	err = svc.repo.Create(ctx, u)
	if err != nil || err == repository.ErrUserDuplicate {
		return domain.User{}, err
	}
	return svc.repo.FindByPhone(ctx, phone)
}
