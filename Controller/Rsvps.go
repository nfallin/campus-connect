package Controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (api *API) rsvpAdd(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "user id must be an integer",
		})
		return
	}

	eventId, err := strconv.Atoi(c.Param("eventId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "event id must be an integer",
		})
		return
	}

	user := api.campus.GetUser(userId)
	if user == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Message: "user not found",
		})
		return
	}

	event := api.campus.GetEvent(eventId)
	if event == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Message: "event not found",
		})
		return
	}

	user.RsvpAdd(event)

	c.Status(http.StatusOK)
}

func (api *API) rsvpRemove(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "user id must be an integer",
		})
		return
	}

	eventId, err := strconv.Atoi(c.Param("eventId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "event id must be an integer",
		})
		return
	}

	user := api.campus.GetUser(userId)
	if user == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Message: "user not found",
		})
		return
	}

	event := api.campus.GetEvent(eventId)
	if event == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Message: "event not found",
		})
		return
	}

	user.RsvpRemove(event)

	c.Status(http.StatusOK)
}

func (api *API) getRSVPs(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "user id must be an integer",
		})
		return
	}

	user := api.campus.GetUser(userId)
	if user == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Message: "user not found",
		})
		return
	}

	rsvps := user.GetRsvpList()
	events := make([]EventResponse, 0, len(rsvps))
	for eventID := range rsvps {
		event := api.campus.GetEvent(eventID)
		if event == nil {
			continue
		}

		events = append(events, toEventResponse(event))
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}
