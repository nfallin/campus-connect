package Model

import (
	"testing"
	"time"
)

// Test that CreateUser creates a user, assigns an ID, initializes its RSVP list, and adds it to campus
func TestCreateUser(t *testing.T) {
	c := NewCampus()
	alice := c.CreateUser("Alice")

	if alice.name != "Alice" {
		t.Errorf("expected Alice's name to be 'Alice', got '%s'", alice.name)
	}

	if alice.id != 1 {
		t.Errorf("expected Alice's ID to be 1, got %d", alice.id)
	}

	if len(alice.rsvpList) != 0 {
		t.Errorf("expected a new user's RSVP list to be empty, got %v", alice.rsvpList)
	}

	if c.GetUser(alice.id) != alice {
		t.Errorf("expected Alice to be added to campus users, but she was not found")
	}
}

// Test that CreateEvent creates an event, assigns an ID, initializes attendees, stores its fields, and adds it to campus
func TestCreateEvent(t *testing.T) {
	c := NewCampus()
	host := c.CreateUser("Host")
	event := c.CreateEvent("Meeting", time.Now(), "Room 101", "Club A", host, "host@example.com", "https://example.com/events/meeting")

	if event.name != "Meeting" {
		t.Errorf("expected event name to be 'Meeting', got '%s'", event.name)
	}

	if event.id != 1 {
		t.Errorf("expected event ID to be 1, got %d", event.id)
	}

	if len(event.attendees) != 0 {
		t.Errorf("expected a new event's attendees list to be empty, got %v", event.attendees)
	}

	if event.host != host {
		t.Errorf("expected event host to be %v, got %v", host, event.host)
	}

	if event.contactEmail != "host@example.com" {
		t.Errorf("expected event contact email to be 'host@example.com', got '%s'", event.contactEmail)
	}

	if event.websiteURL != "https://example.com/events/meeting" {
		t.Errorf("expected event website URL to be 'https://example.com/events/meeting', got '%s'", event.websiteURL)
	}

	if c.GetEvent(event.id) != event {
		t.Errorf("expected event to be added to campus events, but it was not found")
	}
}

// Test that a new campus initializes correctly with empty maps and starting IDs
func TestNewCampus(t *testing.T) {
	c := NewCampus()

	if len(c.users) != 0 {
		t.Errorf("expected empty users map, got %v", c.users)
	}

	if len(c.events) != 0 {
		t.Errorf("expected empty events map, got %v", c.events)
	}

	if c.nextUserID != 1 {
		t.Errorf("expected nextUserID to start at 1, got %d", c.nextUserID)
	}

	if c.nextEventID != 1 {
		t.Errorf("expected nextEventID to start at 1, got %d", c.nextEventID)
	}
}

// Test that GetNextUserID generates sequential user IDs
func testGetNextUserID(t *testing.T) {
	c := NewCampus()
	alice := c.CreateUser("Alice")
	if alice.id != 1 {
		t.Errorf("expected Alice's ID to be 1, got %d", alice.id)
	}

	bob := c.CreateUser("Bob")
	if bob.id != 2 {
		t.Errorf("expected Bob's ID to be 2, got %d", bob.id)
	}
}

// Test that GetNextEventID generates sequential event IDs
func testGetNextEventID(t *testing.T) {
	c := NewCampus()
	event1 := c.CreateEvent("Event 1", time.Now(), "Location 1", "Club A", c.CreateUser("Host 1"), "host1@example.com", "https://example.com/events/1")

	if event1.id != 1 {
		t.Errorf("expected Event 1's ID to be 1, got %d", event1.id)
	}

	event2 := c.CreateEvent("Event 2", time.Now(), "Location 2", "Club B", c.CreateUser("Host 2"), "host2@example.com", "https://example.com/events/2")
	if event2.id != 2 {
		t.Errorf("expected Event 2's ID to be 2, got %d", event2.id)
	}
}

// Test that GetUsers returns all users in the campus
func TestGetUsers(t *testing.T) {
	c := NewCampus()
	alice := c.CreateUser("Alice")
	bob := c.CreateUser("Bob")

	users := c.GetUsers()

	if len(users) != 2 || users[alice.id] != alice || users[bob.id] != bob {
		t.Errorf("expected users map to contain Alice and Bob, got %v", users)
	}
}

// Test that GetUser Returns an existing user and nil for an unknown ID
func TestGetUser(t *testing.T) {
	c := NewCampus()
	alice := c.CreateUser("Alice")

	if c.GetUser(alice.id) != alice {
		t.Errorf("expected to retrieve Alice, got %v", c.GetUser(alice.id))
	}

	if c.GetUser(999) != nil {
		t.Errorf("expected to retrieve nil for unknown user ID, got %v", c.GetUser(999))
	}
}

func TestGetUserByUsername(t *testing.T) {
	c := NewCampus()
	alice := c.CreateUser("Alice")

	if c.GetUserByName("Alice") != alice {
		t.Errorf("expected to retrieve Alice by username, got %v", c.GetUserByName("Alice"))
	}

	if c.GetUserByName("Unknown") != nil {
		t.Errorf("expected nil for an unknown username, got %v", c.GetUserByName("Unknown"))
	}
}

func TestCreateUserRejectsDuplicateUsername(t *testing.T) {
	c := NewCampus()
	alice := c.CreateUser("Alice")
	duplicate := c.CreateUser("Alice")

	if duplicate != nil {
		t.Errorf("expected duplicate username to be rejected, got %v", duplicate)
	}

	if len(c.users) != 1 {
		t.Errorf("expected one user after duplicate creation, got %d", len(c.users))
	}

	if c.nextUserID != 2 {
		t.Errorf("expected duplicate creation not to consume a user ID, next ID is %d", c.nextUserID)
	}

	if c.GetUserByName("Alice") != alice {
		t.Errorf("expected original user to remain associated with username")
	}
}

