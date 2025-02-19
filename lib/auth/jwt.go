package auth

import (
	"bwa-news/config"
	"bwa-news/internal/core/domain/entity"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Jwt interface {
	GenerateToken(data *entity.JWTData) (string, int64, error)
	VerifyAccessToken(token string) (*entity.JWTData, error)
}

type Options struct {
	signingkey string
	issuer     string
}

func (o *Options) GenerateToken(data *entity.JWTData) (string, int64, error) {
	// Prepared data that relate to time
	now := time.Now().Local()
	expiresAt := now.Add(time.Hour * 24)

	// Assign all data to registered claims
	data.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(expiresAt)
	data.RegisteredClaims.Issuer = o.issuer
	data.RegisteredClaims.NotBefore = jwt.NewNumericDate(now)
	acToken := jwt.NewWithClaims(jwt.SigningMethodHS256, data)

	accessToken, err := acToken.SignedString([]byte(o.signingkey))
	if err != nil {
		return "", 0, err
	}

	return accessToken, expiresAt.Unix(), nil
}

func (o *Options) VerifyAccessToken(token string) (*entity.JWTData, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method invalid")
		}
		return []byte(o.signingkey), nil
	})

	if err != nil {
		return nil, err
	}

	if parsedToken.Valid {
		claim, ok := parsedToken.Claims.(jwt.MapClaims)

		if !ok || !parsedToken.Valid {
			return nil, err
		}

		jwtData := &entity.JWTData{
			UserID: claim["user_id"].(float64),
		}

		return jwtData, nil
	}

	return nil, fmt.Errorf("token is not valid")
}

func NewJwt(cfg *config.Config) Jwt {
	opt := new(Options)
	opt.signingkey = cfg.App.JwtSecretKey
	opt.issuer = cfg.App.JwtIssuer

	return opt
}
