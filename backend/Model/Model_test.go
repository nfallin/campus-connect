package Model

import (
	"campus_connect/Types"
	"database/sql"
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestRepository initializes a new Repository with a temporary database for testing purposes.
func testRepository(t *testing.T) *Repository {
	t.Helper()
	db, err := Init(t.TempDir() + "\\test.db")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return NewRepository(db)
}

func TestInitMigratesExistingDatabaseForOptionalEndTime(t *testing.T) {
	path := filepath.Join(t.TempDir(), "existing.db")
	db, err := open(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(usersTable); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(`
		CREATE TABLE events (
			event_id INTEGER PRIMARY KEY AUTOINCREMENT,
			event_name TEXT NOT NULL,
			date_time TEXT NOT NULL,
			location TEXT NOT NULL,
			club_name TEXT NOT NULL,
			description TEXT NOT NULL,
			contact_email TEXT,
			website TEXT,
			host_user_id INTEGER NOT NULL,
			UNIQUE(event_name),
			FOREIGN KEY(host_user_id) REFERENCES users(user_id) ON DELETE RESTRICT
		)
	`); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	migratedDB, err := Init(path)
	if err != nil {
		t.Fatal(err)
	}
	defer migratedDB.Close()

	var endDateTime sql.NullString
	err = migratedDB.QueryRow(`SELECT end_date_time FROM events LIMIT 1`).Scan(&endDateTime)
	if err != nil && err != sql.ErrNoRows {
		t.Fatalf("read migrated end time column: %v", err)
	}
}

func TestRepositoryUsers(t *testing.T) {
	repository := testRepository(t)
	request := Types.CreateUserRequest{Name: "Alice", Email: "alice@example.com"}

	// attempt to create a valid user
	user, created, err := repository.CreateUser(request)
	if err != nil || !created || user.ID == 0 {
		t.Fatalf("create user = %+v, %v, %v", user, created, err)
	}

	// attempt to create a duplicate user
	existing, created, err := repository.CreateUser(request)
	if err != nil || created || existing.ID != user.ID {
		t.Fatalf("duplicate user = %+v, %v, %v", existing, created, err)
	}

	// attempt to get the user
	got, err := repository.GetUser(user.ID)
	if err != nil || got != user {
		t.Fatalf("get user = %+v, %v", got, err)
	}

	// attempt to get all users
	users, err := repository.GetAllUsers()
	if err != nil || len(users) != 1 {
		t.Fatalf("get all users = %+v, %v", users, err)
	}

	name := "Updated Alice"
	email := "updated@example.com"
	// attempt to update a user with new information
	updated, err := repository.UpdateUser(user.ID, Types.UpdateUserRequest{Name: &name, Email: &email})
	if err != nil || updated.Name != name || updated.Email != email {
		t.Fatalf("update user = %+v, %v", updated, err)
	}

	// attempt to delete the user
	if err := repository.DeleteUser(user.ID); err != nil {
		t.Fatal(err)
	}
	// confirm the user was deleted
	if _, err := repository.GetUser(user.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("get deleted user error = %v", err)
	}
	// attempt to delete a user that doesn't exist
	if err := repository.DeleteUser(user.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("delete missing user error = %v", err)
	}

	// attempt to update a user to match the name and email of an existing user
	duplicateA, _, err := repository.CreateUser(Types.CreateUserRequest{Name: "A", Email: "a@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	duplicateB, _, err := repository.CreateUser(Types.CreateUserRequest{Name: "B", Email: "b@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repository.UpdateUser(duplicateB.ID, Types.UpdateUserRequest{Name: stringPointer("A"), Email: stringPointer(duplicateA.Email)}); err == nil {
		t.Fatal("expected duplicate user update error")
	}
}

func TestRepositoryEventsAndFilters(t *testing.T) {
	repository := testRepository(t)
	host, _, err := repository.CreateUser(Types.CreateUserRequest{Name: "Host", Email: "host@example.com"})
	if err != nil {
		t.Fatal(err)
	}

	// attempt to create a valid event
	when := time.Now().UTC().Add(24 * time.Hour)
	endWhen := when.Add(time.Hour)
	event, err := repository.CreateEvent(Types.CreateEventRequest{
		Name: "Meeting", DateTime: when, Location: "Room 1", ClubName: "Club A", Description: "Description",
		HostUserID: &host.ID, EndDateTime: &endWhen,
	})
	if err != nil || event.Host == nil || event.AttendeeCount != 0 ||
		event.EndDateTime == nil || !event.EndDateTime.Equal(endWhen.Truncate(time.Second)) {
		t.Fatalf("create event = %+v, %v", event, err)
	}
	pastEnd := when.Add(-time.Minute)
	if _, err := repository.CreateEvent(Types.CreateEventRequest{
		Name: "Invalid end time", DateTime: when, EndDateTime: &pastEnd,
		Location: "Room", ClubName: "Club", Description: "Description", HostUserID: &host.ID,
	}); !errors.Is(err, ErrInvalidEvent) {
		t.Fatalf("create event with end time before start error = %v", err)
	}
	tooSoonEnd := when.Add(29 * time.Minute)
	if _, err := repository.CreateEvent(Types.CreateEventRequest{
		Name: "End time too soon", DateTime: when, EndDateTime: &tooSoonEnd,
		Location: "Room", ClubName: "Club", Description: "Description", HostUserID: &host.ID,
	}); !errors.Is(err, ErrInvalidEvent) {
		t.Fatalf("create event with end time less than 30 minutes after start error = %v", err)
	}
	if err := validateEventEndDateTime(when, when.Add(30*time.Minute)); err != nil {
		t.Fatalf("validate event with a 30-minute duration error = %v", err)
	}
	duplicate := Types.CreateEventRequest{
		Name: "Meeting", DateTime: when.Add(time.Hour), Location: "Room 2",
		ClubName: "Club B", Description: "Another meeting", HostUserID: &host.ID,
	}
	if _, err := repository.CreateEvent(duplicate); !errors.Is(err, ErrDuplicateEventName) {
		t.Fatalf("duplicate event name error = %v", err)
	}

	// attempt to get an event
	all, err := repository.GetEvents(Types.EventFilters{})
	if err != nil || len(all) != 1 {
		t.Fatalf("get events = %+v, %v", all, err)
	}

	// attempt to filter the list of events by various criteria
	start := when.Add(-time.Hour)
	end := when.Add(time.Hour)
	hostID := host.ID
	for _, filters := range []Types.EventFilters{
		{StartDate: &start},
		{EndDate: &end},
		{ClubName: "Club A"},
		{Location: "Room 1"},
		{HostUserID: &hostID},
		{ClubName: "Missing"},
	} {
		events, err := repository.GetEvents(filters)
		if err != nil {
			t.Fatalf("get filtered events (%+v): %v", filters, err)
		}
		if filters.ClubName == "Missing" && len(events) != 0 {
			t.Fatalf("expected no events, got %+v", events)
		}
	}

	// attempt to update an event with new information
	name := "Updated meeting"
	location := "Room 2"
	updated, err := repository.UpdateEvent(event.ID, Types.UpdateEventRequest{Name: &name, Location: &location})
	if err != nil || updated.Name != name || updated.Location != location {
		t.Fatalf("update event = %+v, %v", updated, err)
	}
	updatedEnd := when.Add(2 * time.Hour)
	updated, err = repository.UpdateEvent(event.ID, Types.UpdateEventRequest{
		EndDateTime: &Types.OptionalDateTime{Set: true, Value: &updatedEnd},
	})
	if err != nil || updated.EndDateTime == nil ||
		!updated.EndDateTime.Equal(updatedEnd.Truncate(time.Second)) {
		t.Fatalf("update event end time = %+v, %v", updated, err)
	}
	tooSoonEnd = when.Add(29 * time.Minute)
	if _, err := repository.UpdateEvent(event.ID, Types.UpdateEventRequest{
		EndDateTime: &Types.OptionalDateTime{Set: true, Value: &tooSoonEnd},
	}); !errors.Is(err, ErrInvalidEvent) {
		t.Fatalf("update event with end time less than 30 minutes after start error = %v", err)
	}
	updated, err = repository.UpdateEvent(event.ID, Types.UpdateEventRequest{
		EndDateTime: &Types.OptionalDateTime{Set: true},
	})
	if err != nil || updated.EndDateTime != nil {
		t.Fatalf("clear event end time = %+v, %v", updated, err)
	}
	otherEvent, err := repository.CreateEvent(Types.CreateEventRequest{
		Name: "Another event", DateTime: when.Add(2 * time.Hour), Location: "Room 3",
		ClubName: "Club C", Description: "Another event", HostUserID: &host.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repository.UpdateEvent(otherEvent.ID, Types.UpdateEventRequest{Name: &updated.Name}); !errors.Is(err, ErrDuplicateEventName) {
		t.Fatalf("duplicate event update error = %v", err)
	}

	// attempt to get a nonexistent event
	if _, err := repository.GetEvent(999); !errors.Is(err, ErrNotFound) {
		t.Fatalf("get missing event error = %v", err)
	}
	validRequest := Types.CreateEventRequest{
		Name: "Valid event", DateTime: when, Location: "Room", ClubName: "Club", Description: "Description", HostUserID: &host.ID,
	}
	maxLengthRequest := validRequest
	maxLengthRequest.Name = strings.Repeat("T", maxEventTitleLength)
	maxLengthRequest.ClubName = strings.Repeat("O", maxEventOrganizationLength)
	maxLengthRequest.Location = strings.Repeat("L", maxEventLocationLength)
	maxLengthRequest.Description = strings.Repeat("D", maxEventDescriptionLength)
	maxLengthRequest.ContactEmail = strings.Repeat("e", maxEventEmailLength)
	maxLengthRequest.WebsiteURL = strings.Repeat("w", maxEventWebsiteLength)
	if _, err := repository.CreateEvent(maxLengthRequest); err != nil {
		t.Fatalf("create event at text field limits: %v", err)
	}

	for _, test := range []struct {
		name   string
		update func(*Types.CreateEventRequest)
	}{
		{"title", func(request *Types.CreateEventRequest) {
			request.Name = strings.Repeat("T", maxEventTitleLength+1)
		}},
		{"organization", func(request *Types.CreateEventRequest) {
			request.ClubName = strings.Repeat("O", maxEventOrganizationLength+1)
		}},
		{"location", func(request *Types.CreateEventRequest) {
			request.Location = strings.Repeat("L", maxEventLocationLength+1)
		}},
		{"description", func(request *Types.CreateEventRequest) {
			request.Description = strings.Repeat("D", maxEventDescriptionLength+1)
		}},
		{"email", func(request *Types.CreateEventRequest) {
			request.ContactEmail = strings.Repeat("e", maxEventEmailLength+1)
		}},
		{"website", func(request *Types.CreateEventRequest) {
			request.WebsiteURL = strings.Repeat("w", maxEventWebsiteLength+1)
		}},
	} {
		t.Run("event text limit "+test.name, func(t *testing.T) {
			request := validRequest
			test.update(&request)
			if _, err := repository.CreateEvent(request); !errors.Is(err, ErrInvalidEvent) {
				t.Fatalf("create event with overlong %s error = %v", test.name, err)
			}
		})
	}

	for name, request := range map[string]Types.CreateEventRequest{
		"missing name":        {DateTime: validRequest.DateTime, Location: validRequest.Location, ClubName: validRequest.ClubName, HostUserID: validRequest.HostUserID},
		"missing date":        {Name: validRequest.Name, Location: validRequest.Location, ClubName: validRequest.ClubName, HostUserID: validRequest.HostUserID},
		"missing location":    {Name: validRequest.Name, DateTime: validRequest.DateTime, ClubName: validRequest.ClubName, HostUserID: validRequest.HostUserID},
		"missing club name":   {Name: validRequest.Name, DateTime: validRequest.DateTime, Location: validRequest.Location, HostUserID: validRequest.HostUserID},
		"misisng description": {Name: validRequest.Name, DateTime: validRequest.DateTime, Location: validRequest.Location, ClubName: validRequest.ClubName, HostUserID: validRequest.HostUserID},
		"missing host":        {Name: validRequest.Name, DateTime: validRequest.DateTime, Location: validRequest.Location, ClubName: validRequest.ClubName},
		"non-positive host":   {Name: validRequest.Name, DateTime: validRequest.DateTime, Location: validRequest.Location, ClubName: validRequest.ClubName, HostUserID: intPointer(0)},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := repository.CreateEvent(request); !errors.Is(err, ErrInvalidEvent) {
				t.Fatalf("create event error = %v", err)
			}
		})
	}
	// attempt to create an event with a nonexistent host
	nonexistentHost := validRequest
	nonexistentHost.HostUserID = intPointer(999)
	if _, err := repository.CreateEvent(nonexistentHost); !errors.Is(err, ErrNotFound) {
		t.Fatalf("invalid host error = %v", err)
	}
	// attempt to update an event with a nonexistent ID
	if _, err := repository.UpdateEvent(999, Types.UpdateEventRequest{}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("update missing event error = %v", err)
	}
	// attempt to update an event with a nonexistent host
	if _, err := repository.UpdateEvent(event.ID, Types.UpdateEventRequest{HostUserID: intPointer(999)}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("invalid replacement host error = %v", err)
	}
	if _, err := repository.UpdateEvent(event.ID, Types.UpdateEventRequest{Name: stringPointer("")}); !errors.Is(err, ErrInvalidEvent) {
		t.Fatalf("empty replacement name error = %v", err)
	}
	if _, err := repository.UpdateEvent(event.ID, Types.UpdateEventRequest{
		Name: stringPointer(strings.Repeat("T", maxEventTitleLength+1)),
	}); !errors.Is(err, ErrInvalidEvent) {
		t.Fatalf("overlong replacement name error = %v", err)
	}
	if _, err := repository.UpdateEvent(event.ID, Types.UpdateEventRequest{HostUserID: intPointer(0)}); !errors.Is(err, ErrInvalidEvent) {
		t.Fatalf("non-positive replacement host error = %v", err)
	}
	// attempt to create or update an event with a past or far future date
	if _, err := repository.CreateEvent(Types.CreateEventRequest{
		Name: "Past event", DateTime: time.Now().Add(-time.Minute), Location: "Room", ClubName: "Club",
		Description: "Description", HostUserID: &host.ID,
	}); !errors.Is(err, ErrInvalidEvent) {
		t.Fatalf("past event error = %v", err)
	}
	if _, err := repository.CreateEvent(Types.CreateEventRequest{
		Name: "Far future event", DateTime: time.Now().AddDate(0, 7, 0), Location: "Room", ClubName: "Club",
		Description: "Description", HostUserID: &host.ID,
	}); !errors.Is(err, ErrInvalidEvent) {
		t.Fatalf("far future event error = %v", err)
	}
	past := time.Now().Add(-time.Minute)
	if _, err := repository.UpdateEvent(event.ID, Types.UpdateEventRequest{DateTime: &past}); !errors.Is(err, ErrInvalidEvent) {
		t.Fatalf("past update error = %v", err)
	}
	farFuture := time.Now().AddDate(0, 7, 0)
	if _, err := repository.UpdateEvent(event.ID, Types.UpdateEventRequest{DateTime: &farFuture}); !errors.Is(err, ErrInvalidEvent) {
		t.Fatalf("far future update error = %v", err)
	}
	// attempt to delete an event
	if err := repository.DeleteEvent(event.ID); err != nil {
		t.Fatal(err)
	}
	// attempt to delete a nonexistent event
	if err := repository.DeleteEvent(event.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("delete missing event error = %v", err)
	}
}

func TestRepositoryRSVPs(t *testing.T) {
	repository := testRepository(t)
	host, _, _ := repository.CreateUser(Types.CreateUserRequest{Name: "Host", Email: "host@example.com"})
	attendee, _, _ := repository.CreateUser(Types.CreateUserRequest{Name: "Attendee", Email: "attendee@example.com"})
	event, err := repository.CreateEvent(Types.CreateEventRequest{
		Name: "Meeting", DateTime: time.Now().UTC().Add(24 * time.Hour), Location: "Room", ClubName: "Club", Description: "Description", HostUserID: &host.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := repository.AddRSVP(attendee.ID, event.ID); err != nil {
		t.Fatal(err)
	}
	if err := repository.AddRSVP(attendee.ID, event.ID); err != nil {
		t.Fatal(err)
	}
	rsvps, err := repository.GetUserRSVPs(attendee.ID)
	if err != nil || len(rsvps) != 1 || rsvps[0].AttendeeCount != 1 {
		t.Fatalf("get RSVPs = %+v, %v", rsvps, err)
	}
	if err := repository.RemoveRSVP(attendee.ID, event.ID); err != nil {
		t.Fatal(err)
	}
	if err := repository.RemoveRSVP(999, event.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing RSVP user error = %v", err)
	}
	if err := repository.AddRSVP(attendee.ID, 999); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing RSVP event error = %v", err)
	}
	if _, err := repository.GetUserRSVPs(999); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing RSVP list user error = %v", err)
	}
}

func TestRepositoryDeleteRelationships(t *testing.T) {
	repository := testRepository(t)
	host, _, err := repository.CreateUser(Types.CreateUserRequest{Name: "Host", Email: "host@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	attendee, _, err := repository.CreateUser(Types.CreateUserRequest{Name: "Attendee", Email: "attendee@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	event, err := repository.CreateEvent(Types.CreateEventRequest{
		Name: "Meeting", DateTime: time.Now().UTC().Add(24 * time.Hour), Location: "Room", ClubName: "Club", Description: "Description", HostUserID: &host.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	// attempt to add an RSVP
	if err := repository.AddRSVP(attendee.ID, event.ID); err != nil {
		t.Fatal(err)
	}
	// attempt to delete the host user, which should be restricted due to the event
	if err := repository.DeleteUser(host.ID); err == nil {
		t.Fatal("expected host deletion to be restricted")
	}
	// attempt to delete the attendee user, which should cascade and remove the RSVP
	if err := repository.DeleteUser(attendee.ID); err != nil {
		t.Fatalf("delete RSVP attendee: %v", err)
	}
	var rsvpCount int
	if err := repository.db.QueryRow(`SELECT COUNT(*) FROM rsvps WHERE user_id = ?`, attendee.ID).Scan(&rsvpCount); err != nil {
		t.Fatal(err)
	}
	// confirm the RSVP was removed
	if rsvpCount != 0 {
		t.Fatalf("expected attendee RSVP to cascade, found %d", rsvpCount)
	}

	// attempt to add another RSVP
	anotherAttendee, _, err := repository.CreateUser(Types.CreateUserRequest{Name: "Another", Email: "another@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.AddRSVP(anotherAttendee.ID, event.ID); err != nil {
		t.Fatal(err)
	}
	// attempt to delete the event, which should cascade and remove the RSVP
	if err := repository.DeleteEvent(event.ID); err != nil {
		t.Fatal(err)
	}
	if err := repository.db.QueryRow(`SELECT COUNT(*) FROM rsvps WHERE event_id = ?`, event.ID).Scan(&rsvpCount); err != nil {
		t.Fatal(err)
	}
	// confirm the RSVP was removed
	if rsvpCount != 0 {
		t.Fatalf("expected event RSVP to cascade, found %d", rsvpCount)
	}
}

// test that the database rejects 0 and negative IDs for users, events, and RSVPS
func TestRepositoryRejectsNonPositiveIDs(t *testing.T) {
	repository := testRepository(t)
	if _, err := repository.GetUser(0); !errors.Is(err, ErrInvalidID) {
		t.Fatalf("GetUser zero error = %v", err)
	}
	if _, err := repository.GetEvent(-1); !errors.Is(err, ErrInvalidID) {
		t.Fatalf("GetEvent negative error = %v", err)
	}
	if err := repository.DeleteUser(0); !errors.Is(err, ErrInvalidID) {
		t.Fatalf("DeleteUser zero error = %v", err)
	}
	if err := repository.DeleteEvent(-1); !errors.Is(err, ErrInvalidID) {
		t.Fatalf("DeleteEvent negative error = %v", err)
	}
	if err := repository.AddRSVP(0, 1); !errors.Is(err, ErrInvalidID) {
		t.Fatalf("AddRSVP zero user error = %v", err)
	}
	if err := repository.RemoveRSVP(1, -1); !errors.Is(err, ErrInvalidID) {
		t.Fatalf("RemoveRSVP negative event error = %v", err)
	}
	if _, err := repository.GetUserRSVPs(0); !errors.Is(err, ErrInvalidID) {
		t.Fatalf("GetUserRSVPs zero error = %v", err)
	}
}

// test that the database rejects requests when closed
func TestRepositoryDatabaseErrors(t *testing.T) {
	repository := testRepository(t)
	repository.db.Close()

	if _, _, err := repository.CreateUser(Types.CreateUserRequest{Name: "A", Email: "a@example.com"}); err == nil {
		t.Fatal("expected create user database error")
	}
	if _, err := repository.GetUser(1); err == nil {
		t.Fatal("expected get user database error")
	}
	if _, err := repository.GetAllUsers(); err == nil {
		t.Fatal("expected list users database error")
	}
	if _, err := repository.UpdateUser(1, Types.UpdateUserRequest{}); err == nil {
		t.Fatal("expected update user database error")
	}
	if err := repository.DeleteUser(1); err == nil {
		t.Fatal("expected delete user database error")
	}
	if _, err := repository.CreateEvent(Types.CreateEventRequest{HostUserID: intPointer(1)}); err == nil {
		t.Fatal("expected create event database error")
	}
	if _, err := repository.GetEvent(1); err == nil {
		t.Fatal("expected get event database error")
	}
	if _, err := repository.GetEvents(Types.EventFilters{}); err == nil {
		t.Fatal("expected list events database error")
	}
	if _, err := repository.UpdateEvent(1, Types.UpdateEventRequest{}); err == nil {
		t.Fatal("expected update event database error")
	}
	if err := repository.DeleteEvent(1); err == nil {
		t.Fatal("expected delete event database error")
	}
	if err := repository.AddRSVP(1, 1); err == nil {
		t.Fatal("expected add RSVP database error")
	}
	if err := repository.RemoveRSVP(1, 1); err == nil {
		t.Fatal("expected remove RSVP database error")
	}
	if _, err := repository.GetUserRSVPs(1); err == nil {
		t.Fatal("expected list RSVPs database error")
	}
}

// confirm that the database rejects malformed dates for events
func TestRepositoryMalformedDate(t *testing.T) {
	repository := testRepository(t)
	host, _, err := repository.CreateUser(Types.CreateUserRequest{Name: "Host", Email: "host@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = repository.db.Exec(`
		INSERT INTO events (
			event_name, date_time, location, club_name, description, host_user_id
		) VALUES (?, ?, ?, ?, ?, ?)
	`, "Bad date", "not-a-date", "Room", "Club", "Description", host.ID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repository.GetEvent(1); err == nil {
		t.Fatal("expected malformed event date error")
	}
	if _, err := repository.GetEvents(Types.EventFilters{}); err == nil {
		t.Fatal("expected malformed event list date error")
	}
}

// confirm that the database init and table creation functions error when given invalid paths or closed databases
func TestRepositoryDriverErrors(t *testing.T) {
	if _, err := Init(t.TempDir()); err == nil {
		t.Fatal("expected database initialization error for directory path")
	}
	db, err := Init(t.TempDir() + "\\valid.db")
	if err != nil {
		t.Fatal(err)
	}
	db.Close()
	if err := createTables(db); err == nil {
		t.Fatal("expected table creation error on closed database")
	}
}

func intPointer(value int) *int {
	return &value
}

func stringPointer(value string) *string {
	return &value
}
