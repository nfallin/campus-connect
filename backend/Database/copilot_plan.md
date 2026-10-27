That structure is a good next step. The current backend is already organized around the same three concepts, so the main change would be replacing the `Campus` maps and in-memory relationship management with database queries.

## Recommended schema

I would make the tables slightly more explicit:

```sql
CREATE TABLE users (
    user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_name TEXT NOT NULL,
    user_email TEXT NOT NULL,
    UNIQUE (user_name, user_email)
);

CREATE TABLE events (
    event_id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_name TEXT NOT NULL,
    date_time DATETIME NOT NULL,
    location TEXT NOT NULL,
    club_name TEXT NOT NULL,
    contact_email TEXT,
    website TEXT,
    host_user_id INTEGER NOT NULL,
    FOREIGN KEY (host_user_id)
        REFERENCES users(user_id)
        ON DELETE RESTRICT
);

CREATE TABLE rsvps (
    user_id INTEGER NOT NULL,
    event_id INTEGER NOT NULL,
    PRIMARY KEY (user_id, event_id),
    FOREIGN KEY (user_id)
        REFERENCES users(user_id)
        ON DELETE CASCADE,
    FOREIGN KEY (event_id)
        REFERENCES events(event_id)
        ON DELETE CASCADE
);
```

A few important details:

- `UNIQUE (user_name, user_email)` allows two users with the same name as long as their emails differ.
- The same person cannot accidentally be created twice with the exact same name/email pair.
- The composite primary key on `rsvps` prevents duplicate RSVPs.
- Foreign keys enforce that events have valid hosts and RSVPs refer to existing users/events.
- `ON DELETE CASCADE` makes RSVP cleanup automatic if an event or user is deleted.
- I would use plural table names such as `users`, `events`, and `rsvps`, although singular names are also valid if the project prefers them.

The schema does not necessarily need an email uniqueness constraint by itself. Your stated identity rule is name plus email, so uniqueness should apply to the pair, not email alone.

## API change for users

The current API accepts only:

```json
{
  "name": "Alice"
}
```

It should become:

```json
{
  "name": "Alice",
  "email": "alice@example.com"
}
```

`POST /users` would then behave like this:

1. Look up a user by both `name` and `email`.
2. If found, return the existing user with `200 OK`.
3. If not found, insert a new user and return it with `201 Created`.

The response should include the email:

```json
{
  "id": 1,
  "name": "Alice",
  "email": "alice@example.com"
}
```

The API documentation in `API.md` would need to be updated in all of these places:

- User data type
- `POST /users` request body
- `POST /users` response examples
- `GET /users/:userId` response
- Expected usage instructions
- Any event host examples containing a user object

The lookup should be performed by the database using an exact query such as:

```sql
SELECT user_id, user_name, user_email
FROM users
WHERE user_name = ? AND user_email = ?;
```

The frontend should store the returned numeric user ID, just as it does today. The user enters name and email when identifying themselves; subsequent event and RSVP requests can continue using `userId`.

## What changes in the backend

### 1. Add a database package

The backend would need a database connection initialized during startup. A likely structure would be:

```text
backend/
    Database/
        Database.go
        migrations.go
    Model/
    Controller/
```

SQLite would probably be the simplest choice for this application:

- No separate database server is required.
- It stores data in a local file.
- It supports foreign keys, transactions, indexes, and constraints.
- It is easy to use in local development and tests.

The application would open something like `campus-connect.db`, enable foreign keys, create the tables if necessary, and pass the database handle into the API.

The exact SQLite driver would need to be chosen. A pure-Go driver is attractive because it avoids requiring a C compiler during builds, but that is a dependency decision worth making deliberately.

### 2. Replace `Model.Campus`

Currently, `Campus` owns:

- `users map[int]*User`
- `events map[int]*Event`
- generated IDs
- RSVP relationships
- event sorting/filtering
- user lookup

With a database, those responsibilities move to SQL-backed repository methods. For example:

```go
CreateUser(name, email string) (*User, error)
GetUser(id int) (*User, error)
GetUserByNameAndEmail(name, email string) (*User, error)

ListEvents() ([]*Event, error)
GetEvent(id int) (*Event, error)
CreateEvent(...) (*Event, error)
UpdateEvent(...) (*Event, error)

AddRSVP(userID, eventID int) error
RemoveRSVP(userID, eventID int) error
ListUserRSVPs(userID int) ([]*Event, error)
```

The current `User` and `Event` types could remain as response/domain types, but they should no longer be responsible for mutating maps. Methods like:

```go
user.RsvpAdd(event)
user.RsvpRemove(event)
```

would be replaced by database operations. RSVP behavior belongs in a repository or service because it changes a database row rather than two in-memory maps.

