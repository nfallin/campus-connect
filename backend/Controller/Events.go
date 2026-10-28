package Controller

import (
	"campus_connect/Model"
	"campus_connect/Types"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (api *API) getEvents(c *gin.Context) {
	events, err := api.repository.GetEvents(Types.EventFilters{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, Types.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"events": events})
}

func (api *API) createEvent(c *gin.Context) {
	var request Types.CreateEventRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{Message: err.Error()})
		return
	}

	event, err := api.repository.CreateEvent(request)
	if errors.Is(err, Model.ErrInvalidEvent) {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{Message: err.Error()})
		return
	}
	if errors.Is(err, Model.ErrNotFound) {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{Message: "host user not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, Types.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, event)
}

func (api *API) updateEvent(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{Message: "event id must be an integer"})
		return
	}

	var request Types.UpdateEventRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{Message: err.Error()})
		return
	}

	event, err := api.repository.UpdateEvent(eventID, request)
	if errors.Is(err, Model.ErrInvalidID) {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{Message: err.Error()})
		return
	}
	if errors.Is(err, Model.ErrInvalidEvent) {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{Message: err.Error()})
		return
	}
	if errors.Is(err, Model.ErrNotFound) {
		c.JSON(http.StatusNotFound, Types.ErrorResponse{Message: "event not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, Types.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, event)
}

func (api *API) getEvent(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{Message: "event id must be an integer"})
		return
	}

	event, err := api.repository.GetEvent(eventID)
	if errors.Is(err, Model.ErrInvalidID) {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{Message: err.Error()})
		return
	}
	if errors.Is(err, Model.ErrNotFound) {
		c.JSON(http.StatusNotFound, Types.ErrorResponse{Message: "event not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, Types.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, event)
}
