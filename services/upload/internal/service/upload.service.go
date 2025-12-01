package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	db "github.com/Anurag-S1ngh/carbon-backend/pkg/db/generated"
	"github.com/Anurag-S1ngh/carbon-backend/pkg/rabbitmq"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UploadService struct {
	dbQueries      *db.Queries
	logger         *slog.Logger
	rabbitMqConfig *rabbitmq.RabbitMQConfig
}

func NewUploadService(dbQueries *db.Queries, rabbitmqConfig *rabbitmq.RabbitMQConfig, logger *slog.Logger) *UploadService {
	return &UploadService{
		dbQueries:      dbQueries,
		logger:         logger,
		rabbitMqConfig: rabbitmqConfig,
	}
}

func (s *UploadService) Upload(c *gin.Context, userUUID uuid.UUID, repoName string) error {
	user, err := s.dbQueries.GetUserByID(c, pgtype.UUID{
		Bytes: userUUID,
		Valid: true,
	})
	if err != nil {
		s.logger.Error("error while fetching user", "error", err)
		return err
	}

	githubAccessToken := user.GithubAccessToken.String
	if githubAccessToken == "" {
		return errors.New("user not connected to github")
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", user.GithubUsername.String, repoName)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		s.logger.Error("error while creating the request", "error", err)
		return err
	}
	req.Header.Set("Authorization", "Bearer "+githubAccessToken)

	res, err := client.Do(req)
	if err != nil {
		s.logger.Error("Error while sending request to check if user owns the repo", "error", err)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New("invalid repo")
	}

	s.logger.Debug("resonse", "res", res.Body)

	q, err := s.rabbitMqConfig.DeclareQueue("carbon:deploy")
	if err != nil {
		s.logger.Error("error while creating a queue", "error", err)
		return err
	}

	var body struct {
		UserUUID uuid.UUID `json:"userUUID"`
		RepoName string    `json:"repoName"`
	}
	body.RepoName = repoName
	body.UserUUID = userUUID

	bytes, err := json.Marshal(body)
	if err != nil {
		s.logger.Error("error while encoding body", "error", err)
		return err
	}

	err = s.rabbitMqConfig.Publish(q, bytes)
	if err != nil {
		s.logger.Error("error while publishing message to the queue", "error", err)
		return err
	}

	return nil
}
