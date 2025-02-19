package repository

import (
	"bwa-news/internal/core/domain/entity"
	"bwa-news/internal/core/domain/model"
	"context"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type UserRepostiory interface {
	GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error)
	UpdatePassword(ctx context.Context, newPass string, id int64) error
}

type userRepostiory struct {
	db *gorm.DB
}

func (u *userRepostiory) GetUserByID(ctx context.Context, id int64) (*entity.UserEntity, error) {
	var modelUser model.User

	err = u.db.Where("id = ?", id).First(&modelUser).Error
	if err != nil {
		code = "[REPOSITORY] GetUserByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return &entity.UserEntity{
		ID:       id,
		Username: modelUser.Username,
		Email:    modelUser.Email,
	}, nil

}
func (u *userRepostiory) UpdatePassword(ctx context.Context, newPass string, id int64) error {

	err = u.db.Model(&model.User{}).Where("id = ?", id).Update("password", newPass).Error
	if err != nil {
		code = "[REPOSITORY] UpdatePassword - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func NewUserRepository(db *gorm.DB) UserRepostiory {
	return &userRepostiory{
		db: db,
	}
}
