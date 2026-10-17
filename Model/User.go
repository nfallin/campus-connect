package Model

type User struct {
	name     string
	id       int
	rsvpList map[int]struct{}
}

func (u User) GetName() string {
	return u.name
}

// returns the list of events that u has RSVP'd to
func (u User) GetRsvpList() map[int]struct{} {
	return u.rsvpList
}

// add e to u's rsvpList and u to e's attendees list
func (u *User) RsvpAdd(e *Event) {
	u.rsvpList[e.id] = struct{}{}
	e.attendees[u.id] = struct{}{}
}

// remove e from u's rsvpList and u from e's attendees list
func (u *User) RsvpRemove(e *Event) {
	delete(u.rsvpList, e.id)
	delete(e.attendees, u.id)
}