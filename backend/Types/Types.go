package Types

import (
	"bytes"
	"encoding/json"
	"time"
)

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type UpdateUserRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateEventRequest struct {
	Name         string     `json:"name" binding:"required"`
	DateTime     time.Time  `json:"dateTime" binding:"required"`
	EndDateTime  *time.Time `json:"endDateTime"`
	Location     string     `json:"location" binding:"required"`
	ClubName     string     `json:"clubName" binding:"required"`
	Description  string     `json:"description" binding:"required"`
	HostUserID   *int       `json:"hostUserId" binding:"required"`
	ContactEmail string     `json:"contactEmail"`
	WebsiteURL   string     `json:"websiteURL"`
}

type UpdateEventRequest struct {
	Name         *string           `json:"name"`
	DateTime     *time.Time        `json:"dateTime"`
	EndDateTime  *OptionalDateTime `json:"endDateTime,omitempty"`
	Location     *string           `json:"location"`
	ClubName     *string           `json:"clubName"`
	Description  *string           `json:"description"`
	HostUserID   *int              `json:"hostUserId"`
	ContactEmail *string           `json:"contactEmail"`
	WebsiteURL   *string           `json:"websiteURL"`
}

type OptionalDateTime struct {
	Value *time.Time
	Set   bool
}

func (value *OptionalDateTime) UnmarshalJSON(data []byte) error {
	value.Set = true
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		value.Value = nil
		return nil
	}
	if string(bytes.TrimSpace(data)) == `""` {
		value.Value = nil
		return nil
	}

	var dateTime time.Time
	if err := json.Unmarshal(data, &dateTime); err != nil {
		return err
	}
	value.Value = &dateTime
	return nil
}

type EventResponse struct {
	ID            int           `json:"id"`
	Name          string        `json:"name"`
	DateTime      time.Time     `json:"dateTime"`
	EndDateTime   *time.Time    `json:"endDateTime,omitempty"`
	Location      string        `json:"location"`
	ClubName      string        `json:"clubName"`
	Description   string        `json:"description"`
	Host          *UserResponse `json:"host,omitempty"`
	ContactEmail  string        `json:"contactEmail"`
	WebsiteURL    string        `json:"websiteURL"`
	AttendeeCount int           `json:"attendeeCount"`
	HostUserID    int           `json:"hostUserId"`
}

type EventFilters struct {
	StartDate  *time.Time
	EndDate    *time.Time
	ClubName   string
	Location   string
	HostUserID *int
}

type ErrorResponse struct {
	Message string `json:"message"`
}
