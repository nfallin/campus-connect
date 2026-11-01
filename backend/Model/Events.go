package Model

import (
	"campus_connect/Types"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Events
func (repository *Repository) CreateEvent(request Types.CreateEventRequest) (Types.EventResponse, error) {
	if request.HostUserID == nil {
		return Types.EventResponse{}, fmt.Errorf("%w: host user id is required", ErrInvalidEvent)
	}
	if request.Name == "" {
		return Types.EventResponse{}, fmt.Errorf("%w: name is required", ErrInvalidEvent)
	}
	if request.DateTime.IsZero() {
		return Types.EventResponse{}, fmt.Errorf("%w: date time is required", ErrInvalidEvent)
	}
	if request.Location == "" {
		return Types.EventResponse{}, fmt.Errorf("%w: location is required", ErrInvalidEvent)
	}
	if request.ClubName == "" {
		return Types.EventResponse{}, fmt.Errorf("%w: club name is required", ErrInvalidEvent)
	}
	if request.Description == "" {
		return Types.EventResponse{}, fmt.Errorf("%w: description is required", ErrInvalidEvent)
	}
	if err := requirePositiveID(*request.HostUserID); err != nil {
		return Types.EventResponse{}, fmt.Errorf("%w: host user id must be positive", ErrInvalidEvent)
	}

	if _, err := repository.GetUser(*request.HostUserID); err != nil {
		return Types.EventResponse{}, err
	}

	result, err := repository.db.Exec(`
		INSERT INTO events (
			event_name, date_time, location, club_name, description,
			contact_email, website, host_user_id
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, request.Name, request.DateTime.Format(time.RFC3339), request.Location,
		request.ClubName, request.Description, request.ContactEmail, request.WebsiteURL, *request.HostUserID)

	if err != nil {
		return Types.EventResponse{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Types.EventResponse{}, err
	}

	return repository.GetEvent(int(id))
}

func (repository *Repository) GetEvent(id int) (Types.EventResponse, error) {
	if err := requirePositiveID(id); err != nil {
		return Types.EventResponse{}, err
	}
	var event Types.EventResponse
	var dateTime string
	var host Types.UserResponse

	err := repository.db.QueryRow(`
		SELECT e.event_id, e.event_name, e.date_time, e.location,
			e.club_name, e.description, e.contact_email, e.website,
			e.host_user_id,
			u.user_id, u.user_name, u.user_email,
			(SELECT COUNT(*) FROM rsvps r WHERE r.event_id = e.event_id)
		FROM events e
		JOIN users u ON u.user_id = e.host_user_id
		WHERE e.event_id = ?
	`, id).Scan(&event.ID, &event.Name, &dateTime, &event.Location,
		&event.ClubName, &event.Description, &event.ContactEmail, &event.WebsiteURL,
		&event.HostUserID, &host.ID, &host.Name, &host.Email,
		&event.AttendeeCount)

	if errors.Is(err, sql.ErrNoRows) {
		return Types.EventResponse{}, ErrNotFound
	}

	if err != nil {
		return Types.EventResponse{}, err
	}

	event.DateTime, err = time.Parse(time.RFC3339, dateTime)
	if err != nil {
		return Types.EventResponse{}, err
	}

	event.Host = &host

	return event, nil
}

func (repository *Repository) GetEvents(filters Types.EventFilters) ([]Types.EventResponse, error) {
	query := `
		SELECT e.event_id, e.event_name, e.date_time, e.location,
			e.club_name, e.description, e.contact_email, e.website,
			e.host_user_id, u.user_id, u.user_name, u.user_email,
			(SELECT COUNT(*) FROM rsvps r WHERE r.event_id = e.event_id)
		FROM events e
		JOIN users u ON u.user_id = e.host_user_id
	`
	args := make([]any, 0)
	conditions := make([]string, 0)

	if filters.StartDate != nil {
		conditions = append(conditions, "e.date_time >= ?")
		args = append(args, filters.StartDate.Format(time.RFC3339))
	}

	if filters.EndDate != nil {
		conditions = append(conditions, "e.date_time <= ?")
		args = append(args, filters.EndDate.Format(time.RFC3339))
	}

	if filters.ClubName != "" {
		conditions = append(conditions, "e.club_name = ?")
		args = append(args, filters.ClubName)
	}

	if filters.Location != "" {
		conditions = append(conditions, "e.location = ?")
		args = append(args, filters.Location)
	}

	if filters.HostUserID != nil {
		conditions = append(conditions, "e.host_user_id = ?")
		args = append(args, *filters.HostUserID)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY e.date_time, e.event_id"

	rows, err := repository.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	events := make([]Types.EventResponse, 0)

	for rows.Next() {
		var event Types.EventResponse
		var dateTime string
		var host Types.UserResponse

		if err := rows.Scan(&event.ID, &event.Name, &dateTime, &event.Location,
			&event.ClubName, &event.Description, &event.ContactEmail, &event.WebsiteURL,
			&event.HostUserID, &host.ID, &host.Name, &host.Email, &event.AttendeeCount); err != nil {
			return nil, err
		}

		event.DateTime, err = time.Parse(time.RFC3339, dateTime)
		if err != nil {
			return nil, err
		}

		event.Host = &host
		events = append(events, event)
	}

	return events, rows.Err()
}

func (repository *Repository) UpdateEvent(id int, request Types.UpdateEventRequest) (Types.EventResponse, error) {
	if err := requirePositiveID(id); err != nil {
		return Types.EventResponse{}, err
	}
	current, err := repository.GetEvent(id)

	if err != nil {
		return Types.EventResponse{}, err
	}

	if request.Name != nil {
		if *request.Name == "" {
			return Types.EventResponse{}, fmt.Errorf("%w: name cannot be empty", ErrInvalidEvent)
		}
		current.Name = *request.Name
	}

	if request.DateTime != nil {
		if request.DateTime.IsZero() {
			return Types.EventResponse{}, fmt.Errorf("%w: date time cannot be zero", ErrInvalidEvent)
		}
		current.DateTime = *request.DateTime
	}

	if request.Location != nil {
		if *request.Location == "" {
			return Types.EventResponse{}, fmt.Errorf("%w: location cannot be empty", ErrInvalidEvent)
		}
		current.Location = *request.Location
	}

	if request.ClubName != nil {
		if *request.ClubName == "" {
			return Types.EventResponse{}, fmt.Errorf("%w: club name cannot be empty", ErrInvalidEvent)
		}
		current.ClubName = *request.ClubName
	}

	if request.Description != nil {
		if *request.Description == "" {
			return Types.EventResponse{}, fmt.Errorf("%w: description cannot be empty", ErrInvalidEvent)
		}
		current.Description = *request.Description
	}

	if request.ContactEmail != nil {
		current.ContactEmail = *request.ContactEmail
	}

	if request.WebsiteURL != nil {
		current.WebsiteURL = *request.WebsiteURL
	}

	if request.HostUserID != nil {
		if err := requirePositiveID(*request.HostUserID); err != nil {
			return Types.EventResponse{}, fmt.Errorf("%w: host user id must be positive", ErrInvalidEvent)
		}
		if _, err := repository.GetUser(*request.HostUserID); err != nil {
			return Types.EventResponse{}, err
		}

		current.HostUserID = *request.HostUserID
	}

	_, err = repository.db.Exec(`
		UPDATE events
		SET event_name = ?, date_time = ?, location = ?, club_name = ?, description = ?,
			contact_email = ?, website = ?, host_user_id = ?
		WHERE event_id = ?
	`, current.Name, current.DateTime.Format(time.RFC3339), current.Location,
		current.ClubName, current.Description, current.ContactEmail, current.WebsiteURL,
		current.HostUserID, id)

	if err != nil {
		return Types.EventResponse{}, err
	}

	return repository.GetEvent(id)
}

func (repository *Repository) DeleteEvent(id int) error {
	if err := requirePositiveID(id); err != nil {
		return err
	}
	result, err := repository.db.Exec(`DELETE FROM events WHERE event_id = ?`, id)

	if err != nil {
		return err
	}

	return requireAffected(result)
}
