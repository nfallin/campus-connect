package Model

import (
	"campus_connect/Types"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf16"
)

const (
	maxEventTitleLength        = 75
	maxEventOrganizationLength = 50
	maxEventLocationLength     = 50
	maxEventDescriptionLength  = 500
	maxEventEmailLength        = 30
	maxEventWebsiteLength      = 60
	maxEventInterests          = 5
)

// Events
func validateEventFieldLength(value, field string, maxLength int) error {
	length := 0
	for _, character := range value {
		length += utf16.RuneLen(character)
	}
	if length > maxLength {
		return fmt.Errorf("%w: %s must be %d characters or fewer", ErrInvalidEvent, field, maxLength)
	}
	return nil
}

func validateEventFieldLengths(event Types.EventResponse) error {
	fields := []struct {
		value string
		name  string
		limit int
	}{
		{event.Name, "event title", maxEventTitleLength},
		{event.ClubName, "organization name", maxEventOrganizationLength},
		{event.Location, "location", maxEventLocationLength},
		{event.Description, "description", maxEventDescriptionLength},
		{event.ContactEmail, "email", maxEventEmailLength},
		{event.WebsiteURL, "website", maxEventWebsiteLength},
	}
	for _, field := range fields {
		if err := validateEventFieldLength(field.value, field.name, field.limit); err != nil {
			return err
		}
	}
	return nil
}

func formatOptionalDateTime(value *time.Time) any {
	if value == nil {
		return nil
	}
	return value.Format(time.RFC3339)
}

func normalizeEventInterests(interests []Types.EventInterest) []Types.EventInterest {
	if interests == nil {
		return []Types.EventInterest{}
	}
	return interests
}

func validateEventInterests(interests []Types.EventInterest) error {
	if len(interests) > maxEventInterests {
		return fmt.Errorf("%w: choose no more than %d interests", ErrInvalidEvent, maxEventInterests)
	}
	seen := make(map[Types.EventInterest]struct{}, len(interests))
	for _, interest := range interests {
		if !Types.IsValidEventInterest(interest) {
			return fmt.Errorf("%w: invalid interest %q", ErrInvalidEvent, interest)
		}
		if _, exists := seen[interest]; exists {
			return fmt.Errorf("%w: duplicate interest %q", ErrInvalidEvent, interest)
		}
		seen[interest] = struct{}{}
	}
	return nil
}

func storeEventInterests(event *Types.EventResponse, serialized string) error {
	if err := json.Unmarshal([]byte(serialized), &event.Interests); err != nil {
		return err
	}
	event.Interests = normalizeEventInterests(event.Interests)
	return nil
}

func validateEventEndDateTime(start, end time.Time) error {
	if end.Sub(start) < 30*time.Minute {
		return fmt.Errorf("%w: end date time must be at least 30 minutes after start date time", ErrInvalidEvent)
	}
	return nil
}

