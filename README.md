# Campus Connect API

This document defines the API contract for the Campus Connect MVP.

## Running the API

All request and response bodies use JSON. Event date/time values use ISO 8601
format in UTC, for example:

```text
2026-10-25T18:30:00Z
```

IDs are positive integers.

## Endpoint summary

| Method | Path | Purpose |
| --- | --- | --- |
| `POST` | `/users` | Create or return a user by username |
| `GET` | `/users/:userId` | Get a user by ID |
| `GET` | `/events` | List or filter events happening on campus |
| `GET` | `/events/:eventId` | Get one event by ID |
| `POST` | `/events` | Create an event |
| `POST` | `/users/:userId/rsvps/:eventId` | RSVP a user to an event |
| `DELETE` | `/users/:userId/rsvps/:eventId` | Remove a user's RSVP |
| `GET` | `/users/:userId/rsvps` | List a user's RSVP'd events |

## Data shapes

### Event

```json
{
  "id": 1,
  "name": "Game Night",
  "dateTime": "2026-10-25T18:30:00Z",
  "location": "Room 101",
  "clubName": "Game Development Club",
  "attendeeCount": 12
}
```

### User

```json
{
  "id": 1,
  "username": "alice"
}
```

`username` must be unique within the campus.

### Create user request

```json
{
  "username": "alice"
}
```

## User endpoints

### Create a user

```http
POST /users
Content-Type: application/json
```

Creates a user if the username is new. If the username already exists, returns
the existing user. The frontend should save the returned user object, especially
`id`, because RSVP endpoints use the numeric ID.

#### Response: `201 Created` or `200 OK`

```json
{
  "id": 1,
  "username": "alice"
}
```

#### Errors

- `400 Bad Request` if the JSON is malformed or the username is empty.
- `200 OK` if the username already exists.

### Get a user by ID

```http
GET /users/:userId
```

Returns the user represented by the ID.

#### Response: `200 OK`

```json
{
  "id": 1,
  "username": "alice"
}
```

#### Errors

- `400 Bad Request` if `userId` is not a valid integer.
- `404 Not Found` if no user has that ID.

### Create event request

```json
{
  "name": "Game Night",
  "dateTime": "2026-10-25T18:30:00Z",
  "location": "Room 101",
  "clubName": "Game Development Club"
}
```

## Endpoints

### List or filter campus events

```http
GET /events
```

Returns every event currently available on campus. Optional query parameters
can narrow the results by an inclusive date range, club name, or both:

```http
GET /events?startDate=2026-10-25T00:00:00Z&endDate=2026-10-31T23:59:59Z
GET /events?clubName=Game%20Development%20Club
GET /events?startDate=2026-10-25T00:00:00Z&endDate=2026-10-31T23:59:59Z&clubName=Game%20Development%20Club
```

`startDate` and `endDate` must be ISO 8601 timestamps. Either boundary may be
omitted. `clubName` is an exact, case-sensitive match in the MVP. 
Results are sorted by `dateTime`, then by `id`.

#### Response: `200 OK`

```json
{
  "events": [
    {
      "id": 1,
      "name": "Game Night",
      "dateTime": "2026-10-25T18:30:00Z",
      "location": "Room 101",
      "clubName": "Game Development Club",
      "attendeeCount": 12
    }
  ]
}
```

An empty campus returns:

```json
{
  "events": []
}
```

An empty filtered result is the same:

```json
{
  "events": []
}
```

#### Errors

- `400 Bad Request` if a date is not a valid ISO 8601 timestamp.
- `400 Bad Request` if `startDate` is after `endDate`.

### Get one event

```http
GET /events/:eventId
```

Returns the name, date/time, location, host/club label, and attendee count for
one event.

#### Response: `200 OK`

```json
{
  "id": 1,
  "name": "Game Night",
  "dateTime": "2026-10-25T18:30:00Z",
  "location": "Room 101",
  "clubName": "Game Development Club",
  "attendeeCount": 12
}
```

#### Errors

- `400 Bad Request` if `eventId` is not a valid integer.
- `404 Not Found` if no event has that ID.

### Create an event

```http
POST /events
Content-Type: application/json
```

Creates an event with all details needed for the Host MVP feature.

#### Request

```json
{
  "name": "Game Night",
  "dateTime": "2026-10-25T18:30:00Z",
  "location": "Room 101",
  "clubName": "Game Development Club"
}
```

#### Response: `201 Created`

```json
{
  "id": 1,
  "name": "Game Night",
  "dateTime": "2026-10-25T18:30:00Z",
  "location": "Room 101",
  "clubName": "Game Development Club",
  "attendeeCount": 0
}
```

#### Errors

- `400 Bad Request` if JSON is malformed, a required field is missing, or
  `dateTime` is not a valid ISO 8601 timestamp.

### RSVP to an event

```http
POST /users/:userId/rsvps/:eventId
```

Adds the event to the user's RSVP list, adds the user to the event's attendence list, and increments the event's attendee
count.

#### Response: `200 OK`

```json
{
  "userId": 1,
  "eventId": 1,
  "rsvpd": true
}
```

#### Errors

- `400 Bad Request` if either path ID is not a valid integer.
- `404 Not Found` if the user or event does not exist.

### Remove an RSVP

```http
DELETE /users/:userId/rsvps/:eventId
```

Removes the event from the user's RSVP list, removes the user from the event's attendence list, and decrements the event's
attendee count.

#### Response: `204 No Content`

The response has no body.

#### Errors

- `400 Bad Request` if either path ID is not a valid integer.
- `404 Not Found` if the user or event does not exist.

### List a user's RSVP'd events

```http
GET /users/:userId/rsvps
```

Returns complete event objects so the frontend can render the user's RSVP list
without making one request per event.

#### Response: `200 OK`

```json
{
  "events": [
    {
      "id": 1,
      "name": "Game Night",
      "dateTime": "2026-10-25T18:30:00Z",
      "location": "Room 101",
      "clubName": "Game Development Club",
      "attendeeCount": 12
    }
  ]
}
```

If the user has no RSVPs, return an empty `events` array.

#### Errors

- `400 Bad Request` if `userId` is not a valid integer.
- `404 Not Found` if the user does not exist.

## Common error shape

All error responses should use the same shape:

```json
{
  "error": "event not found"
}
```

The message may be more specific, but the top-level `error` field should remain
consistent so the frontend can display or log failures predictably.

## User identity and return flow

When a user enters a username, call `POST /users`. It returns the existing user
if the username is already registered, or creates and returns a new user.

The frontend should persist the returned ID locally when practical, such as in
`localStorage`, and can verify it later with `GET /users/:userId`. If local
storage is unavailable or cleared, repeat the request:

```text
POST /users
-> { "id": 1, "username": "alice" }

POST /users/1/rsvps/4
```

## Typical frontend workflow

The following sequence is sufficient to build the MVP website or app:

1. **Identify the user**
   - Call `POST /users` with the username and store the returned user object.
   - Persist the returned `id` locally when possible. If local storage is
     cleared, repeat the request.
2. **Discover events**
   - Call `GET /events` to render the campus event list.
3. **Inspect an event**
   - Call `GET /events/:eventId` for the selected event's complete details.
4. **Host an event**
   - Submit the event form to `POST /events`.
   - Use the returned event object to show the newly created event.
5. **RSVP**
   - Call `POST /users/:userId/rsvps/:eventId` using the IDs returned by the
     user and event endpoints.
6. **Show the user's RSVP list**
   - Call `GET /users/:userId/rsvps`.
7. **Cancel an RSVP**
   - Call `DELETE /users/:userId/rsvps/:eventId`.