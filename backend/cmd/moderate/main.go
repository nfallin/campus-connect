package main

import (
	"campus_connect/Model"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	list := flag.Bool("list", false, "list pending events")
	approve := flag.Int("approve", 0, "approve a pending event by ID")
	deny := flag.Int("deny", 0, "deny a pending event by ID")
	dbPath := flag.String("db", "campus_connect.db", "path to the SQLite database")
	flag.Parse()

	operations := 0
	if *list {
		operations++
	}
	if *approve > 0 {
		operations++
	}
	if *deny > 0 {
		operations++
	}
	if operations != 1 {
		fmt.Fprintln(os.Stderr, "exactly one of -list, -approve ID, or -deny ID is required")
		os.Exit(2)
	}

	db, err := Model.Init(*dbPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer db.Close()
	repository := Model.NewRepository(db)

	switch {
	case *list:
		events, err := repository.GetPendingEvents()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, event := range events {
			fmt.Printf(
				"{\n\tID: %d\n\tPending: %t\n\tName: %s\n\tDateTime: %s\n\tEndDateTime: %v\n\t"+
					"Location: %s\n\tClubName: %s\n\tDescription: %s\n\tHost: %+v\n\t"+
					"ContactEmail: %s\n\tWebsiteURL: %s\n\tAttendeeCount: %d\n\tHostUserID: %d\n}\n",
				event.ID, event.Pending, event.Name, event.DateTime.Format(time.RFC3339),
				event.EndDateTime, event.Location, event.ClubName, event.Description,
				event.Host, event.ContactEmail, event.WebsiteURL, event.AttendeeCount,
				event.HostUserID,
			)
		}
	case *approve > 0:
		if err := repository.ApproveEvent(*approve); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("approved event %d\n", *approve)
	case *deny > 0:
		if err := repository.DeleteEvent(*deny); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("denied event %d\n", *deny)
	}
}
