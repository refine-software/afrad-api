package database

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Constraints map[string]string

type DBError struct {
	Message string
	Repo    string
	Method  string
	Code    string
	Err     error
	Column  string
}

func NewDBError(err error, message, repo, method, code, constraint string) DBError {
	return DBError{
		Message: message,
		Repo:    repo,
		Method:  method,
		Code:    code,
		Err:     err,
		Column:  constraint,
	}
}

func (e DBError) Error() string {
	return e.Message
}

var (
	ErrDuplicate      = "duplicate record"
	ErrForeignKey     = "related record not found"
	ErrNotFound       = "record not found"
	ErrInvalidInput   = "invalid input"
	ErrStringTooLong  = "input string is too long"
	ErrInvalidFormat  = "invalid input format"
	ErrNullViolation  = "missing required field"
	ErrCheckViolation = "invalid value for one of the fields"
	ErrUnknown        = "unknown database error"
)

var (
	StringDataRightTruncationCode = "22001" // A string is longer than the column allows.
	InvalidTextRepresentationCode = "22P02" // PostgreSQL can’t cast a string to the target data type.
	NotNullViolationCode          = "23502"
	ForeignKeyViolationCode       = "23503"
	UniqueViolationCode           = "23505"
	CheckViolationCode            = "23514"
)

// Parse interprets a database error and returns a more descriptive application-level error.
// It maps known PostgreSQL error codes (e.g., unique violations, foreign key violations)
// to domain-specific errors using the provided constraints map, and attaches contextual
// information such as the repository and method where the error occurred.
func Parse(err error, repo, method string, constraints Constraints) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return NewDBError(err, ErrNotFound, repo, method, "", "")
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case UniqueViolationCode:
			return NewDBError(
				err,
				ErrDuplicate,
				repo,
				method,
				UniqueViolationCode,
				constraints[UniqueViolationCode],
			)
		case ForeignKeyViolationCode:
			return NewDBError(
				err,
				ErrForeignKey,
				repo,
				method,
				ForeignKeyViolationCode,
				constraints[ForeignKeyViolationCode],
			)
		case NotNullViolationCode:
			return NewDBError(
				err,
				ErrNullViolation,
				repo,
				method,
				NotNullViolationCode,
				constraints[NotNullViolationCode],
			)
		case CheckViolationCode:
			return NewDBError(
				err,
				ErrCheckViolation,
				repo,
				method,
				CheckViolationCode,
				constraints[CheckViolationCode],
			)
		case StringDataRightTruncationCode:
			return NewDBError(
				err,
				ErrStringTooLong,
				repo,
				method,
				StringDataRightTruncationCode,
				constraints[StringDataRightTruncationCode],
			)
		case InvalidTextRepresentationCode:
			return NewDBError(
				err,
				ErrInvalidFormat,
				repo,
				method,
				InvalidTextRepresentationCode,
				constraints[InvalidTextRepresentationCode],
			)
		}
	}

	return NewDBError(err, ErrUnknown, repo, method, "", "")
}

func IsDBNotFoundErr(err error) bool {
	var dbErr DBError
	ok := errors.As(err, &dbErr)
	if !ok {
		return false
	}

	return dbErr.Message == ErrNotFound
}
