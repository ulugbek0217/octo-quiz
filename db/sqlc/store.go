package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries
	Pool *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return Store{
		Queries: New(),
		Pool:    db,
	}
}
