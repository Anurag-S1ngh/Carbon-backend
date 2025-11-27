package db

import (
	"context"
	"fmt"

	db "github.com/Anurag-S1ngh/carbon-backend/pkg/db/generated"
	"github.com/jackc/pgx/v5"
)

func NewDatabaseConnection(dbURL string) (*pgx.Conn, error) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		fmt.Println("error while creating a connection to database", "error", err)
		return nil, err
	}

	return conn, nil
}

func DatabaseQueries(dbConn *pgx.Conn) *db.Queries {
	queries := db.New(dbConn)

	return queries
}
