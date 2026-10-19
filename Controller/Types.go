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
	Name     string    `json:"name" binding:"required"`
	DateTime time.Time `json:"dateTime" binding:"required"`
	Location string    `json:"location" binding:"required"`
	ClubName string    `json:"clubName" binding:"required"`
}

type EventResponse struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	DateTime      time.Time `json:"dateTime"`
	Location      string    `json:"location"`
	ClubName      string    `json:"clubName"`
	AttendeeCount int       `json:"attendeeCount"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}