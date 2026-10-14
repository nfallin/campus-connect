package Model

/*
	TestCreateEntitiesAndLookups	Created users, clubs, and events are reachable through campus lookup methods
	TestCreateEventAndRSVP	A user can RSVP to a campus event
	TestRSVPBidirectionalConsistency	RSVP appears in both the User and Event
	TestRemoveRSVPBidirectionalConsistency	Removing RSVP removes it from both sides
	TestRemoveEventCleansUpRSVPs	Deleting an event removes it from campus and every user's RSVP list
	TestRemoveUserCleansUpRSVPs	Deleting a user removes that user from every event's attendees
	TestRemoveClubCleansUpEventsAndRSVPs	Deleting a club removes its events and cleans up their RSVPs
	TestMultipleUsersRSVP	Multiple users can RSVP to the same event
	TestMultipleUsersAndEvents	Users can RSVP to different combinations of events without affecting each other
*/