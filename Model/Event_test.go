package Model

import (
	"testing"
	"time"
)

// Test that GetName() eturns the event name
func TestGetName(t *testing.T) {
	c := NewCampus()
	event := c.CreateEvent("Meeting", time.Now(), "Room 101", "Club A")

	if name := event.GetName(); name != "Meeting" {
		t.Errorf("expected event name to be 'Meeting', got '%s'", name)
	}
}

// Test that GetID() returns the event ID
func TestGetID(t *testing.T) {
	c := NewCampus()
	event := c.CreateEvent("Meeting", time.Now(), "Room 101", "Club A")

	if id := event.GetID(); id != 1 {
		t.Errorf("expected event ID to be 1, got %d", id)
	}
}

// Test that GetDateTime() returns the correct event date/time
func TestGetDateTime(t *testing.T) {
	c := NewCampus()
	now := time.Now()
	event := c.CreateEvent("Meeting", now, "Room 101", "Club A")

	if dateTime := event.GetDateTime(); !dateTime.Equal(now) {
		t.Errorf("expected event date/time to be %v, got %v", now, dateTime)
	}
}

// Test that GetLocation() returns the correct event location
func TestGetLocation(t *testing.T) {
	c := NewCampus()
	event := c.CreateEvent("Meeting", time.Now(), "Room 101", "Club A")

	if location := event.GetLocation(); location != "Room 101" {
		t.Errorf("expected event location to be 'Room 101', got '%s'", location)
	}
}

// Test that GetAttendees() returns the event's attendees
func TestGetAttendees(t *testing.T) {
	c := NewCampus()
	event := c.CreateEvent("Meeting", time.Now(), "Room 101", "Club A")
	alice := c.CreateUser("Alice")
	bob := c.CreateUser("Bob")

	alice.RsvpAdd(event)
	bob.RsvpAdd(event)

	if attendees := event.GetAttendees(); len(attendees) != 2 {
		t.Errorf("expected 2 attendees, got %d", len(attendees))
	}
}

// Test that GetNumAttendees() returns 0 when nobody is attending
func TestGetNumAttendeesEmpty(t *testing.T) {
	c := NewCampus()
	event := c.CreateEvent("Meeting", time.Now(), "Room 101", "Club A")

	if numAttendees := event.GetNumAttendees(); numAttendees != 0 {
		t.Errorf("expected 0 attendees, got %d", numAttendees)
	}
}

// Test that GetNumAttendees() returns the correct number of attendees
func TestGetNumAttendees(t *testing.T) {
	c := NewCampus()
	event := c.CreateEvent("Meeting", time.Now(), "Room 101", "Club A")
	alice := c.CreateUser("Alice")
	bob := c.CreateUser("Bob")

	alice.RsvpAdd(event)
	bob.RsvpAdd(event)

	if numAttendees := event.GetNumAttendees(); numAttendees != 2 {
		t.Errorf("expected 2 attendees, got %d", numAttendees)
	}
}

// Test that GetClubName() returns the name of the club hosting the event
func TestGetClubName(t *testing.T) {
	c := NewCampus()
	event := c.CreateEvent("Meeting", time.Now(), "Room 101", "Club A")

	if clubName := event.GetClubName(); clubName != "Club A" {
		t.Errorf("expected club name to be 'Club A', got '%s'", clubName)
	}
}
