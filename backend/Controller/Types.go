package Controller

import (
	"time"
)

type CreateUserRequest struct {
	Name string `json:"name" binding:"required"`
}

type UserResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CreateEventRequest struct {
	Name         string    `json:"name" binding:"required"`
	Description  string    `json:"description" binding:"required"`
	DateTime     time.Time `json:"dateTime" binding:"required"`
	Location     string    `json:"location" binding:"required"`
	ClubName     string    `json:"clubName" binding:"required"`
	HostUserID   *int      `json:"hostUserId" binding:"required"`
	ContactEmail string    `json:"contactEmail"`
	WebsiteURL   string    `json:"websiteURL"`
}

type UpdateEventRequest struct {
	Name         *string    `json:"name"`
	Description  *string    `json:"description"`
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
	Description   string        `json:"description"`
	DateTime      time.Time     `json:"dateTime"`
	Location      string        `json:"location"`
	ClubName      string        `json:"clubName"`
	Host          *UserResponse `json:"host,omitempty"`
	ContactEmail  string        `json:"contactEmail"`
	WebsiteURL    string        `json:"websiteURL"`
	AttendeeCount int           `json:"attendeeCount"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
