package Model

import (
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("record not found")
var ErrInvalidID = errors.New("id must be a positive integer")
var ErrInvalidEvent = errors.New("invalid event")

func requirePositiveID(id int) error {
	if id <= 0 {
		return ErrInvalidID
	}
	return nil
}

func requireAffected(result sql.Result) error {
	affected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrNotFound
	}

	return nil
}
