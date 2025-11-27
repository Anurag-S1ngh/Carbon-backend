package jwt

import (
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtConfig struct {
	secret string
	logger *slog.Logger
}

func NewJwtConfig(secret string, logger *slog.Logger) *JwtConfig {
	return &JwtConfig{
		secret: secret,
		logger: logger,
	}
}

func (j *JwtConfig) GenerateJwt(userIDString string, expiresInHour uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userIDString,
		"iat":    time.Now().Unix(),
		"exp":    time.Now().Add(time.Hour * time.Duration(expiresInHour)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		j.logger.Error("error while signing token", "error", err)
		return "", err
	}

	return tokenString, nil
}

func (j *JwtConfig) VerifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(j.secret), nil
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
