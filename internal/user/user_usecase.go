package user

import (
	"context"
	"errors"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

var (
	privateKey []byte = []byte(os.Getenv("JWT_PRIVATE_KEY"))
)

type UserRepo interface {
	IsUserExists(ctx context.Context, userID int) bool
}

type UserUsecase struct {
	repo UserRepo
}

func NewUserUsecase(repo UserRepo) *UserUsecase {
	return &UserUsecase{
		repo: repo,
	}
}

func (uu UserUsecase) IsUserExists(ctx context.Context, userID int) bool {
	if userID == 0 {
		return false
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return uu.repo.IsUserExists(ctx, userID)
}

func (uu UserUsecase) DecriptJWT(token string) (map[string]interface{}, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("auth invalid")
		}
		return privateKey, nil
	})

	data := make(map[string]interface{})
	if err != nil {
		return data, err
	}

	if !parsedToken.Valid {
		return data, errors.New("invalid token")
	}
	return parsedToken.Claims.(jwt.MapClaims), nil
}
