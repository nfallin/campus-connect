package Model

/*
	TestNewCampus	Initializes empty user, club, and event maps with IDs starting at 1
	TestGetNextUserID	Generates sequential user IDs
	TestGetNextClubID	Generates sequential club IDs
	TestGetNextEventID	Generates sequential event IDs
	TestGetUsers	Returns all campus users
	TestGetUser	Returns an existing user and nil for an unknown ID
	TestGetClubs	Returns all campus clubs
	TestGetClub	Returns an existing club and nil for an unknown ID
	TestGetEvents	Returns all campus events
	TestGetEvent	Returns an existing event and nil for an unknown ID
	TestCreateUser	Creates a user, assigns an ID, initializes its RSVP list, and adds it to campus
	TestAddUser	Adds a user to the campus user map
	TestCreateClub	Creates a club, assigns an ID, and adds it to campus
	TestAddClub	Adds a club to the campus club map
	TestCreateEvent	Creates an event, assigns an ID, initializes attendees, stores its fields, and adds it to campus
	TestAddEvent	Adds an event to the campus event map
	TestCreateMultipleEntities	Assigns independent sequential IDs for users, clubs, and events
	TestRemoveEvent	Removes an event from campus and cleans up attending users
	TestRemoveClub	Removes a club and all of its events while preserving unrelated events
	TestRemoveUser	Removes a user and cleans up the user's RSVPs from events
*/