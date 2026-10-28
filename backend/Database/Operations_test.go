package Database

import (
	"campus_connect/Types"
	"database/sql"
	"errors"
	"testing"
	"time"
)

func testDatabase(t *testing.T) *sql.DB {
	t.Helper()
	db, err := Init(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestOperationsCreateAndFindUser(t *testing.T) {
	db := testDatabase(t)
	repository := NewRepository(db)

	request := Types.CreateUserRequest{Name: "Alice", Email: "alice@example.com"}
	user, created, err := repository.CreateUser(request)
	if err != nil {
		t.Fatal(err)
	}
	if !created || user.Name != request.Name || user.Email != request.Email {
		t.Fatalf("unexpected created user: %+v, created=%v", user, created)
	}

	existing, created, err := repository.CreateUser(request)
	if err != nil {
		t.Fatal(err)
	}
	if created || existing.ID != user.ID {
		t.Fatalf("expected existing user, got %+v, created=%v", existing, created)
	}
}

func TestOperationsEventUpdateAndRSVP(t *testing.T) {
	db := testDatabase(t)
	repository := NewRepository(db)

	host, _, err := repository.CreateUser(Types.CreateUserRequest{
		Name: "Host", Email: "host@example.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	attendee, _, err := repository.CreateUser(Types.CreateUserRequest{
		Name: "Attendee", Email: "attendee@example.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	event, err := repository.CreateEvent(Types.CreateEventRequest{
		Name:       "Meeting",
		DateTime:   time.Date(2026, 10, 27, 18, 0, 0, 0, time.UTC),
		Location:   "Room 1",
		ClubName:   "Club",
		HostUserID: &host.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	updated, err := repository.UpdateEvent(event.ID, Types.UpdateEventRequest{
		Location: stringPointer("Room 2"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Location != "Room 2" {
		t.Fatalf("event was not updated: %+v", updated)
	}

	if err := repository.AddRSVP(attendee.ID, event.ID); err != nil {
		t.Fatal(err)
	}
	if err := repository.AddRSVP(attendee.ID, event.ID); err != nil {
		t.Fatal(err)
	}
	rsvps, err := repository.GetUserRSVPs(attendee.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(rsvps) != 1 || rsvps[0].AttendeeCount != 1 {
		t.Fatalf("unexpected RSVPs: %+v", rsvps)
	}

	if err := repository.RemoveRSVP(attendee.ID, event.ID); err != nil {
		t.Fatal(err)
	}
}

func TestOperationsMissingRecords(t *testing.T) {
	db := testDatabase(t)
	repository := NewRepository(db)

	_, err := repository.GetUser(999)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func stringPointer(value string) *string {
	return &value
}
