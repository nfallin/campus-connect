package main

import (
	"campus_connect/Controller"
	"campus_connect/Database"
)

func main() {

	db, err := Database.Init("campus_connect.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repository := Database.NewRepository(db)
	if err := Controller.Serve(repository); err != nil {
		panic(err)
	}
}
