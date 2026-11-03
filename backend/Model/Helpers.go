package Model

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

var ErrNotFound = errors.New("record not found")
var ErrInvalidID = errors.New("id must be a positive integer")
var ErrInvalidEvent = errors.New("invalid event")
var ErrDuplicateEventName = errors.New("Event name must be unique.")

func requirePositiveID(id int) error {
	if id <= 0 {
		return ErrInvalidID
	}
	return nil
}

func requireEventDateTime(dateTime time.Time) error {
	now := time.Now()
	if dateTime.Before(now) {
		return fmt.Errorf("%w: date time cannot be in the past", ErrInvalidEvent)
	}
	if dateTime.After(now.AddDate(0, 6, 0)) {
		return fmt.Errorf("%w: date time cannot be more than six months in the future", ErrInvalidEvent)
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

func translateEventWriteError(err error) error {
	var sqliteError *sqlite.Error
	if errors.As(err, &sqliteError) && sqliteError.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
		return ErrDuplicateEventName
	}
	return err
}
