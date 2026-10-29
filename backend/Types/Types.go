package Types

import (
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
	Name         string    `json:"name" binding:"required"`
	DateTime     time.Time `json:"dateTime" binding:"required"`
	Location     string    `json:"location" binding:"required"`
	ClubName     string    `json:"clubName" binding:"required"`
	HostUserID   *int      `json:"hostUserId" binding:"required"`
	ContactEmail string    `json:"contactEmail"`
	WebsiteURL   string    `json:"websiteURL"`
}

type UpdateEventRequest struct {
	Name         *string    `json:"name"`
	DateTime     *time.Time `json:"dateTime"`
	Location     *string    `json:"location"`
	ClubName     *string    `json:"clubName"`
	HostUserID   *int       `json:"hostUserId"`
	ContactEmail *string    `json:"contactEmail"`
	WebsiteURL   *string    `json:"websiteURL"`
}

type EventResponse struct {
	ID            int           `json:"id"`
	Name          string        `json:"name"`
	DateTime      time.Time     `json:"dateTime"`
	Location      string        `json:"location"`
	ClubName      string        `json:"clubName"`
	Host          *UserResponse `json:"host,omitempty"`
	ContactEmail  string        `json:"contactEmail"`
	WebsiteURL    string        `json:"websiteURL"`
	AttendeeCount int           `json:"attendeeCount"`
	HostUserID    int           `json:"-"`
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
