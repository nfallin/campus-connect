package main

import (
	"campus_connect/Controller"
	"campus_connect/Model"
)

func main() {

	c := Model.NewCampus()

	Controller.Serve()
}
