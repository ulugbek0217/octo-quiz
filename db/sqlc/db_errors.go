package db

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	NoDataFound = "P0002"
)

func ErrorCode(err error) string {
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		return pgError.Code
	}
	return ""
}
