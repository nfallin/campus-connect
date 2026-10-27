package Database

// users cannot have the same name and email address, but they can have the same name.
// they are differentiated by the PAIR of values.
const usersTable = `
	CREATE TABLE IF NOT EXISTS users (
		user_id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_name TEXT NOT NULL,
		user_email TEXT NOT NULL,
		UNIQUE(user_name, user_email)
	);
`

// users cannot be deleted while they host events
const eventsTable = `
	CREATE TABLE IF NOT EXISTS events (
		event_id INTEGER PRIMARY KEY AUTOINCREMENT,
		event_name TEXT NOT NULL,
		date_time TEXT NOT NULL,
		location TEXT NOT NULL,
		club_name TEXT NOT NULL,
		contact_email TEXT,
		website TEXT,
		host_user_id INTEGER NOT NULL,
		FOREIGN KEY(host_user_id) REFERENCES users(user_id) ON DELETE RESTRICT
	);
`

// rsvps are deleted when their user or host are deleted.
// rsvps don't need their own id
const rsvpTable = `
	CREATE TABLE IF NOT EXISTS rsvps (
		user_id INTEGER NOT NULL,
		event_id INTEGER NOT NULL,
		PRIMARY KEY(user_id, event_id),
		FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE,
		FOREIGN KEY(event_id) REFERENCES events(event_id) ON DELETE CASCADE
	);
`