package Model

import (
	"time"
)

// campus struct definition and constructor function
type Campus struct {
	users       map[int]*User
	events      map[int]*Event
	nextEventID int
	nextUserID  int
}

func NewCampus() *Campus {
	return &Campus{
		users:       make(map[int]*User),
		events:      make(map[int]*Event),
		nextEventID: 1,
		nextUserID:  1,
	}
}

// ---

// Getters for the 'events', ”, and 'users' maps in the campus as well as individual items.
func (c Campus) GetEvents() map[int]*Event {
	return c.events
}

func (c Campus) GetEvent(eventID int) *Event {
	return c.events[eventID]
}

func (c Campus) GetUsers() map[int]*User {
	return c.users
}

func (c Campus) GetUser(userID int) *User {
	return c.users[userID]
}

func (c Campus) GetUserByName(name string) *User {
	for _, user := range c.users {
		if user.name == name {
			return user
		}
	}

	return nil
}

// ---

// incrementally generate unique IDs for events and users
func (c *Campus) getNextEventID() int {
	id := c.nextEventID
	c.nextEventID++
	return id
}

func (c *Campus) getNextUserID() int {
	id := c.nextUserID
	c.nextUserID++
	return id
}

// ---

// Constructor methods for creating new events, and users.
func (c *Campus) CreateEvent(name string, dateTime time.Time, location string, clubName string) *Event {
	newEvent := &Event{
		name:      name,
		id:        c.getNextEventID(),
		dateTime:  dateTime,
		location:  location,
		attendees: map[int]struct{}{},
		clubName:  clubName,
	}

	c.events[newEvent.id] = newEvent
	return newEvent
}

func (c *Campus) CreateUser(name string) *User {
	if c.GetUserByName(name) != nil {
		return nil
	}

	newUser := &User{
		name:     name,
		id:       c.getNextUserID(),
		rsvpList: make(map[int]struct{}),
	}

	c.users[newUser.id] = newUser
	return newUser
}

// ---

// Remove methods for events, and users. These methods ensure that all references to the removed entity are cleaned up from other entities
func (c *Campus) RemoveEvent(event *Event) {
	// clean up all users that are attending this event before removing the event
	for userID := range event.attendees {
		if user := c.users[userID]; user != nil {
			user.RsvpRemove(event)
		}
	}
	delete(c.events, event.id)
}

func (c *Campus) RemoveUser(user *User) {
	// clean up all events that the user has RSVP'd to before removing the user
	for eventID := range user.rsvpList {
		if event := c.events[eventID]; event != nil {
			user.RsvpRemove(event)
		}
	}
	delete(c.users, user.id)
}

// ---
