package main

import (
	"campus_connect/Controller"
	"campus_connect/Model"
)

func main() {

	db, err := Model.Init("campus_connect.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repository := Model.NewRepository(db)
	if err := Controller.Serve(repository); err != nil {
		panic(err)
	}
}
