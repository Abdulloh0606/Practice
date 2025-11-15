package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPgxCon() (*pgxpool.Pool, error) {
	ctx := context.Background()

	dsn := "postgres://postgres:abdulloh@localhost:5432/minitrello?sslmode=disable"

	conn, err := pgxpool.New(ctx, dsn)
	if err != nil {
		fmt.Println("conn error: ", err)
		return nil, err
	}

	return conn, nil
}