### 3. Remove or reduce the API mutex

Currently `API` uses a `sync.RWMutex` because multiple requests access the same in-memory maps.

A database connection is designed for concurrent use, so the application-level read/write lock would generally no longer be needed. Removing it would allow independent requests to proceed concurrently.

However, transactions would still be needed for operations that must be atomic. For example, an RSVP operation should ensure that:

- The user exists.
- The event exists.
- The RSVP row is inserted.

The foreign keys and a transaction can enforce this reliably.

The comment in `Web.go` saying that the locks will stop being necessary with a database would then become obsolete.

## Endpoint impact

### `POST /users`

Change request:

```json
{
  "name": "Alice",
  "email": "alice@example.com"
}
```

Potential validation:

- Both fields required.
- Trim whitespace.
- Validate email format.
- Decide whether email matching is case-insensitive.

I would recommend treating emails case-insensitively for lookup. For example, `Alice@example.com` and `alice@example.com` should usually identify the same email address. That can be handled by normalizing the email before storage, or by storing a separate normalized value.

One design option is:

```text
user_email          -- original/display value
normalized_email    -- lowercase lookup value
```

Then the unique constraint can apply to `(user_name, normalized_email)`. This is more robust, but may be unnecessary for the first implementation.

### Event endpoints

The JSON shape can remain mostly unchanged. `hostUserId` would still identify the host:

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

The backend would insert the event using `host_user_id`, then query the host user when constructing the response.

One naming mismatch to consider is that the database proposal calls the field `website`, while the existing API calls it `websiteURL`. I would keep `websiteURL` in the HTTP API for compatibility and use `website` or `website_url` only internally in SQL.

### RSVP endpoints

The URLs can remain unchanged:

```text
POST   /users/:userId/rsvps/:eventId
DELETE /users/:userId/rsvps/:eventId
GET    /users/:userId/rsvps
```

The implementation changes from map mutation to SQL:

```sql
INSERT INTO rsvps (user_id, event_id)
VALUES (?, ?)
ON CONFLICT (user_id, event_id) DO NOTHING;
```

For removal:

```sql
DELETE FROM rsvps
WHERE user_id = ? AND event_id = ?;
```

The current idempotent behavior can therefore be preserved exactly.

For listing RSVPs, the current code iterates through a user’s RSVP map. The database version would use a join:

```sql
SELECT e.event_id, e.event_name, e.date_time, ...
FROM events e
JOIN rsvps r ON r.event_id = e.event_id
WHERE r.user_id = ?
ORDER BY e.date_time, e.event_id;
```

## Event host representation

The event response currently embeds a host object:

```json
"host": {
  "id": 1,
  "name": "Alice"
}
```

After adding email, it should probably become:

```json
"host": {
  "id": 1,
  "name": "Alice",
  "email": "alice@example.com"
}
```

That requires the event query to join the users table:

```sql
SELECT
    e.event_id,
    e.event_name,
    e.date_time,
    e.location,
    e.club_name,
    e.contact_email,
    e.website,
    u.user_id,
    u.user_name,
    u.user_email
FROM events e
JOIN users u ON u.user_id = e.host_user_id;
```

The API’s `UserResponse` type in `Types.go` would gain an `Email` field.

## Frontend impact

The frontend API wrapper in `api.ts` probably does not need structural changes. It already sends arbitrary JSON.

The frontend will need changes wherever it creates or identifies a user:

- Add an email input.
- Send `{ name, email }` to `POST /users`.
- Retain the returned `id`.
- Display or use the returned email where appropriate.

The larger `App.tsx` file should be checked for any assumptions that a user consists only of `id` and `name`.

The frontend does not need to know that the backend now uses a database. As long as the API response shape remains compatible, event and RSVP functionality should continue working.

## Important design decisions to settle before implementation

### Database file location

For local development:

```text
backend/campus-connect.db
```

For deployment, it would be better to configure the path through an environment variable:

```text
CAMPUS_CONNECT_DB_PATH
```

That avoids hardcoding a writable path and makes testing easier.

### Persistence across tests

Database-backed tests should use a separate temporary or in-memory database per test. Otherwise, tests could affect one another.

A useful approach is:

- Production: file-backed SQLite.
- Unit/integration tests: in-memory SQLite database.
- Each test creates its own schema.

### Schema initialization

For an MVP, creating tables on startup may be enough. A more maintainable approach is a small migration system, even if it initially has only one migration:

```text
001_create_users_events_rsvps
```

This matters if the database schema changes later, especially since adding the user email column would be a migration if existing data already exists.

### Existing in-memory data

There is no persistent production data to migrate currently, so the first database version can start empty. If sample data is eventually added, it should be inserted through a seed mechanism rather than hardcoded into the controllers.

### Deletion endpoints

