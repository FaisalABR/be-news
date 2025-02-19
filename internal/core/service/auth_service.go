package service

import (
	"bwa-news/config"
	"bwa-news/internal/adapter/repository"
	"bwa-news/internal/core/domain/entity"
	"bwa-news/lib/auth"
	"bwa-news/lib/conv"
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
)

var err error
var code string

type AuthService interface {
	GetUserByEmail(ctx context.Context, req entity.LoginRequest) (*entity.AccessToken, error)
}

type authService struct {
	authRepository repository.AuthRepository
	cfg            *config.Config
	jwtToken       auth.Jwt
}

func (a *authService) GetUserByEmail(ctx context.Context, req entity.LoginRequest) (*entity.AccessToken, error) {
	result, err := a.authRepository.GetUserByEmail(ctx, req)
	if err != nil {
		code = "[AUTHSERVICE] GetUserByEmail - 1"
		log.Errorw(code, err)
		return nil, err
	}

	if checkPass := conv.CheckPasswordHash(req.Password, result.Password); !checkPass {
		code = "[AUTHSERVICE] GetUserByEmail - 2"
		log.Errorw(code, "Invalid Password")
		err = errors.New("Invalid password")
		return nil, err
	}

	jwtData := entity.JWTData{
		UserID: float64(result.ID),
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
			ID:        string(result.ID),
		},
	}

	accesToken, expiresAt, err := a.jwtToken.GenerateToken(&jwtData)
	if err != nil {
		code = "[AUTHSERVICE] GetUserByEmail - 3"
		log.Errorw(code, err)
		return nil, err
	}

	resp := entity.AccessToken{
		AccessToken: accesToken,
		ExpiresAt:   expiresAt,
	}

	return &resp, nil
}

func NewAuthService(authRepository repository.AuthRepository, cfg *config.Config, jwtToken auth.Jwt) AuthService {
	return &authService{
		authRepository: authRepository,
		cfg:            cfg,
		jwtToken:       jwtToken,
	}
}
