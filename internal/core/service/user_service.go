package service

import (
	"bwa-news/internal/adapter/repository"
	"bwa-news/internal/core/domain/entity"
	"bwa-news/lib/conv"
	"context"

	"github.com/gofiber/fiber/v2/log"
)

type UserService interface {
	GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error)
	UpdatePassword(ctx context.Context, newPass string, id int64) error
}

type userService struct {
	userRepostiory repository.UserRepostiory
}

func (u *userService) GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error) {
	result, err := u.userRepostiory.GetUserByID(ctx, id)
	if err != nil {
		code = "[SERVICE] GetUserByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return result, nil
}
func (u *userService) UpdatePassword(ctx context.Context, newPass string, id int64) error {
	password, err := conv.HashPassword(newPass)
	if err != nil {
		code = "[SERVICE] UpdatePassword - 1"
		log.Errorw(code, err)
		return err
	}

	err = u.userRepostiory.UpdatePassword(ctx, password, id)
	if err != nil {
		code = "[SERVICE] UpdatePassword - 2"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewUserService(userRepository repository.UserRepostiory) UserService {
	return &userService{
		userRepostiory: userRepository,
	}
}