The current API does not expose event or user deletion. The database should still define deletion behavior now because foreign keys depend on it:

- Deleting an event should delete its RSVPs.
- Deleting a user should probably be restricted if they host events.
- Deleting a user could cascade their RSVPs, but not events they host.

## Suggested implementation sequence

1. Decide on SQLite and the driver.
2. Add database initialization and schema creation.
3. Add database-backed user operations, including name/email lookup.
4. Update user request/response types and handlers.
5. Add database-backed event operations.
6. Add database-backed RSVP operations.
7. Remove the in-memory `Campus` ownership and API mutex.
8. Update frontend user creation/identification.
9. Update `API.md`.
10. Add database/repository tests covering:
   - Same name, different emails
   - Same name and email returning the existing user
   - Missing/invalid email
   - Event host foreign-key validation
   - Duplicate RSVP behavior
   - RSVP cleanup
   - Persistence after restarting the server

The biggest conceptual shift is that `User` and `Event` stop being mutable containers that know about each other. They become records loaded from the database, while repositories/services handle relationships and transactions. That separation will make the application more reliable and will also make future features—searching, filtering, authentication, and pagination—much easier to add.








In Go, SQLite is typically used through Go’s standard [`database/sql`](https://pkg.go.dev/database/sql) package plus a SQLite driver.

## 1. Choose a SQLite driver

For this project, I would recommend a pure-Go driver:

```bash
cd backend
go get modernc.org/sqlite
```

This avoids CGO and is generally easier to build on Windows and in CI.

The other common option is:

```bash
go get github.com/mattn/go-sqlite3
```

`go-sqlite3` is widely used, but it requires CGO and a C compiler, which adds setup complexity.

## 2. Open the database

With `modernc.org/sqlite`:

```go
package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
```

Then from application startup:

```go
db, err := database.Open("campus-connect.db")
if err != nil {
	log.Fatal(err)
}
defer db.Close()
```

`sql.Open` creates a database handle, but it may not establish a connection immediately. `db.Ping()` verifies that the database is usable.

The blank import:

```go
_ "modernc.org/sqlite"
```

registers the SQLite driver with `database/sql`.

## 3. Create the tables

You can initialize the schema with one `Exec` call containing multiple statements:

```go
func InitializeSchema(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			user_id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_name TEXT NOT NULL,
			user_email TEXT NOT NULL,
			UNIQUE (user_name, user_email)
		);

		CREATE TABLE IF NOT EXISTS events (
			event_id INTEGER PRIMARY KEY AUTOINCREMENT,
			event_name TEXT NOT NULL,
			date_time TEXT NOT NULL,
			location TEXT NOT NULL,
			club_name TEXT NOT NULL,
			contact_email TEXT,
			website TEXT,
			host_user_id INTEGER NOT NULL,
			FOREIGN KEY (host_user_id)
				REFERENCES users(user_id)
				ON DELETE RESTRICT
		);

		CREATE TABLE IF NOT EXISTS rsvps (
			user_id INTEGER NOT NULL,
			event_id INTEGER NOT NULL,
			PRIMARY KEY (user_id, event_id),
			FOREIGN KEY (user_id)
				REFERENCES users(user_id)
				ON DELETE CASCADE,
			FOREIGN KEY (event_id)
				REFERENCES events(event_id)
				ON DELETE CASCADE
		);
	`)

	return err
}
```

A typical startup sequence would be:

```go
db, err := database.Open("campus-connect.db")
if err != nil {
	log.Fatal(err)
}
defer db.Close()

if err := database.InitializeSchema(db); err != nil {
	log.Fatal(err)
}
```

For a first version, `CREATE TABLE IF NOT EXISTS` is enough. Later, you could add a proper migration system.

## 4. Insert a user

Use placeholders rather than concatenating values into SQL:

```go
func CreateUser(db *sql.DB, name, email string) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO users (user_name, user_email)
		VALUES (?, ?)
	`, name, email)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}
```

The `?` placeholders protect against SQL injection and properly handle escaping.

## 5. Find a user by name and email

```go
type User struct {
	ID    int
	Name  string
	Email string
}

func GetUserByNameAndEmail(
	db *sql.DB,
	name string,
	email string,
) (*User, error) {
	var user User

	err := db.QueryRow(`
		SELECT user_id, user_name, user_email
		FROM users
		WHERE user_name = ? AND user_email = ?
	`, name, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
```

Then `POST /users` can implement the existing “create or return” behavior:

```go
existingUser, err := GetUserByNameAndEmail(db, request.Name, request.Email)
if err != nil {
	// return 500
}

if existingUser != nil {
	// return 200 with existing user
}

userID, err := CreateUser(db, request.Name, request.Email)
// return 201 with new user
```

