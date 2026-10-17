# Campus Connect API

This document defines the API contract for the Campus Connect MVP.

For this MVP, every user must choose a unique username. Display names,
username changes, and authentication are deferred until a later version.
## Running the API

All request and response bodies use JSON. Event date/time values use ISO 8601
format in UTC, for example:

```text
2026-10-25T18:30:00Z
```

IDs are positive integers.

## MVP coverage

| Requirement | Endpoint |
| --- | --- |
| Discover campus events | `GET /events` |
| See a specific named event | `GET /events/:eventId` |
| View an event's date and time | Included in `GET /events` and `GET /events/:eventId` |
| View an event's location | Included in `GET /events` and `GET /events/:eventId` |
| Create or recover a user | `POST /users` and `GET /users?username=:username` |
| RSVP to an event | `POST /users/:userId/rsvps/:eventId` |
| See events a user RSVP'd to | `GET /users/:userId/rsvps` |
| Create an event | `POST /events` |
| Add details to an event | `POST /events` |

## Endpoint summary

| Method | Path | Purpose |
| --- | --- | --- |
| `POST` | `/users` | Create a user from a username |
| `GET` | `/users/:userId` | Get a user by ID |
| `GET` | `/users?username=:username` | Find a user by username |
| `GET` | `/events` | List events happening on campus |
| `GET` | `/events/:eventId` | Get one event by ID |
| `POST` | `/events` | Create an event |
| `POST` | `/users/:userId/rsvps/:eventId` | RSVP a user to an event |
| `DELETE` | `/users/:userId/rsvps/:eventId` | Remove a user's RSVP |
| `GET` | `/users/:userId/rsvps` | List a user's RSVP'd events |

The RSVP removal endpoint is included so the frontend can correct an RSVP;
removing an RSVP is not a separate MVP discovery feature.

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

`attendeeCount` is included for display purposes. Individual attendee
identities are not required by the MVP.

### User

```json
{
  "id": 1,
  "username": "alice"
}
```

`username` is required and must be unique. It is the user's stable identifier
for the MVP, but it is not a password and does not provide authentication.
The MVP does not have a separate display name.

### Create user request

```json
{
  "username": "alice"
}
```

Usernames should be normalized consistently by the backend, for example by
trimming whitespace and treating them case-insensitively. The response should
always return the canonical username.

## User endpoints

### Create a user

```http
POST /users
Content-Type: application/json
```

Creates a new user and returns the generated ID.

#### Response: `201 Created`

```json
{
  "id": 1,
  "username": "alice"
}
```

#### Errors

- `400 Bad Request` if the JSON is malformed or the username is empty.
- `409 Conflict` if the username is already in use.

### Get a user by ID

```http
GET /users/:userId
```

Returns the user represented by the ID. The frontend can use this endpoint
when it has persisted the ID from an earlier session.

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

### Find a user by username

```http
GET /users?username=alice
```

Returns the existing user with that exact username. This is the MVP's
username-based return flow: if a frontend only has the username the user
entered after a page reload, it can look up the user and recover the ID needed
by the RSVP endpoints.

#### Response: `200 OK`

```json
{
  "id": 1,
  "username": "alice"
}
```

#### Errors

- `400 Bad Request` if `username` is missing or empty.
- `404 Not Found` if no user has that username.

The frontend should first call `GET /users?username=...` when a user enters an
existing username. If the lookup returns `404`, it can offer to create the
user with `POST /users`. It should not blindly create a new user every time a
user returns, because that would produce duplicate identities.

### Create event request

```json
{
  "name": "Game Night",
  "dateTime": "2026-10-25T18:30:00Z",
  "location": "Room 101",
  "clubName": "Game Development Club"
}
```

All four fields are required. `clubName` is the event's host or organization
label; the MVP does not require a separate club resource.

## Endpoints

### List campus events

```http
GET /events
```

Returns every event currently available on campus.

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

Adds the event to the user's RSVP list and increments the event's attendee
count. The operation should be idempotent: RSVP'ing to the same event twice
must not create a duplicate RSVP or increment the count twice.

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

Removes the event from the user's RSVP list and decrements the event's
attendee count. Removing an RSVP that does not exist should be harmless.

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

The MVP has no authentication system or passwords. A username is only a
stable way to find the same user again; it does not prove that the person
entering it owns that identity.

When a user first enters a username:

1. Call `GET /users?username=:username`.
2. If it returns `200`, save the returned `id` as the active user ID.
3. If it returns `404`, call `POST /users` to create the user, then save the
   returned `id`.

The frontend should persist the returned ID locally when practical, such as in
`localStorage`, and can verify it later with `GET /users/:userId`. If local
storage is unavailable or cleared, the username lookup recovers the ID:

```text
GET /users?username=alice
-> { "id": 1, "username": "alice" }

POST /users/1/rsvps/4
```

A future authentication layer can replace this lookup and the explicit
`:userId` with the authenticated user from a session or token without changing
the event response shapes.

## Out of scope for this MVP

- Display names and username changes
- Passwords, authentication, and account recovery
- Searching or filtering events
- Editing or deleting events
- Club management
- Attendee identity lists
- Social connections, messaging, or friend features
- Pagination
