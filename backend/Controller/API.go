package Controller

import (
	"campus_connect/Database"
)

type API struct {
	repository *Database.Repository
}

func NewAPI(repository *Database.Repository) *API {
	return &API{
		repository: repository,
	}
}
