package service

import "log/slog"

type AuthService struct {
	logger *slog.Logger
}

func NewAuthService(logger *slog.Logger) *AuthService {
	return &AuthService{
		logger: logger,
	}
}
