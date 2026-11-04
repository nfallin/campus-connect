# System Specifications

## Functional Specifications
- The frontend must allow the user to enter an event title, date/time, location, club name, contact email, and event website in separate input fields.
- The frontend must allow the user to create new events in the backend according to the information they enter.
- The frontend must refuse to allow the user to create an event unless its title, date/time, location, description, club name, and host user ID fields are non-empty.
- The frontend must refuse to allow the user to create an event whose date/time is not a valid timestamp.
- The frontend must refuse to allow the user to create an event whose date/time is in the past or over 6 months into the future.
- The frontend must refuse to allow the user to create an event whose contact email does not match the conventional email format X@Y.Z.
- The frontend must refuse to allow the user to create an event whose website link does not match a valid URL.
- The frontend must provide functional hyperlinks for events containing website links.
- The frontend must refuse to allow the user to create an event unless they have signed into an account
- The frontend must refuse to allow the user to RSVP to an event unless they have signed into an account
- The frontend must refuse to allow a user to edit an event unless they are the event's creator
- The frontend must display an error message if the user attempts to create or rename an event using a name that is already in use
- Event names must be unique.
- User name/email pairs must be unique.
- The API must reject any requests that do not match the request formats outlined in `API.md`.
- The API must reject event creation requests that do not contain a host ID identifying an existing user.
- The frontend must display all events successfully retrieved from the backend, including their title, date, time, location, description, club name, contact email, and event website.

## Non-Functional Specifications
- The frontend must obtain the list of stored events within 3 seconds of initialization as long as the backend is running and reachable.
- The frontend UI layout must be accessable and usable on a mobile device 
- The backend must create and store event objects when supplied with a valid request as outlined in `API.md` within 3 seconds.
- The backend must return the appropriate response to any valid request as outlined in `API.md` within 3 seconds.
- The frontend must complete sign-in and RSVP actions within 3 seconds as long as the backend is running and reachable.