// Test that GetEvents returns all events in the campus
func TestGetEvents(t *testing.T) {
	c := NewCampus()
	event1 := c.CreateEvent("Event 1", time.Now(), "Location 1", "Club A", c.CreateUser("Host 1"), "host1@example.com", "https://example.com/events/1")
	event2 := c.CreateEvent("Event 2", time.Now(), "Location 2", "Club B", c.CreateUser("Host 2"), "host2@example.com", "https://example.com/events/2")

	events := c.GetEvents()

	if len(events) != 2 || events[event1.id] != event1 || events[event2.id] != event2 {
		t.Errorf("expected events map to contain Event 1 and Event 2, got %v", events)
	}
}

func TestFilterEventsByDateRange(t *testing.T) {
	c := NewCampus()
	start := time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC)
	inRange := c.CreateEvent("In range", start.AddDate(0, 0, 1), "Room 101", "Club A", c.CreateUser("Host 1"), "host1@example.com", "https://example.com/events/in-range")
	c.CreateEvent("Before range", start.AddDate(0, 0, -1), "Room 102", "Club A", c.CreateUser("Host 2"), "host2@example.com", "https://example.com/events/before")
	c.CreateEvent("After range", start.AddDate(0, 0, 3), "Room 103", "Club A", c.CreateUser("Host 3"), "host3@example.com", "https://example.com/events/after")
	end := start.AddDate(0, 0, 2)

	events := c.FilterEvents(&start, &end, "")

	if len(events) != 1 || events[0] != inRange {
		t.Errorf("expected only the event in the date range, got %v", events)
	}
}

func TestFilterEventsByClubName(t *testing.T) {
	c := NewCampus()
	clubEvent := c.CreateEvent("Club event", time.Now(), "Room 101", "Club A", c.CreateUser("Host 1"), "host1@example.com", "https://example.com/events/club")
	c.CreateEvent("Other event", time.Now(), "Room 102", "Club B", c.CreateUser("Host 2"), "host2@example.com", "https://example.com/events/other")

	events := c.FilterEvents(nil, nil, "Club A")

	if len(events) != 1 || events[0] != clubEvent {
		t.Errorf("expected only Club A's event, got %v", events)
	}
}

// Test that GetEvent Returns an existing event and nil for an unknown ID
func TestGetEvent(t *testing.T) {
	c := NewCampus()
	event1 := c.CreateEvent("Event 1", time.Now(), "Location 1", "Club A", c.CreateUser("Host"), "host@example.com", "https://example.com/events/1")

	if c.GetEvent(event1.id) != event1 {
		t.Errorf("expected to retrieve Event 1, got %v", c.GetEvent(event1.id))
	}

	if c.GetEvent(999) != nil {
		t.Errorf("expected to retrieve nil for unknown event ID, got %v", c.GetEvent(999))
	}
}

// Test that CreateMultipleEntities assigns independent sequential IDs for users and events
func TestCreateMultipleEntities(t *testing.T) {
	c := NewCampus()
	alice := c.CreateUser("Alice")
	bob := c.CreateUser("Bob")

	if alice.id != 1 || bob.id != 2 {
		t.Errorf("expected Alice's ID to be 1 and Bob's ID to be 2, got %d and %d", alice.id, bob.id)
	}

	event1 := c.CreateEvent("Event 1", time.Now(), "Location 1", "Club A", alice, "alice@example.com", "https://example.com/events/1")
	event2 := c.CreateEvent("Event 2", time.Now(), "Location 2", "Club B", bob, "bob@example.com", "https://example.com/events/2")

	if event1.id != 1 || event2.id != 2 {
		t.Errorf("expected Event 1's ID to be 1 and Event 2's ID to be 2, got %d and %d", event1.id, event2.id)
	}
}

// Test that RemoveEvent removes an event from campus and cleans up attending users
func TestRemoveEvent(t *testing.T) {
	c := NewCampus()
	alice := c.CreateUser("Alice")
	event := c.CreateEvent("Meeting", time.Now(), "Room 101", "Club A", alice, "alice@example.com", "https://example.com/events/meeting")
	alice.RsvpAdd(event)

	c.RemoveEvent(event)

	if c.GetEvent(event.id) != nil {
		t.Errorf("expected event to be removed from campus, but it was still found")
	}

	if _, exists := alice.rsvpList[event.id]; exists {
		t.Errorf("expected event to be removed from Alice's RSVP list, but it was still found")
	}
}

// Test that RemoveUser removes a user and cleans up the user's RSVPs from events
func TestRemoveUser(t *testing.T) {
	c := NewCampus()
	alice := c.CreateUser("Alice")
	event := c.CreateEvent("Meeting", time.Now(), "Room 101", "Club A", alice, "alice@example.com", "https://example.com/events/meeting")
	alice.RsvpAdd(event)

	c.RemoveUser(alice)

	if c.GetUser(alice.id) != nil {
		t.Errorf("expected user to be removed from campus, but they were still found")
	}

	if c.GetUserByName("Alice") != nil {
		t.Errorf("expected removed user's username to be available, but it was still found")
	}

	if _, exists := event.attendees[alice.id]; exists {
		t.Errorf("expected Alice to be removed from event's attendees list, but she was still found")
	}
}
