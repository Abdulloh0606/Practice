package repository

import (

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	postgres *pgxpool.Pool
}

func NewRepository(postgres *pgxpool.Pool) *Repository {
	return &Repository{
		postgres: postgres,
	}
}


