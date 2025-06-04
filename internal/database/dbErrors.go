package database

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrDuplicate      = errors.New("duplicate record")
	ErrForeignKey     = errors.New("related record not found")
	ErrNotFound       = errors.New("record not found")
	ErrInvalidInput   = errors.New("invalid input")
	ErrStringTooLong  = errors.New("input string is too long")
	ErrInvalidFormat  = errors.New("invalid input format")
	ErrNullViolation  = errors.New("missing required field")
	ErrCheckViolation = errors.New("invalid value for one of the fields")
	ErrUnknown        = errors.New("unknown database error")
)

var (
	stringDataRightTruncationCode = "22001" // A string is longer than the column allows.
	invalidTextRepresentationCode = "22P02" // PostgreSQL can’t cast a string to the target data type.
	notNullViolationCode          = "23502"
	foreignKeyViolationCode       = "23503"
	uniqueViolationCode           = "23505"
	checkViolationCode            = "23514"
)

func Parse(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case uniqueViolationCode:
			return ErrDuplicate
		case foreignKeyViolationCode:
			return ErrForeignKey
		case notNullViolationCode:
			return ErrNullViolation
		case checkViolationCode:
			return ErrCheckViolation
		case stringDataRightTruncationCode:
			return ErrStringTooLong
		case invalidTextRepresentationCode:
			return ErrInvalidFormat
		default:
			return ErrUnknown
		}
	}

	return err
}
