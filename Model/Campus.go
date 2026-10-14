package Model

import (
	"image"
	"time"
)

// campus struct definition and constructor function
type Campus struct {
	clubs       map[int]*Club
	users       map[int]*User
	events		map[int]*Event
	nextEventID int
	nextClubID  int
	nextUserID  int
}

func NewCampus() *Campus {
	return &Campus{
		clubs:       make(map[int]*Club),
		users:       make(map[int]*User),
		events:		make(map[int]*Event),
		nextEventID: 1,
		nextClubID:  1,
		nextUserID:  1,
	}
}
// ---

// Getters for the 'events', 'clubs', and 'users' maps in the campus as well as individual items.
func (c Campus) GetEvents() map[int]*Event {
	return c.events
}

func (c Campus) GetEvent(eventID int) *Event {
	return c.events[eventID]
}

func (c Campus) GetClubs() map[int]*Club {
	return c.clubs
}

func (c Campus) GetClub(clubID int) *Club {
	return c.clubs[clubID]
}

func (c Campus) GetUsers() map[int]*User {
	return c.users
}

func (c Campus) GetUser(userID int) *User {
	return c.users[userID]
}
// ---

// incrementally generate unique IDs for events, clubs, and users
func (c *Campus) getNextEventID() int {
	id := c.nextEventID
	c.nextEventID++
	return id
}

func (c *Campus) getNextClubID() int {
	id := c.nextClubID
	c.nextClubID++
	return id
}

func (c *Campus) getNextUserID() int {
	id := c.nextUserID
	c.nextUserID++
	return id
}
// ---

// Constructor methods for creating new events, clubs, and users.
func (c *Campus) CreateEvent(name string, dateTime time.Time, location string, poster image.Image, clubID int) *Event {
	newEvent := &Event{
		name:      name,
		id:        c.getNextEventID(),
		dateTime:  dateTime,
		location:  location,
		poster:    poster,
		attendees: map[int]struct{}{},
		clubID:    clubID,
	}

	c.addEvent(newEvent)
	return newEvent
}

func (c *Campus) addEvent(event *Event) {
	c.events[event.id] = event
}

func (c *Campus) CreateUser(name string) *User {
	newUser := &User{
		name:     name,
		id:       c.getNextUserID(),
		rsvpList: make(map[int]struct{}),
	}

	c.addUser(newUser)
	return newUser
}

func (c *Campus) addUser(user *User) {
	c.users[user.id] = user
}

func (c *Campus) CreateClub(name string) *Club {
	newClub := &Club{
		name:   name,
		id:     c.getNextClubID(),
	}

	c.addClub(newClub)
	return newClub
}

func (c *Campus) addClub(club *Club) {
	c.clubs[club.id] = club
}
// ---

// Remove methods for events, clubs, and users. These methods ensure that all references to the removed entity are cleaned up from other entities
func (c *Campus) RemoveEvent(event *Event) {
	// clean up all users that are attending this event before removing the event
	for userID := range event.attendees {
		if user := c.users[userID]; user != nil {
			user.RsvpRemove(event)
		}
	}
	delete(c.events, event.id)
}
	
func (c *Campus) RemoveClub(club *Club) {
	// delete up all events that belong to this club before removing the club
	for _, event := range c.events {
		if event.clubID == club.id {
			c.RemoveEvent(event)
		}
	}
	delete(c.clubs, club.id)
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