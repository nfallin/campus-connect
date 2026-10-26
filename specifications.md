# System Specifications

## Functional Specifications
- The frontend must allow the user to enter an event title, date/time, location, location, clube name, contact email, and event website in separate interactable fields.
- The frotnend must allow the user to create new events in the backend according to the information they enter.
- The frontend must refuse to allow the user to create an event unless its title, date/time, location, club name, and host user ID fields are non-empty.
- The frontend must refuse to allow the user to create an event whose date/time is not a valid timestamp.
- The frontend must refuse to allow the user to create an event whose contact email does not match the conventional email format X@Y.Z.
- The frontend must refuse to allow the user to create an event whose website link does not match a valid URL.
- The frontend must provide functional hyperlinks for events containing website links.
- The API must reject any requests that do not match the request formats outlined in `API.md`.
- The API must reject event creation requests that do not contain a host ID identifying an existing user.

## Non-Functional Specifications
- The frontend must obtain the list of stored events within 3 seconds of initialization as long as the backend is running and reachable. 
- The backend must create and store event objects when supplied with a valid request as outlined in `API.md` within 3 seconds.
- The backend must return the appropriate response to any valid request as outlined in `API.md` within 3 seconds.