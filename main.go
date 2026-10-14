package main

import (
	"campus_connect/Model"
	"campus_connect/Web"
	"fmt"
	"time"
)

func main() {

	c := Model.NewCampus()

	// three clubs
	gameDevClub := c.CreateClub("Game Development Club")
	roboticsClub := c.CreateClub("Robotics Club")
	photographyClub := c.CreateClub("Photography Club")

	// five users
	alice := c.CreateUser("Alice")
	bob := c.CreateUser("Bob")
	charlie := c.CreateUser("Charlie")
	diana := c.CreateUser("Diana")
	ethan := c.CreateUser("Ethan")

	// four events
	gameNight := c.CreateEvent("Game Night", time.Now().AddDate(0, 0, 3), "Room 101", nil, gameDevClub.GetID())
	robotWorkshop := c.CreateEvent("Robot Workshop", time.Now().AddDate(0, 0, 5), "Room 202", nil, roboticsClub.GetID())
	photoWalk := c.CreateEvent("Photo Walk", time.Now().AddDate(0, 0, 7), "Campus Grounds", nil, photographyClub.GetID())
	campusSocial := c.CreateEvent("Campus Social", time.Now().AddDate(0, 0, 10), "Room 101", nil, gameDevClub.GetID())

	// everyone RSVPs to some events
	alice.RsvpAdd(gameNight)
	alice.RsvpAdd(photoWalk)
	bob.RsvpAdd(gameNight)
	bob.RsvpAdd(robotWorkshop)
	charlie.RsvpAdd(gameNight)
	charlie.RsvpAdd(campusSocial)
	diana.RsvpAdd(robotWorkshop)
	diana.RsvpAdd(photoWalk)
	ethan.RsvpAdd(robotWorkshop)
	ethan.RsvpAdd(campusSocial)

	// duplicate and cancel invalid rsvps
	fmt.Printf("Game Night attendees before alice tries to rsvp a second time: %v\n", gameNight.GetAttendees())
	alice.RsvpAdd(gameNight)
	fmt.Printf("Game Night attendees after alice tries to rsvp a second time and before ethan tries to remove his rsvp: %v\n", gameNight.GetAttendees())
	ethan.RsvpRemove(gameNight)
	fmt.Printf("Game Night attendees: %v\n", gameNight.GetAttendees())

	fmt.Printf("Bob's rsvp list before the robot workshop is deleted: %v\n", bob.GetRsvpList())
	c.RemoveEvent(robotWorkshop)
	fmt.Printf("Bob's rsvp list after the robot workshop is deleted: %v\n", bob.GetRsvpList())

	fmt.Printf("photo walk atendees before diana is deleted: %v\n", photoWalk.GetAttendees())
	c.RemoveUser(diana)
	fmt.Printf("photo walk atendees after diana is deleted: %v\n", photoWalk.GetAttendees())

	fmt.Printf("ethan's rsvp list before the game dev club is deleted: %v\n", ethan.GetRsvpList())
	c.RemoveClub(gameDevClub)
	fmt.Printf("ethan's rsvp list after the game dev club is deleted: %v\n", ethan.GetRsvpList())

	Web.Serve()
}
