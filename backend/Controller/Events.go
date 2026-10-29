package Controller

import (
	"campus_connect/Model"
	"net/http"
	"sort"
	"strconv"
	"strings"

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
	if strings.TrimSpace(request.Description) == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "description is required"})
		return
	}

	host := api.campus.GetUser(*request.HostUserID)
	if host == nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "host user not found"})
		return
	}
	event := api.campus.CreateEvent(
		request.Name,
		request.DateTime,
		request.Location,
		request.ClubName,
		host,
		request.ContactEmail,
		request.WebsiteURL,
		request.Description,
	)

	c.JSON(http.StatusCreated, toEventResponse(event))
}

func (api *API) updateEvent(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "event id must be an integer"})
		return
	}

	event := api.campus.GetEvent(eventID)
	if event == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Message: "event not found"})
		return
	}

	var request UpdateEventRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	name := event.GetName()
	if request.Name != nil {
		name = *request.Name
	}
	description := event.GetDescription()
	if request.Description != nil {
		if strings.TrimSpace(*request.Description) == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: "description is required"})
			return
		}
		description = *request.Description
	}
	dateTime := event.GetDateTime()
	if request.DateTime != nil {
		dateTime = *request.DateTime
	}
	location := event.GetLocation()
	if request.Location != nil {
		location = *request.Location
	}
	clubName := event.GetClubName()
	if request.ClubName != nil {
		clubName = *request.ClubName
	}
	host := event.GetHost()
	if request.HostUserID != nil {
		host = api.campus.GetUser(*request.HostUserID)
		if host == nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Message: "host user not found"})
			return
		}
	}
	contactEmail := event.GetContactEmail()
	if request.ContactEmail != nil {
		contactEmail = *request.ContactEmail
	}
	websiteURL := event.GetWebsiteURL()
	if request.WebsiteURL != nil {
		websiteURL = *request.WebsiteURL
	}

	api.campus.UpdateEvent(event, name, dateTime, location, clubName, host, contactEmail, websiteURL, description)
	c.JSON(http.StatusOK, toEventResponse(event))
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
	var host *UserResponse
	if event.GetHost() != nil {
		hostResponse := toUserResponse(event.GetHost())
		host = &hostResponse
	}

	return EventResponse{
		ID:            event.GetID(),
		Name:          event.GetName(),
		Description:   event.GetDescription(),
		DateTime:      event.GetDateTime(),
		Location:      event.GetLocation(),
		ClubName:      event.GetClubName(),
		Host:          host,
		ContactEmail:  event.GetContactEmail(),
		WebsiteURL:    event.GetWebsiteURL(),
		AttendeeCount: event.GetNumAttendees(),
	}
}
