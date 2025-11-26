package jwt

import (
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	secret string
	logger *slog.Logger
}

func (j *JWTConfig) GeneratToken(userIDString string, expiresInHour uint8) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userIDString,
		"exp":    time.Now().Add(time.Hour * time.Duration(expiresInHour)),
	})

	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		j.logger.Error("error while signing token", "error", err)
		return "", err
	}

	return tokenString, nil
}

func (j *JWTConfig) VerifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return j.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		j.logger.Error("error while verifying token", "error", err)
		return "", err
	}

	var userIDString string
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userIDString = claims["userID"].(string)
	} else {
		j.logger.Error("error while verifying token", "error", err)
		return "", err
	}

	return userIDString, nil
}
