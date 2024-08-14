package db_utils

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
)

// resolveResultToProperError tries to catch all errors uncaught by gorm and translate them to gorm errors
func resolveError(err error) error {
	if err == nil {
		return nil
	}

	// There was an error

	var pgErr *pgconn.PgError
	if errors.As(result.Error, &pgErr) {
		if pgErr.Severity == "ERROR" && pgErr.Code == "23505" { // gorm, stupid as always, and can not handle unique key constarint violations...
			// Log original error? I think it's auto logged...
			return gorm.ErrDuplicatedKey
		}

		if pgErr.Severity == "ERROR" && pgErr.Code == "23503" { // gorm, stupid as always, and can not handle foreign key violation...
			return ErrForeignKeyViolation // <- stupid gorm does not even have error of this type
		}

	}

	return result.Error

}
