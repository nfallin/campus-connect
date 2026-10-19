package Controller

import (
	"campus_connect/Model"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (api *API) getEvents(c *gin.Context) {
	events := api.campus.GetEvents()
	eventResponses := make([]EventResponse, 0, len(events))

	for _, event := range events {
		eventResponses = append(eventResponses, toEventResponse(event))
	}

	sort.Slice(eventResponses, func(i, j int) bool {
		if eventResponses[i].DateTime.Equal(eventResponses[j].DateTime) {
			return eventResponses[i].ID < eventResponses[j].ID
		}
		return eventResponses[i].DateTime.Before(eventResponses[j].DateTime)
	})

	c.JSON(http.StatusOK, gin.H{"events": eventResponses})
}

func (api *API) createEvent(c *gin.Context) {
	var request CreateEventRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	event := api.campus.CreateEvent(
		request.Name,
		request.DateTime,
		request.Location,
		request.ClubName,
	)

	c.JSON(http.StatusCreated, toEventResponse(event))
}

func (api *API) getEvent(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "event id must be an integer",
		})
		return
	}

	event := api.campus.GetEvent(eventID)
	if event == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Message: "event not found",
		})
		return
	}

	c.JSON(http.StatusOK, toEventResponse(event))
}

func toEventResponse(event *Model.Event) EventResponse {
	return EventResponse{
		ID:            event.GetID(),
		Name:          event.GetName(),
		DateTime:      event.GetDateTime(),
		Location:      event.GetLocation(),
		ClubName:      event.GetClubName(),
		AttendeeCount: event.GetNumAttendees(),
	}
}
