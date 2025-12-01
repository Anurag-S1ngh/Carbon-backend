package service

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	db "github.com/Anurag-S1ngh/carbon-backend/pkg/db/generated"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type GitService struct {
	dbQueries *db.Queries
	logger    *slog.Logger
}

func NewGitService(dbQueries *db.Queries, logger *slog.Logger) *GitService {
	return &GitService{
		dbQueries: dbQueries,
		logger:    logger,
	}
}

func (s *GitService) ListAllRepos(c *gin.Context, userUUID uuid.UUID) ([]map[string]any, error) {
	var repos []map[string]any
	user, err := s.dbQueries.GetUserByID(c, pgtype.UUID{
		Bytes: [16]byte(userUUID),
		Valid: true,
	})
	if err != nil {
		s.logger.Error("error while fetching user by id", "error", err)
		return repos, err
	}

	githubAccessToken := user.GithubAccessToken.String
	if githubAccessToken == "" {
		return repos, errors.New("user not connected to github")
	}

	url := "https://api.github.com/user/repos"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+githubAccessToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	res, err := client.Do(req)
	if err != nil {
		s.logger.Error("error while sending request to github to fetch all repos", "error", err)
		return repos, err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&repos); err != nil {
		s.logger.Error("error while decoding json", "error", err)
		return repos, err
	}

	return repos, nil
}