You should still keep the database uniqueness constraint because two concurrent requests could both check for an existing user before either inserts. The database constraint is the final protection.

## 6. Read events

For multiple rows, use `Query`:

```go
rows, err := db.Query(`
	SELECT
		e.event_id,
		e.event_name,
		e.date_time,
		e.location,
		e.club_name,
		e.contact_email,
		e.website,
		e.host_user_id,
		u.user_name,
		u.user_email
	FROM events e
	JOIN users u ON u.user_id = e.host_user_id
	ORDER BY e.date_time, e.event_id
`)
if err != nil {
	return err
}
defer rows.Close()

for rows.Next() {
	var event Event

	err := rows.Scan(
		&event.ID,
		&event.Name,
		&event.DateTime,
		&event.Location,
		&event.ClubName,
		&event.ContactEmail,
		&event.Website,
		&event.HostUserID,
		&event.HostName,
		&event.HostEmail,
	)
	if err != nil {
		return err
	}

	// Add event to response slice.
}

if err := rows.Err(); err != nil {
	return err
}
```

For exactly one row, use `QueryRow`, as with the user lookup.

## 7. Insert an event

```go
func CreateEvent(db *sql.DB, event Event) (int64, error) {
	result, err := db.Exec(`
		INSERT INTO events (
			event_name,
			date_time,
			location,
			club_name,
			contact_email,
			website,
			host_user_id
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
		event.Name,
		event.DateTime,
		event.Location,
		event.ClubName,
		event.ContactEmail,
		event.Website,
		event.HostUserID,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}
```

For dates, you can store RFC 3339 strings:

```go
event.DateTime.Format(time.RFC3339)
```

and parse them after reading:

```go
dateTime, err := time.Parse(time.RFC3339, storedDateTime)
```

SQLite does not have a dedicated timestamp type, so storing normalized RFC 3339 text is a practical choice.

## 8. Add and remove RSVPs

The composite primary key makes duplicate RSVPs impossible:

```go
func AddRSVP(db *sql.DB, userID, eventID int) error {
	_, err := db.Exec(`
		INSERT INTO rsvps (user_id, event_id)
		VALUES (?, ?)
		ON CONFLICT (user_id, event_id) DO NOTHING
	`, userID, eventID)

	return err
}
```

Remove one:

```go
func RemoveRSVP(db *sql.DB, userID, eventID int) error {
	_, err := db.Exec(`
		DELETE FROM rsvps
		WHERE user_id = ? AND event_id = ?
	`, userID, eventID)

	return err
}
```

List a user’s RSVP events:

```go
rows, err := db.Query(`
	SELECT
		e.event_id,
		e.event_name,
		e.date_time,
		e.location,
		e.club_name,
		e.contact_email,
		e.website,
		e.host_user_id
	FROM events e
	JOIN rsvps r ON r.event_id = e.event_id
	WHERE r.user_id = ?
	ORDER BY e.date_time, e.event_id
`, userID)
```

## 9. Use transactions when multiple operations must succeed together

A transaction is useful when an operation involves multiple statements:

```go
tx, err := db.Begin()
if err != nil {
	return err
}
defer tx.Rollback()

// Perform tx.Exec or tx.QueryRow operations here.

if err := tx.Commit(); err != nil {
	return err
}
```

For a simple RSVP insert, the foreign keys and primary key are enough. For more complicated operations, use a transaction.

Example:

```go
func CreateUserWithInitialData(db *sql.DB, name, email string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO users (user_name, user_email)
		VALUES (?, ?)
	`, name, email)
	if err != nil {
		return err
	}

	// Other dependent operations could go here.

	return tx.Commit()
}
```

## 10. How this would fit your project

A reasonable backend structure would be:

```text
backend/
    Database/
        Database.go
        Schema.go
    Repository/
        Users.go
        Events.go
        RSVPs.go
    Controller/
        API.go
        Users.go
        Events.go
        Rsvps.go
        Types.go
    main.go
```

The flow would become:

```text
HTTP controller
    -> repository method
        -> database/sql
            -> SQLite
```

For example:

```go
type API struct {
	db *sql.DB
}

func NewAPI(db *sql.DB) *API {
	return &API{db: db}
}
```

Then the controller would no longer call:

```go
api.campus.GetUser(...)
api.campus.CreateEvent(...)
user.RsvpAdd(event)
```

Instead, it would call repository functions such as:

```go
api.users.GetByID(...)
api.events.Create(...)
api.rsvps.Add(...)
```

The `sync.RWMutex` currently used by `API` would likely be removed. `*sql.DB` is safe for concurrent use and manages a pool of database connections.

The main thing to remember is that `*sql.DB` should be created once at startup and shared throughout the application. Do not open and close a new database connection for every request.