package Model

import (
	"database/sql"
	"errors"
)

// SaveEventImage stores or replaces an event's image.
func (repository *Repository) SaveEventImage(id int, data []byte, contentType string) error {
	if err := requirePositiveID(id); err != nil {
		return err
	}

	result, err := repository.db.Exec(`
		UPDATE events
		SET image_data = ?, image_content_type = ?
		WHERE event_id = ?
	`, data, contentType, id)

	if err != nil {
		return err
	}

	return requireAffected(result)
}

// GetEventImage retrieves an event's image and its content type.
func (repository *Repository) GetEventImage(id int) ([]byte, string, error) {
	if err := requirePositiveID(id); err != nil {
		return nil, "", err
	}

	var data []byte
	var contentType string

	err := repository.db.QueryRow(`
		SELECT image_data, image_content_type
		FROM events
		WHERE event_id = ? AND image_data IS NOT NULL
	`, id).Scan(&data, &contentType)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, "", ErrNotFound
	}
	if err != nil {
		return nil, "", err
	}

	return data, contentType, nil
}

// DeleteEventImage removes an event's image.
func (repository *Repository) DeleteEventImage(id int) error {
	if err := requirePositiveID(id); err != nil {
		return err
	}

	result, err := repository.db.Exec(`
		UPDATE events
		SET image_data = NULL, image_content_type = NULL
		WHERE event_id = ?
	`, id)

	if err != nil {
		return err
	}

	return requireAffected(result)
}
