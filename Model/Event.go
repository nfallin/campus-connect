package Model

import (
	"time"
)

type Event struct {
	name      string
	id        int
	dateTime  time.Time
	location  string
	attendees map[int]struct{}
	clubName  string
}

func (e *Event) GetName() string {
	return e.name
}

func (e *Event) GetID() int {
	return e.id
}

func (e *Event) GetDateTime() time.Time {
	return e.dateTime
}

func (e *Event) GetLocation() string {
	return e.location
}

func (e *Event) GetAttendees() map[int]struct{} {
	return e.attendees
}

func (e *Event) GetNumAttendees() int {
	return len(e.attendees)
}

func (e *Event) GetClubName() string {
	return e.clubName
}