func (repository *Repository) CreateEvent(request Types.CreateEventRequest, hostUserID int) (Types.EventResponse, error) {
	if err := validateEventInterests(request.Interests); err != nil {
		return Types.EventResponse{}, err
	}
	if request.Name == "" {
		return Types.EventResponse{}, fmt.Errorf("%w: name is required", ErrInvalidEvent)
	}
	if request.DateTime.IsZero() {
		return Types.EventResponse{}, fmt.Errorf("%w: date time is required", ErrInvalidEvent)
	}
	if err := requireEventDateTime(request.DateTime); err != nil {
		return Types.EventResponse{}, err
	}
	if request.EndDateTime != nil {
		if err := validateEventEndDateTime(request.DateTime, *request.EndDateTime); err != nil {
			return Types.EventResponse{}, err
		}
		if err := requireEventDateTime(*request.EndDateTime); err != nil {
			return Types.EventResponse{}, err
		}
	}
	if request.Location == "" {
		return Types.EventResponse{}, fmt.Errorf("%w: location is required", ErrInvalidEvent)
	}
	if request.Description == "" {
		return Types.EventResponse{}, fmt.Errorf("%w: description is required", ErrInvalidEvent)
	}
	if err := validateEventFieldLengths(Types.EventResponse{
		Name: request.Name, ClubName: request.ClubName, Location: request.Location,
		Description: request.Description, ContactEmail: request.ContactEmail,
		WebsiteURL: request.WebsiteURL,
	}); err != nil {
		return Types.EventResponse{}, err
	}
	if err := requirePositiveID(hostUserID); err != nil {
		return Types.EventResponse{}, fmt.Errorf("%w: host user id must be positive", ErrInvalidEvent)
	}

	if _, err := repository.GetUser(hostUserID); err != nil {
		return Types.EventResponse{}, err
	}

	interestsJSON, err := json.Marshal(normalizeEventInterests(request.Interests))
	if err != nil {
		return Types.EventResponse{}, err
	}
	result, err := repository.db.Exec(`
		INSERT INTO events (
			event_name, date_time, end_date_time, location, club_name, description,
			contact_email, website, interests, pending, host_user_id
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, TRUE, ?)
	`, request.Name, request.DateTime.Format(time.RFC3339), formatOptionalDateTime(request.EndDateTime),
		request.Location, request.ClubName, request.Description, request.ContactEmail, request.WebsiteURL, string(interestsJSON), hostUserID)

	if err != nil {
		return Types.EventResponse{}, translateEventWriteError(err)
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
	var endDateTime sql.NullString
	var interestsJSON string
	var host Types.UserResponse

	err := repository.db.QueryRow(`
		SELECT e.event_id, e.event_name, e.date_time, e.end_date_time, e.location,
			e.club_name, e.description, e.contact_email, e.website, e.interests,
			e.pending, e.host_user_id,
			u.user_id, u.user_name, u.user_email,
			(SELECT COUNT(*) FROM rsvps r WHERE r.event_id = e.event_id)
		FROM events e
		JOIN users u ON u.user_id = e.host_user_id
		WHERE e.event_id = ?
	`, id).Scan(&event.ID, &event.Name, &dateTime, &endDateTime, &event.Location,
		&event.ClubName, &event.Description, &event.ContactEmail, &event.WebsiteURL, &interestsJSON,
		&event.Pending, &event.HostUserID, &host.ID, &host.Username, &host.Email,
		&event.AttendeeCount)

	if errors.Is(err, sql.ErrNoRows) {
		return Types.EventResponse{}, ErrNotFound
	}

	if err != nil {
		return Types.EventResponse{}, err
	}
	if err := storeEventInterests(&event, interestsJSON); err != nil {
		return Types.EventResponse{}, err
	}

	event.DateTime, err = time.Parse(time.RFC3339, dateTime)
	if err != nil {
		return Types.EventResponse{}, err
	}
	if endDateTime.Valid {
		parsedEndDateTime, parseErr := time.Parse(time.RFC3339, endDateTime.String)
		if parseErr != nil {
			return Types.EventResponse{}, parseErr
		}
		event.EndDateTime = &parsedEndDateTime
	}

	event.Host = &host

	return event, nil
}

func (repository *Repository) GetEvents(filters Types.EventFilters) ([]Types.EventResponse, error) {
	query := `
		SELECT e.event_id, e.event_name, e.date_time, e.end_date_time, e.location,
			e.club_name, e.description, e.contact_email, e.website, e.interests,
			e.pending, e.host_user_id, u.user_id, u.user_name, u.user_email,
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

	if filters.ViewerUserID != nil {
		conditions = append([]string{"(e.pending = FALSE OR e.host_user_id = ?)"}, conditions...)
		args = append([]any{*filters.ViewerUserID}, args...)
	} else {
		conditions = append([]string{"e.pending = FALSE"}, conditions...)
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
		var endDateTime sql.NullString
		var interestsJSON string
		var host Types.UserResponse

		if err := rows.Scan(&event.ID, &event.Name, &dateTime, &endDateTime, &event.Location,
			&event.ClubName, &event.Description, &event.ContactEmail, &event.WebsiteURL, &interestsJSON,
			&event.Pending, &event.HostUserID, &host.ID, &host.Username, &host.Email, &event.AttendeeCount); err != nil {
			return nil, err
		}
		if err := storeEventInterests(&event, interestsJSON); err != nil {
			return nil, err
		}

		event.DateTime, err = time.Parse(time.RFC3339, dateTime)
		if err != nil {
			return nil, err
		}
		if endDateTime.Valid {
			parsedEndDateTime, parseErr := time.Parse(time.RFC3339, endDateTime.String)
			if parseErr != nil {
				return nil, parseErr
			}
			event.EndDateTime = &parsedEndDateTime
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
	edited := request.Name != nil || request.DateTime != nil || request.EndDateTime != nil ||
		request.Location != nil || request.ClubName != nil || request.Description != nil ||
		request.ContactEmail != nil || request.WebsiteURL != nil || request.Interests != nil

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
		if err := requireEventDateTime(*request.DateTime); err != nil {
			return Types.EventResponse{}, err
		}
		current.DateTime = *request.DateTime
	}
	if request.EndDateTime != nil && request.EndDateTime.Set {
		if request.EndDateTime.Value == nil {
			current.EndDateTime = nil
		} else {
			requestedEndDateTime := request.EndDateTime.Value
			if err := validateEventEndDateTime(current.DateTime, *requestedEndDateTime); err != nil {
				return Types.EventResponse{}, err
			}
			if err := requireEventDateTime(*requestedEndDateTime); err != nil {
				return Types.EventResponse{}, err
			}
			current.EndDateTime = requestedEndDateTime
		}
	}
	if current.EndDateTime != nil {
		if err := validateEventEndDateTime(current.DateTime, *current.EndDateTime); err != nil {
			return Types.EventResponse{}, err
		}
	}

	if request.Location != nil {
		if *request.Location == "" {
			return Types.EventResponse{}, fmt.Errorf("%w: location cannot be empty", ErrInvalidEvent)
		}
		current.Location = *request.Location
	}

	if request.ClubName != nil {
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
	if request.Interests != nil {
		current.Interests = normalizeEventInterests(*request.Interests)
	}
	if err := validateEventInterests(current.Interests); err != nil {
		return Types.EventResponse{}, err
	}

	if err := validateEventFieldLengths(current); err != nil {
		return Types.EventResponse{}, err
	}
	interestsJSON, err := json.Marshal(normalizeEventInterests(current.Interests))
	if err != nil {
		return Types.EventResponse{}, err
	}

	_, err = repository.db.Exec(`
		UPDATE events
		SET event_name = ?, date_time = ?, end_date_time = ?, location = ?, club_name = ?, description = ?,
			contact_email = ?, website = ?, interests = ?, pending = ?, host_user_id = ?
		WHERE event_id = ?
	`, current.Name, current.DateTime.Format(time.RFC3339), formatOptionalDateTime(current.EndDateTime), current.Location,
		current.ClubName, current.Description, current.ContactEmail, current.WebsiteURL, string(interestsJSON),
		edited, current.HostUserID, id)

	if err != nil {
		return Types.EventResponse{}, translateEventWriteError(err)
	}

	return repository.GetEvent(id)
}

func (repository *Repository) GetPendingEvents() ([]Types.EventResponse, error) {
	rows, err := repository.db.Query(`
		SELECT e.event_id, e.event_name, e.date_time, e.end_date_time, e.location,
			e.club_name, e.description, e.contact_email, e.website, e.interests,
			e.pending, e.host_user_id, u.user_id, u.user_name, u.user_email,
			(SELECT COUNT(*) FROM rsvps r WHERE r.event_id = e.event_id)
		FROM events e
		JOIN users u ON u.user_id = e.host_user_id
		WHERE e.pending = TRUE
		ORDER BY e.date_time, e.event_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]Types.EventResponse, 0)
	for rows.Next() {
		var event Types.EventResponse
		var dateTime string
		var endDateTime sql.NullString
		var interestsJSON string
		var host Types.UserResponse
		if err := rows.Scan(&event.ID, &event.Name, &dateTime, &endDateTime, &event.Location,
			&event.ClubName, &event.Description, &event.ContactEmail, &event.WebsiteURL, &interestsJSON,
			&event.Pending, &event.HostUserID, &host.ID, &host.Username, &host.Email,
			&event.AttendeeCount); err != nil {
			return nil, err
		}
		if err := storeEventInterests(&event, interestsJSON); err != nil {
			return nil, err
		}
		event.DateTime, err = time.Parse(time.RFC3339, dateTime)
		if err != nil {
			return nil, err
		}
		if endDateTime.Valid {
			parsed, parseErr := time.Parse(time.RFC3339, endDateTime.String)
			if parseErr != nil {
				return nil, parseErr
			}
			event.EndDateTime = &parsed
		}
		event.Host = &host
		events = append(events, event)
	}
	return events, rows.Err()
}

func (repository *Repository) ApproveEvent(id int) error {
	if err := requirePositiveID(id); err != nil {
		return err
	}
	result, err := repository.db.Exec(`UPDATE events SET pending = FALSE WHERE event_id = ? AND pending = TRUE`, id)
	if err != nil {
		return err
	}
	return requireAffected(result)
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
