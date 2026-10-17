package Model

import (
	"testing"
	"time"
)

// Test that GetRsvpList returns the correct RSVP list for a user
func TestGetRsvpList(t *testing.T) {
	c := NewCampus()
	alice := c.CreateUser("Alice")

	if rsvpList := alice.GetRsvpList(); len(rsvpList) != 0 {
		t.Errorf("expected a new user's RSVP list to be empty, got %v", rsvpList)
	}
}

// Test that RsvpAdd adds an event to the user's RSVP list and that RsvpRemove removes it correctly
func TestRsvpAddRemove(t *testing.T) {
	c := NewCampus()
	alice := c.CreateUser("Alice")
	event := c.CreateEvent("Meeting", time.Now(), "Room 101", "Club A")

	alice.RsvpAdd(event)
	if _, exists := alice.rsvpList[event.id]; !exists {
		t.Errorf("expected event ID %d to be in Alice's RSVP list, but it was not found", event.id)
	}

	alice.RsvpRemove(event)
	if _, exists := alice.rsvpList[event.id]; exists {
		t.Errorf("expected event ID %d to be removed from Alice's RSVP list, but it was still found", event.id)
	}
}

// Test that RSVPing to the same event twice doesn't create duplicate entries
func TestRsvpSameEventTwice(t *testing.T) {
	c := NewCampus()
	alice := c.CreateUser("Alice")
	event := c.CreateEvent("Meeting", time.Now(), "Room 101", "Club A")
	alice.RsvpAdd(event)
	alice.RsvpAdd(event) // Attempt to RSVP again

	if len(alice.rsvpList) != 1 {
		t.Errorf("expected Alice's RSVP list to contain only one entry for the event, but it contained %d entries", len(alice.rsvpList))
	}
}

// Test that removing an event the user hasn't RSVP'd to doesn't affect the RSVP map
func TestRsvpRemoveNonRsvpEvent(t *testing.T) {
	c := NewCampus()
	alice := c.CreateUser("Alice")
	event := c.CreateEvent("Meeting", time.Now(), "Room 101", "Club A")
	alice.RsvpRemove(event) // Attempt to remove an event that Alice hasn't RSVP'd to

	if len(alice.rsvpList) != 0 {
		t.Errorf("expected Alice's RSVP list to remain empty after attempting to remove a non-RSVP'd event, but it contained %d entries", len(alice.rsvpList))
	}
}

// Test that a user can RSVP to multiple events independently
func TestRsvpMultipleEvents(t *testing.T) {
	c := NewCampus()
	alice := c.CreateUser("Alice")
	event1 := c.CreateEvent("Meeting 1", time.Now(), "Room 101", "Club A")
	event2 := c.CreateEvent("Meeting 2", time.Now(), "Room 102", "Club B")

	alice.RsvpAdd(event1)
	alice.RsvpAdd(event2)

	if len(alice.rsvpList) != 2 {
		t.Errorf("expected Alice's RSVP list to contain 2 entries, but it contained %d entries", len(alice.rsvpList))
	}
}
