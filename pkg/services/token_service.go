package services

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"workout-tracker-go-app/pkg/utils"
)

func CreateToken(userId uint, tokenDurationInMinutes int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"exp": utils.CurrentTimePlusMinutesAsUnix(tokenDurationInMinutes),
	})
	secret := os.Getenv("SECRET")
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenString string) (bool, int64, int, error) {
	secret := os.Getenv("SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return false, 0, 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, 0, 0, err
	}

	return token.Valid, int64(claims["exp"].(float64)), int(claims["sub"].(float64)), nil
}
