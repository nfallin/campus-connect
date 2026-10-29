package Model

import (
	"campus_connect/Types"
	"time"
)

// RSVPs
func (repository *Repository) AddRSVP(userID, eventID int) error {
	if err := requirePositiveID(userID); err != nil {
		return err
	}
	if err := requirePositiveID(eventID); err != nil {
		return err
	}
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
	if err := requirePositiveID(userID); err != nil {
		return err
	}
	if err := requirePositiveID(eventID); err != nil {
		return err
	}
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
	if err := requirePositiveID(userID); err != nil {
		return nil, err
	}
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
