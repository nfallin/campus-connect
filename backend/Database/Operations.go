package Database

import (
	"campus_connect/Types"
	"database/sql"
	"errors"
	"strings"
	"time"
)

var ErrNotFound = errors.New("record not found")

// Users
func (repository *Repository) CreateUser(request Types.CreateUserRequest) (Types.UserResponse, bool, error) {
	var user Types.UserResponse

	err := repository.db.QueryRow(`
		SELECT user_id, user_name, user_email
		FROM users
		WHERE user_name = ? AND user_email = ?
	`, request.Name, request.Email).Scan(&user.ID, &user.Name, &user.Email)

	if err == nil {
		return user, false, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return Types.UserResponse{}, false, err
	}

	result, err := repository.db.Exec(`
		INSERT INTO users (user_name, user_email)
		VALUES (?, ?)
	`, request.Name, request.Email)

	if err != nil {
		return Types.UserResponse{}, false, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return Types.UserResponse{}, false, err
	}

	return Types.UserResponse{
		ID:    int(id),
		Name:  request.Name,
		Email: request.Email,
	}, true, nil
}

func (repository *Repository) GetUser(id int) (Types.UserResponse, error) {
	var user Types.UserResponse

	err := repository.db.QueryRow(`
		SELECT user_id, user_name, user_email
		FROM users
		WHERE user_id = ?
	`, id).Scan(&user.ID, &user.Name, &user.Email)

	if errors.Is(err, sql.ErrNoRows) {
		return Types.UserResponse{}, ErrNotFound
	}

	if err != nil {
		return Types.UserResponse{}, err
	}

	return user, nil
}

