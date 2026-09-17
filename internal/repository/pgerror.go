package repository

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/eduardocardona93/golang_shopping/pkg/apperrors"
)

// uniqueViolationCode is the PostgreSQL SQLSTATE for a unique_violation.
const uniqueViolationCode = "23505"

// translateWriteError maps low-level PostgreSQL errors into the sentinel
// apperrors used across the service layer, so callers never need to know
// about database-specific error types.
func translateWriteError(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
		return apperrors.ErrConflict
	}

	return err
}
