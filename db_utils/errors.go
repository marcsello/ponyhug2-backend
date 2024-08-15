package db_utils

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
)

func isPgErrTypeOf(err error, severity, code string) bool {
	if err == nil {
		return false
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Severity == severity && pgErr.Code == code
	}
	return false
}

func IsDuplicatedKeyErr(err error) bool {
	return isPgErrTypeOf(err, "ERROR", "23505")
}

func IsForeignKeyViolationErr(err error) bool {
	return isPgErrTypeOf(err, "ERROR", "23503")
}

func IsNotNullViolation(err error) bool {
	return isPgErrTypeOf(err, "ERROR", "23502")
}
