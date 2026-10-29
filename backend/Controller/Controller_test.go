package Controller

import (
	"bytes"
	"campus_connect/Model"
	"campus_connect/Types"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type controllerBackend struct {
	router *gin.Engine
	db     *sql.DB
}

// Initializes a temporary router / db pair for testing.
func newControllerBackend(t *testing.T) controllerBackend {
	t.Helper()
	db, err := Model.Init(t.TempDir() + "\\controller.db")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	router := newRouter(Model.NewRepository(db))
	return controllerBackend{router: router, db: db}
}

// helper function to make an HTTP request to the router and return the response
func request(t *testing.T, router http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var content *bytes.Reader
	if body == nil {
		content = bytes.NewReader(nil)
	} else {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		content = bytes.NewReader(data)
	}
	req := httptest.NewRequest(method, path, content)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	return recorder
}

func TestUserRoutes(t *testing.T) {
	controllerBackend := newControllerBackend(t)
	body := Types.CreateUserRequest{Name: "Alice", Email: "alice@example.com"}

	// attempt to create a valid user
	if response := request(t, controllerBackend.router, http.MethodPost, "/users", body); response.Code != http.StatusCreated {
		t.Fatalf("create user status = %d", response.Code)
	}
	// attempt to create a duplicate user
	if response := request(t, controllerBackend.router, http.MethodPost, "/users", body); response.Code != http.StatusOK {
		t.Fatalf("duplicate user status = %d", response.Code)
	}
	// attempt to create a user with invalid JSON
	if response := request(t, controllerBackend.router, http.MethodPost, "/users", "{"); response.Code != http.StatusBadRequest {
		t.Fatalf("invalid user JSON status = %d", response.Code)
	}
	// confirm the invalid JSON user was not created
	if response := request(t, controllerBackend.router, http.MethodPost, "/users", "{"); response.Code != http.StatusBadRequest {
		t.Fatalf("invalid user JSON status = %d", response.Code)
	}
	// attempt to get a string user ID
	if response := request(t, controllerBackend.router, http.MethodGet, "/users/nope", nil); response.Code != http.StatusBadRequest {
		t.Fatalf("invalid user ID status = %d", response.Code)
	}
	// attempt to get a 0 user ID
	if response := request(t, controllerBackend.router, http.MethodGet, "/users/0", nil); response.Code != http.StatusBadRequest {
		t.Fatalf("zero user ID status = %d", response.Code)
	}
	// attempt to get a negative user ID
	if response := request(t, controllerBackend.router, http.MethodGet, "/users/-1", nil); response.Code != http.StatusBadRequest {
		t.Fatalf("negative user ID status = %d", response.Code)
	}
	// attempt to get a nonexistent user ID
	if response := request(t, controllerBackend.router, http.MethodGet, "/users/999", nil); response.Code != http.StatusNotFound {
		t.Fatalf("missing user status = %d", response.Code)
	}
	// attempt to get a valid user ID
	if response := request(t, controllerBackend.router, http.MethodGet, "/users/1", nil); response.Code != http.StatusOK {
		t.Fatalf("get user status = %d", response.Code)
	}
}

func TestEventRoutes(t *testing.T) {
	controllerBackend := newControllerBackend(t)
	host, _, err := Model.NewRepository(controllerBackend.db).CreateUser(Types.CreateUserRequest{Name: "Host", Email: "host@example.com"})
	if err != nil {
		t.Fatal(err)
	}

	body := Types.CreateEventRequest{
		Name: "Meeting", DateTime: time.Now().UTC(), Location: "Room", ClubName: "Club", HostUserID: &host.ID,
	}

	// attempt to list an empty events list
	if response := request(t, controllerBackend.router, http.MethodGet, "/events", nil); response.Code != http.StatusOK {
		t.Fatalf("list events status = %d", response.Code)
	}
	// attempt to create a valid event
	if response := request(t, controllerBackend.router, http.MethodPost, "/events", body); response.Code != http.StatusCreated {
		t.Fatalf("create event status = %d", response.Code)
	}
	// attempt to get that event
	if response := request(t, controllerBackend.router, http.MethodGet, "/events/1", nil); response.Code != http.StatusOK {
		t.Fatalf("get event status = %d", response.Code)
	}
	// attempt to get an event with a string ID
	if response := request(t, controllerBackend.router, http.MethodGet, "/events/nope", nil); response.Code != http.StatusBadRequest {
		t.Fatalf("invalid event ID status = %d", response.Code)
	}
	// attempt to get an event with a 0 ID
	if response := request(t, controllerBackend.router, http.MethodGet, "/events/0", nil); response.Code != http.StatusBadRequest {
		t.Fatalf("zero event ID status = %d", response.Code)
	}
	// attempt to get an event with a negative ID
	if response := request(t, controllerBackend.router, http.MethodGet, "/events/-1", nil); response.Code != http.StatusBadRequest {
		t.Fatalf("negative event ID status = %d", response.Code)
	}
	// attempt to get an event with a nonexistent ID
	if response := request(t, controllerBackend.router, http.MethodGet, "/events/999", nil); response.Code != http.StatusNotFound {
		t.Fatalf("missing event status = %d", response.Code)
	}
	// attempt to create an event with missing required fields
	if response := request(t, controllerBackend.router, http.MethodPost, "/events", Types.CreateEventRequest{Name: "bad"}); response.Code != http.StatusBadRequest {
		t.Fatalf("invalid event status = %d", response.Code)
	}
	// attempt to create an event with a nonexistent host user ID
	nonexistentHostEvent := body
	nonexistentHostEvent.HostUserID = intPointer(999)
	if response := request(t, controllerBackend.router, http.MethodPost, "/events", nonexistentHostEvent); response.Code != http.StatusBadRequest {
		t.Fatalf("missing host event status = %d", response.Code)
	}
	// attempt to create an event with invalid JSON
	if response := request(t, controllerBackend.router, http.MethodPost, "/events", "{"); response.Code != http.StatusBadRequest {
		t.Fatalf("invalid event JSON status = %d", response.Code)
	}

	// attempt to update an event with a string ID
	if response := request(t, controllerBackend.router, http.MethodPatch, "/events/nope", Types.UpdateEventRequest{}); response.Code != http.StatusBadRequest {
		t.Fatalf("invalid update ID status = %d", response.Code)
	}
	// attempt to update an event with a nonexistent ID
	if response := request(t, controllerBackend.router, http.MethodPatch, "/events/999", Types.UpdateEventRequest{}); response.Code != http.StatusNotFound {
		t.Fatalf("missing update event status = %d", response.Code)
	}
	// attempt to update an event with a bad JSON
	if response := request(t, controllerBackend.router, http.MethodPatch, "/events/1", "{"); response.Code != http.StatusBadRequest {
		t.Fatalf("invalid update JSON status = %d", response.Code)
	}
	// attempt to update an event with valid data
	if response := request(t, controllerBackend.router, http.MethodPatch, "/events/1", Types.UpdateEventRequest{}); response.Code != http.StatusOK {
		t.Fatalf("update event status = %d", response.Code)
	}
}

// Test the behavior of each endpoint when the database is closed, which should result in 500 errors.
func TestControllerDatabaseErrors(t *testing.T) {
	controllerBackend := newControllerBackend(t)
	controllerBackend.db.Close()

	validUser := Types.CreateUserRequest{Name: "Alice", Email: "alice@example.com"}
	validEvent := Types.CreateEventRequest{Name: "Meeting", DateTime: time.Now().UTC(), Location: "Room", ClubName: "Club", HostUserID: intPointer(1)}

	for _, test := range []struct {
		name   string
		method string
		path   string
		body   any
		want   int
	}{
		{"create user", http.MethodPost, "/users", validUser, http.StatusInternalServerError},
		{"get user", http.MethodGet, "/users/1", nil, http.StatusInternalServerError},
		{"list events", http.MethodGet, "/events", nil, http.StatusInternalServerError},
		{"create event", http.MethodPost, "/events", validEvent, http.StatusInternalServerError},
		{"get event", http.MethodGet, "/events/1", nil, http.StatusInternalServerError},
		{"update event", http.MethodPatch, "/events/1", Types.UpdateEventRequest{}, http.StatusInternalServerError},
		{"add RSVP", http.MethodPost, "/users/1/rsvps/1", nil, http.StatusInternalServerError},
		{"remove RSVP", http.MethodDelete, "/users/1/rsvps/1", nil, http.StatusInternalServerError},
		{"list RSVPs", http.MethodGet, "/users/1/rsvps", nil, http.StatusInternalServerError},
	} {
		t.Run(test.name, func(t *testing.T) {
			response := request(t, controllerBackend.router, test.method, test.path, test.body)
			if response.Code != test.want {
				t.Fatalf("status = %d, want %d", response.Code, test.want)
			}
		})
	}
}

func TestRSVPRoutes(t *testing.T) {
	controllerBackend := newControllerBackend(t)
	repository := Model.NewRepository(controllerBackend.db)
	host, _, err := repository.CreateUser(Types.CreateUserRequest{Name: "Host", Email: "host@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	attendee, _, err := repository.CreateUser(Types.CreateUserRequest{Name: "Attendee", Email: "attendee@example.com"})
	if err != nil {
		t.Fatal(err)
	}
	event, err := repository.CreateEvent(Types.CreateEventRequest{
		Name: "Meeting", DateTime: time.Now().UTC(), Location: "Room", ClubName: "Club", HostUserID: &host.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	// attempt to add and remove RSVPs with various invalid and valid inputs
	if response := request(t, controllerBackend.router, http.MethodPost, "/users/nope/rsvps/1", nil); response.Code != http.StatusBadRequest {
		t.Fatalf("invalid RSVP user ID status = %d", response.Code)
	}
	if response := request(t, controllerBackend.router, http.MethodPost, "/users/1/rsvps/nope", nil); response.Code != http.StatusBadRequest {
		t.Fatalf("invalid RSVP event ID status = %d", response.Code)
	}
	if response := request(t, controllerBackend.router, http.MethodPost, "/users/999/rsvps/1", nil); response.Code != http.StatusNotFound {
		t.Fatalf("missing RSVP user status = %d", response.Code)
	}
	if response := request(t, controllerBackend.router, http.MethodPost, "/users/0/rsvps/1", nil); response.Code != http.StatusBadRequest {
		t.Fatalf("zero RSVP user status = %d", response.Code)
	}
	if response := request(t, controllerBackend.router, http.MethodPost, "/users/1/rsvps/-1", nil); response.Code != http.StatusBadRequest {
		t.Fatalf("negative RSVP event status = %d", response.Code)
	}
	if response := request(t, controllerBackend.router, http.MethodPost, "/users/1/rsvps/999", nil); response.Code != http.StatusNotFound {
		t.Fatalf("missing RSVP event status = %d", response.Code)
	}
	if response := request(t, controllerBackend.router, http.MethodPost, "/users/2/rsvps/1", nil); response.Code != http.StatusOK {
		t.Fatalf("add RSVP status = %d", response.Code)
	}
	if response := request(t, controllerBackend.router, http.MethodPost, "/users/2/rsvps/1", nil); response.Code != http.StatusOK {
		t.Fatalf("duplicate RSVP status = %d", response.Code)
	}
	if response := request(t, controllerBackend.router, http.MethodGet, "/users/nope/rsvps", nil); response.Code != http.StatusBadRequest {
		t.Fatalf("invalid RSVP list user ID status = %d", response.Code)
	}
	if response := request(t, controllerBackend.router, http.MethodGet, "/users/999/rsvps", nil); response.Code != http.StatusNotFound {
		t.Fatalf("missing RSVP list user status = %d", response.Code)
	}
	if response := request(t, controllerBackend.router, http.MethodGet, "/users/2/rsvps", nil); response.Code != http.StatusOK {
		t.Fatalf("get RSVPs status = %d", response.Code)
	}
	if response := request(t, controllerBackend.router, http.MethodDelete, "/users/2/rsvps/1", nil); response.Code != http.StatusOK {
		t.Fatalf("remove RSVP status = %d", response.Code)
	}
	if response := request(t, controllerBackend.router, http.MethodDelete, "/users/2/rsvps/1", nil); response.Code != http.StatusOK {
		t.Fatalf("remove missing RSVP status = %d", response.Code)
	}
	_ = attendee
	_ = event
}

func intPointer(value int) *int {
	return &value
}
