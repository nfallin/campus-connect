package Database

import (
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("record not found")

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