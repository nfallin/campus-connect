# Campus Connect API Contract

This document defines the HTTP and JSON contract for the Campus Connect MVP.
The API is intended to support a frontend for discovering campus events,
creating events, and managing RSVPs.

All request and response bodies use JSON unless a response is
documented as having no body. Event `dateTime` values use RFC 3339 / ISO 8601
format, for example:

```text
2026-10-25T18:30:00Z
```

IDs are positive integers.

## How the frontend talks to the backend

- All API calls go through `src/lib/api.ts`, which prefixes requests with
  `VITE_API_BASE_URL` (defaults to `/api`, see `.env.example`).
- In dev, `vite.config.ts` proxies anything under `/api` to
  `http://localhost:8080` and strips the `/api` prefix, since the backend's
  routes (`/events`, `/users`, ...) don't use that prefix themselves. So a
  frontend call to `/api/events` reaches the backend's `/events` route.
- In production there is no Vite proxy, so either put a reverse proxy in
  front of both services doing the same `/api` rewrite, or set
  `VITE_API_BASE_URL` to the backend's full deployed URL and enable CORS on
  the backend for the frontend's origin.

## Data types

### User

```json
{
  "id": 1,
  "name": "Alice",
  "email": "alice@example.com"
}
```

### Event

```json
{
  "id": 1,
  "name": "Game Night",
  "dateTime": "2026-10-25T18:30:00Z",
  "location": "Room 101",
  "clubName": "Game Development Club",
  "host": {
    "id": 1,
    "name": "Alice",
    "email": "alice@example.com"
  },
  "contactEmail": "alice@example.com",
  "websiteURL": "https://example.com/game-night",
  "attendeeCount": 12
}
```

### Error

```json
{
  "message": "event not found"
}
```

## Endpoints

### Create or return a user

```http
POST /users
Content-Type: application/json
```

#### Request body

```json
{
  "name": "Alice",
  "email": "alice@example.com"
}
```

#### Response: `201 Created`

Returned when the name/email pair is new.

```json
{
  "id": 1,
  "name": "Alice",
  "email": "alice@example.com"
}
```

#### Response: `200 OK`

Returned when the name/email pair already exists. The existing user is returned.

```json
{
  "id": 1,
  "name": "Alice",
  "email": "alice@example.com"
}
```

#### Response: `400 Bad Request`

Returned when the JSON is malformed or the required `name` field is missing
or empty.

### Get a user

```http
GET /users/:userId
```

#### Response: `200 OK`

```json
{
  "id": 1,
  "name": "Alice",
  "email": "alice@example.com"
}
```

#### Response: `400 Bad Request`

Returned when `id` is not an integer.

```json
{
  "message": "user id must be an integer"
}
```

#### Response: `404 Not Found`

Returned when no user has that ID.

```json
{
  "message": "user not found"
}
```

## Events

### List events

```http
GET /events
```

Returns all events. Results are sorted by `dateTime`, then by `id`.

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

If there are no events:

```json
{
  "events": []
}
```

The MVP does not define event filtering. Filtering by date, club, or other
criteria may be added after the MVP.

### Get an event

```http
GET /events/:id
```

#### Response: `200 OK`

Returns one complete event object.

#### Response: `400 Bad Request`

Returned when `id` is not an integer.

```json
{
  "message": "event id must be an integer"
}
```

#### Response: `404 Not Found`

Returned when no event has that ID.

```json
{
  "message": "event not found"
}
```

### Create an event

```http
POST /events
Content-Type: application/json
```

#### Request body

```json
{
  "name": "Game Night",
  "dateTime": "2026-10-25T18:30:00Z",
  "location": "Room 101",
  "clubName": "Game Development Club",
  "hostUserId": 1,
  "contactEmail": "alice@example.com",
  "websiteURL": "https://example.com/game-night"
}
```

`hostUserId` is required and must identify an existing user.
`contactEmail` and `websiteURL` are optional and may be empty.

#### Response: `201 Created`

Returns the newly created event:

```json
{
  "id": 1,
  "name": "Game Night",
  "dateTime": "2026-10-25T18:30:00Z",
  "location": "Room 101",
  "clubName": "Game Development Club",
  "host": {
    "id": 1,
    "name": "Alice",
    "email": "alice@example.com"
  },
  "contactEmail": "alice@example.com",
  "websiteURL": "https://example.com/game-night",
  "attendeeCount": 0
}
```

#### Response: `400 Bad Request`

Returned when the JSON is malformed, a required field is missing, or
`dateTime` is not a valid timestamp.

### Edit an event

```http
PATCH /events/:id
Content-Type: application/json
```

All fields are optional. Omitted fields keep their existing values, while provided fields replace them.

```json
{
  "name": "Updated Game Night",
  "location": "Room 202",
  "contactEmail": "events@example.com",
  "websiteURL": "https://example.com/updated-game-night"
}
```

The response is the updated event object with status `200 OK`. The endpoint
returns `404 Not Found` when the event does not exist and `400 Bad Request`
when the event ID, JSON, or host user ID is invalid.

## RSVPs

### Add an RSVP

```http
POST /users/:userId/rsvps/:eventId
```

Adds the event to the user's RSVP list and the user to the event's attendee
list. The operation is idempotent: adding the same RSVP more than once does
not create duplicates.

#### Response: `200 OK`

The response has no body.

#### Response: `400 Bad Request`

Returned when either path ID is not an integer.

#### Response: `404 Not Found`

Returned when the user or event does not exist.

### Remove an RSVP

```http
DELETE /users/:userId/rsvps/:eventId
```

Removes the event from the user's RSVP list and the user from the event's
attendee list. Removing an RSVP that does not exist is harmless.

#### Response: `200 OK`

The response has no body.

#### Response: `400 Bad Request`

Returned when either path ID is not an integer.

#### Response: `404 Not Found`

Returned when the user or event does not exist.

### List a user's RSVPs

```http
GET /users/:userId/rsvps
```

Returns complete event objects for every event the user has RSVP'd to.

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

If the user has no RSVPs:

```json
{
  "events": []
}
```

#### Response: `400 Bad Request`

Returned when `userId` is not an integer.

#### Response: `404 Not Found`

Returned when the user does not exist.

## Expected Usage

1. Call `POST /users` with the user's name and email.
2. Save the returned user `id`.
3. Call `GET /events` to display available events.
4. Call `GET /events/:id` to obtain one event.
5. Call `POST /events` to create a new event.
6. Call `POST /users/:userId/rsvps/:eventId` to RSVP.
7. Call `GET /users/:userId/rsvps` to display the user's RSVPs.
8. Call `DELETE /users/:userId/rsvps/:eventId` to cancel an RSVP.