func (repository *Repository) GetAllUsers() ([]Types.UserResponse, error) {
	rows, err := repository.db.Query(`
		SELECT user_id, user_name, user_email
		FROM users
		ORDER BY user_id
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]Types.UserResponse, 0)

	for rows.Next() {
		var user Types.UserResponse

		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (repository *Repository) UpdateUser(id int, request Types.UpdateUserRequest) (Types.UserResponse, error) {
	current, err := repository.GetUser(id)

	if err != nil {
		return Types.UserResponse{}, err
	}

	if request.Name != nil {
		current.Name = *request.Name
	}

	if request.Email != nil {
		current.Email = *request.Email
	}

	_, err = repository.db.Exec(`
		UPDATE users
		SET user_name = ?, user_email = ?
		WHERE user_id = ?
	`, current.Name, current.Email, id)

	if err != nil {
		return Types.UserResponse{}, err
	}

	return current, nil
}

func (repository *Repository) DeleteUser(id int) error {
	result, err := repository.db.Exec(`DELETE FROM users WHERE user_id = ?`, id)

	if err != nil {
		return err
	}

	return requireAffected(result)
}

// Events
func (repository *Repository) CreateEvent(request Types.CreateEventRequest) (Types.EventResponse, error) {
	if request.HostUserID == nil {
		return Types.EventResponse{}, errors.New("host user id is required")
	}

	if _, err := repository.GetUser(*request.HostUserID); err != nil {
		return Types.EventResponse{}, err
	}

	result, err := repository.db.Exec(`
		INSERT INTO events (
			event_name, date_time, location, club_name,
			contact_email, website, host_user_id
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, request.Name, request.DateTime.Format(time.RFC3339), request.Location,
		request.ClubName, request.ContactEmail, request.WebsiteURL, *request.HostUserID)

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
	var event Types.EventResponse
	var dateTime string
	var host Types.UserResponse

	err := repository.db.QueryRow(`
		SELECT e.event_id, e.event_name, e.date_time, e.location,
			e.club_name, e.contact_email, e.website,
			e.host_user_id,
			u.user_id, u.user_name, u.user_email,
			(SELECT COUNT(*) FROM rsvps r WHERE r.event_id = e.event_id)
		FROM events e
		JOIN users u ON u.user_id = e.host_user_id
		WHERE e.event_id = ?
	`, id).Scan(&event.ID, &event.Name, &dateTime, &event.Location,
		&event.ClubName, &event.ContactEmail, &event.WebsiteURL,
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
			e.club_name, e.contact_email, e.website,
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
			&event.ClubName, &event.ContactEmail, &event.WebsiteURL,
			&event.HostUserID,
			&host.ID, &host.Name, &host.Email, &event.AttendeeCount); err != nil {
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
	current, err := repository.GetEvent(id)

	if err != nil {
		return Types.EventResponse{}, err
	}

	if request.Name != nil {
		current.Name = *request.Name
	}

	if request.DateTime != nil {
		current.DateTime = *request.DateTime
	}

	if request.Location != nil {
		current.Location = *request.Location
	}

	if request.ClubName != nil {
		current.ClubName = *request.ClubName
	}

	if request.ContactEmail != nil {
		current.ContactEmail = *request.ContactEmail
	}

	if request.WebsiteURL != nil {
		current.WebsiteURL = *request.WebsiteURL
	}

	if request.HostUserID != nil {
		if _, err := repository.GetUser(*request.HostUserID); err != nil {
			return Types.EventResponse{}, err
		}

		current.HostUserID = *request.HostUserID
	}

	_, err = repository.db.Exec(`
		UPDATE events
		SET event_name = ?, date_time = ?, location = ?, club_name = ?,
			contact_email = ?, website = ?, host_user_id = ?
		WHERE event_id = ?
	`, current.Name, current.DateTime.Format(time.RFC3339), current.Location,
		current.ClubName, current.ContactEmail, current.WebsiteURL,
		current.HostUserID, id)

	if err != nil {
		return Types.EventResponse{}, err
	}

	return repository.GetEvent(id)
}

func (repository *Repository) DeleteEvent(id int) error {
	result, err := repository.db.Exec(`DELETE FROM events WHERE event_id = ?`, id)

	if err != nil {
		return err
	}

	return requireAffected(result)
}

// RSVPs
func (repository *Repository) AddRSVP(userID, eventID int) error {
	if _, err := repository.GetUser(userID); err != nil {
		return err
	}

	if _, err := repository.GetEvent(eventID); err != nil {
		return err
	}

	_, err := repository.db.Exec(`
		INSERT INTO rsvps (user_id, event_id)
		VALUES (?, ?)
		ON CONFLICT(user_id, event_id) DO NOTHING
	`, userID, eventID)

	return err
}

func (repository *Repository) RemoveRSVP(userID, eventID int) error {
	if _, err := repository.GetUser(userID); err != nil {
		return err
	}

	if _, err := repository.GetEvent(eventID); err != nil {
		return err
	}

	_, err := repository.db.Exec(`
		DELETE FROM rsvps
		WHERE user_id = ? AND event_id = ?
	`, userID, eventID)

	return err
}

func (repository *Repository) GetUserRSVPs(userID int) ([]Types.EventResponse, error) {
	if _, err := repository.GetUser(userID); err != nil {
		return nil, err
	}

	rows, err := repository.db.Query(`
		SELECT e.event_id, e.event_name, e.date_time, e.location,
			e.club_name, e.contact_email, e.website, e.host_user_id,
			u.user_id, u.user_name, u.user_email,
			(SELECT COUNT(*) FROM rsvps all_rsvps WHERE all_rsvps.event_id = e.event_id)
		FROM events e
		JOIN rsvps r ON r.event_id = e.event_id
		JOIN users u ON u.user_id = e.host_user_id
		WHERE r.user_id = ?
		ORDER BY e.date_time, e.event_id
	`, userID)

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
			&event.ClubName, &event.ContactEmail, &event.WebsiteURL,
			&event.HostUserID, &host.ID, &host.Name, &host.Email,
			&event.AttendeeCount); err != nil {
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
