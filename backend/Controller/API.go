package Controller

import (
	"campus_connect/Model"
)

type API struct {
	repository *Model.Repository
}

func NewAPI(repository *Model.Repository) *API {
	return &API{
		repository: repository,
	}
}
