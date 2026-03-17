package localerrors

import (
	"errors"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	PgUniqueViolation     = "23505"
	PgForeignKeyViolation = "23503"
	PgNotNullViolation    = "23502"
	PgCheckViolation      = "23514"
	PgRestrictViolation   = "23001"
)

func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return fmt.Errorf(format, args...)
	}
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, args...), err)
}

func AsHumaError(err error) huma.StatusError {
	if err == nil {
		return huma.Error500InternalServerError("Internal Server Error")
	}

	var statusErr huma.StatusError
	if errors.As(err, &statusErr) {
		return statusErr
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return huma.Error404NotFound("Not Found")
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErrorToHumaError(pgErr)
	}

	return huma.Error500InternalServerError("Internal Server Error", err)
}

func pgErrorToHumaError(err *pgconn.PgError) huma.StatusError {
	if err == nil {
		return huma.Error500InternalServerError("Internal Server Error")
	}

	switch err.Code {
	case PgUniqueViolation:
		return huma.Error409Conflict("Conflict", errors.New(err.Message))
	case PgForeignKeyViolation:
		return huma.Error404NotFound("Not Found", errors.New(err.Message))
	case PgNotNullViolation:
		return huma.Error400BadRequest("Bad Request", errors.New(err.Message))
	case PgCheckViolation:
		return huma.Error400BadRequest("Bad Request", errors.New(err.Message))
	case PgRestrictViolation:
		return huma.Error400BadRequest("Bad Request", errors.New(err.Message))
	default:
		return huma.Error500InternalServerError("Internal Server Error", errors.New(err.Message))
	}
}
