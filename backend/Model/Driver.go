package Model

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

const enableForeignKeys = `PRAGMA foreign_keys=ON`

func Init(path string) (*sql.DB, error) {
	db, err := open(path)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if err := createTables(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	if _, err := db.Exec(enableForeignKeys); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	if _, err := db.Exec(usersTable); err != nil {
		return err
	}
	if _, err := db.Exec(eventsTable); err != nil {
		return err
	}
	if _, err := db.Exec(rsvpTable); err != nil {
		return err
	}
	// Migrate existing databases.
	if err := ensureColumn(db, "events", "image_data", "BLOB"); err != nil {
		return err
	}
	if err := ensureColumn(db, "events", "image_content_type", "TEXT"); err != nil {
		return err
	}
	
	return nil
}

func ensureColumn(db *sql.DB, table, column, columnType string) error {
	rows, err := db.Query("PRAGMA table_info(" + table + ")")
	if err != nil {
		return err
	}

	found := false

	for rows.Next() {
		var cid, notNull, primaryKey int
		var name, dataType string
		var defaultValue sql.NullString

		if err := rows.Scan(
			&cid, &name, &dataType, &notNull, &defaultValue, &primaryKey,
		); err != nil {
			rows.Close()
			return err
		}

		if name == column {
			found = true
		}
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}
	rows.Close()

	if found {
		return nil
	}

	_, err = db.Exec(
		"ALTER TABLE " + table + " ADD COLUMN " + column + " " + columnType,
	)
	return err
}
