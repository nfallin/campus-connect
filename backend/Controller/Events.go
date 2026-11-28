package Controller

import (
	"campus_connect/Model"
	"campus_connect/Types"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (api *API) getEventInterests(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"interests": Types.AllEventInterests()})
}

func (api *API) getEvents(c *gin.Context) {
	filters := Types.EventFilters{}
	if user, exists := c.Get(authenticatedUserKey); exists {
		viewer := user.(Types.UserResponse).ID
		filters.ViewerUserID = &viewer
	}
	events, err := api.repository.GetEvents(filters)
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
	event, err := api.repository.CreateEvent(request, authenticatedUserID(c))
	if errors.Is(err, Model.ErrInvalidEvent) {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{Message: err.Error()})
		return
	}
	if errors.Is(err, Model.ErrDuplicateEventName) {
		c.JSON(http.StatusConflict, Types.ErrorResponse{Message: err.Error()})
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

	event, err := api.repository.GetEvent(eventID)
	if errors.Is(err, Model.ErrNotFound) {
		c.JSON(http.StatusNotFound, Types.ErrorResponse{Message: "event not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, Types.ErrorResponse{Message: err.Error()})
		return
	}
	if event.HostUserID != authenticatedUserID(c) {
		c.JSON(http.StatusForbidden, Types.ErrorResponse{Message: "only the event host can update this event"})
		return
	}

	var request Types.UpdateEventRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{Message: err.Error()})
		return
	}
	event, err = api.repository.UpdateEvent(eventID, request)
	if errors.Is(err, Model.ErrInvalidID) {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{Message: err.Error()})
		return
	}
	if errors.Is(err, Model.ErrInvalidEvent) {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{Message: err.Error()})
		return
	}
	if errors.Is(err, Model.ErrDuplicateEventName) {
		c.JSON(http.StatusConflict, Types.ErrorResponse{Message: err.Error()})
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

func (api *API) deleteEvent(c *gin.Context) {
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
	if event.HostUserID != authenticatedUserID(c) {
		c.JSON(http.StatusForbidden, Types.ErrorResponse{Message: "only the event host can cancel this event"})
		return
	}

	if err := api.repository.DeleteEvent(eventID); errors.Is(err, Model.ErrInvalidID) {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{Message: err.Error()})
		return
	} else if errors.Is(err, Model.ErrNotFound) {
		c.JSON(http.StatusNotFound, Types.ErrorResponse{Message: "event not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, Types.ErrorResponse{Message: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
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
	isOwner := false
	if user, exists := c.Get(authenticatedUserKey); exists {
		isOwner = event.HostUserID == user.(Types.UserResponse).ID
	}
	if event.Pending && !isOwner {
		c.JSON(http.StatusNotFound, Types.ErrorResponse{Message: "event not found"})
		return
	}
	c.JSON(http.StatusOK, event)
}
