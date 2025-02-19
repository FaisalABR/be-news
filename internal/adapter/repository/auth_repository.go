package repository

import (
	"bwa-news/internal/core/domain/entity"
	"bwa-news/internal/core/domain/model"
	"context"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

var code string
var err error

type AuthRepository interface {
	GetUserByEmail(ctx context.Context, req entity.LoginRequest) (*entity.UserEntity, error)
}

type authRepository struct {
	db *gorm.DB
}

func (a *authRepository) GetUserByEmail(ctx context.Context, req entity.LoginRequest) (*entity.UserEntity, error) {
	var modelUser model.User

	err := a.db.Where("email = ?", req.Email).First(&modelUser).Error
	if err != nil {
		code = "[AUTHREPOSITORY] GetUserByEmail - 1"
		log.Errorw(code, err)
		return nil, err
	}

	resp := entity.UserEntity{
		ID:       modelUser.ID,
		Username: modelUser.Username,
		Email:    modelUser.Email,
		Password: modelUser.Password,
	}

	return &resp, nil
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}
