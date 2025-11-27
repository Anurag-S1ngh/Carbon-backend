package service

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	db "github.com/Anurag-S1ngh/carbon-backend/pkg/db/generated"
	"github.com/Anurag-S1ngh/carbon-backend/pkg/email"
	"github.com/Anurag-S1ngh/carbon-backend/pkg/otp"
	"github.com/Anurag-S1ngh/carbon-backend/pkg/redis"
	"github.com/Anurag-S1ngh/carbon-backend/pkg/token/jwt"
	"github.com/Anurag-S1ngh/carbon-backend/pkg/token/token"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthService struct {
	logger                   *slog.Logger
	dbQueries                *db.Queries
	jwtConfig                *jwt.JwtConfig
	emailConfig              *email.EmailConfig
	redisConfig              *redis.RedisConfig
	otpExpirySeconds         string
	refreshTokenExpiresHours string
	accessTokenExpiresHours  string
}

func NewAuthService(dbQueries *db.Queries,
	otpExpirySeconds,
	refreshTokenExpiresInHour,
	accessTokenExpiresHours string,
	emailConfig *email.EmailConfig,
	redisConfig *redis.RedisConfig,
	jwtConfig *jwt.JwtConfig,
	logger *slog.Logger,
) *AuthService {
	return &AuthService{
		dbQueries:                dbQueries,
		otpExpirySeconds:         otpExpirySeconds,
		refreshTokenExpiresHours: refreshTokenExpiresInHour,
		accessTokenExpiresHours:  accessTokenExpiresHours,
		emailConfig:              emailConfig,
		redisConfig:              redisConfig,
		jwtConfig:                jwtConfig,
		logger:                   logger,
	}
}

func (s *AuthService) VerifyEmail(userEmail string) error {
	otpExpiryInSeconds, _ := strconv.Atoi(s.otpExpirySeconds)
	otpCode := otp.GenerateOTP()

	// FIXME: remove this in production
	s.logger.Debug("OTP", "otp", otpCode)

	key := fmt.Sprintf("otp:%s", userEmail)
	err := s.redisConfig.SetEx(key, otpCode, uint(otpExpiryInSeconds))
	if err != nil {
		return err
	}
	err = s.emailConfig.SendEmail(userEmail, otpCode)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) VerifyOTP(c *gin.Context, userEmail, userOTP string) error {
	accessTokenExpiresInHour, _ := strconv.Atoi(s.accessTokenExpiresHours)
	refreshTokenExpiresInHour, _ := strconv.Atoi(s.refreshTokenExpiresHours)
	if len(userOTP) != 6 {
		return errors.New("invalid otp")
	}

	redisKey := fmt.Sprintf("otp:%s", userEmail)
	originalOTP, err := s.redisConfig.Get(redisKey)
	if err != nil {
		return errors.New("try again later")
	}
	if originalOTP != userOTP {
		return errors.New("invalid otp")
	}

	existingUser, err := s.dbQueries.GetUserByEmail(c, userEmail)
	var user db.User
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			s.logger.Error("error while fetching the user from db", "error", err)
			return errors.New("try again later")
		}
		user, err = s.dbQueries.CreateUser(c, db.CreateUserParams{
			Email:           userEmail,
			ProfileImageUrl: pgtype.Text{},
		})
		if err != nil {
			s.logger.Error("error while creating user", "error", err)
			return errors.New("try again later")
		}
	} else {
		user = existingUser
	}

	s.logger.Debug("existing user", "debug", existingUser)

	refreshToken, err := token.GenerateRandomID(32)
	if err != nil {
		s.logger.Error("error while creating refresh token", "error", err)
		return errors.New("try again later")
	}
	refreshTokenHash := token.GenerateHash(refreshToken)

	accessToken, err := s.jwtConfig.GenerateJwt(user.ID.String(), uint(accessTokenExpiresInHour))
	if err != nil {
		return err
	}

	err = s.dbQueries.InsertRefreshToken(c, db.InsertRefreshTokenParams{
		UserID:    user.ID,
		HashToken: refreshTokenHash,
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now().Add(time.Hour * time.Duration(refreshTokenExpiresInHour)),
			Valid: true,
		},
	})
	if err != nil {
		return err
	}

	c.SetCookie("carbon-access-token", accessToken, accessTokenExpiresInHour*3600, "/", "", false, true)
	c.SetCookie("carbon-refresh-token", refreshToken, refreshTokenExpiresInHour*3600, "/", "", false, true)

	return nil
}

func (s *AuthService) RefreshAccessToken(c *gin.Context, refreshToken string) error {
	refreshTokenHash := token.GenerateHash(refreshToken)

	// FIXME: remove this in prod
	s.logger.Debug("refresh token", "debug", "flag -1")
	originalRefreshToken, err := s.dbQueries.GetRefreshTokenByToken(c, refreshTokenHash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			s.logger.Error("error while fetching refresh token", "error", err)
			return errors.New("try again later")
		}
		return errors.New("sign in again")
	}
	if time.Now().After(originalRefreshToken.ExpiresAt.Time) {
		return errors.New("sign in again")
	}
	// FIXME: remove this in prod
	s.logger.Debug("refresh token", "debug", "flag 0")

	accessTokenExpiresInHour, _ := strconv.Atoi(s.accessTokenExpiresHours)
	accessToken, err := s.jwtConfig.GenerateJwt(originalRefreshToken.UserID.String(), uint(accessTokenExpiresInHour))
	if err != nil {
		return err
	}
	// FIXME: remove this in prod
	s.logger.Debug("refresh token", "debug", "flag 1")

	c.SetCookie("carbon-access-token", accessToken, accessTokenExpiresInHour*3600, "/", "", false, true)

	return nil
}
