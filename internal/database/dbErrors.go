package database

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBError struct {
	Message string
	Repo    string
	Method  string
	Code    string
}

func NewDBError(message, repo, method, code string) *DBError {
	return &DBError{
		Message: message,
		Repo:    repo,
		Method:  method,
		Code:    code,
	}
}

func (e *DBError) Error() string {
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
	stringDataRightTruncationCode = "22001" // A string is longer than the column allows.
	invalidTextRepresentationCode = "22P02" // PostgreSQL can’t cast a string to the target data type.
	notNullViolationCode          = "23502"
	foreignKeyViolationCode       = "23503"
	uniqueViolationCode           = "23505"
	checkViolationCode            = "23514"
)

func Parse(err error, repo, method string) *DBError {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return NewDBError(ErrNotFound, repo, method, "")
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case uniqueViolationCode:
			return NewDBError(ErrDuplicate, repo, method, uniqueViolationCode)
		case foreignKeyViolationCode:
			return NewDBError(ErrForeignKey, repo, method, foreignKeyViolationCode)
		case notNullViolationCode:
			return NewDBError(ErrNullViolation, repo, method, notNullViolationCode)
		case checkViolationCode:
			return NewDBError(ErrCheckViolation, repo, method, checkViolationCode)
		case stringDataRightTruncationCode:
			return NewDBError(ErrStringTooLong, repo, method, stringDataRightTruncationCode)
		case invalidTextRepresentationCode:
			return NewDBError(ErrInvalidFormat, repo, method, invalidTextRepresentationCode)
		default:
			return NewDBError(ErrUnknown, repo, method, "")
		}
	}

	return NewDBError(ErrUnknown, repo, method, "")
}
