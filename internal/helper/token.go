package helper

import (
	"chera_khube/internal/dto"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type SignedDetail struct {
	PhoneNumber string
	Password    string
	jwt.RegisteredClaims
}

func GenerateAllToken(Token dto.Token) (token string, expireTime time.Time) {
	expireTime = time.Now().Add(time.Duration(Token.ExpireHour) * time.Hour)
	claims := SignedDetail{
		PhoneNumber: Token.PhoneNumber,
		Password:    Token.Password,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
		},
	}

	tokenModel := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, _ = tokenModel.SignedString([]byte(Token.SecretKey))

	return token, expireTime
}
