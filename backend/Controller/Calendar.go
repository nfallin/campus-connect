package Controller

import (
	"campus_connect/Model"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"campus_connect/Types"

	"github.com/gin-gonic/gin"
)

func (api *API) downloadEventCalendar(c *gin.Context) {
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

	if event.Pending && !isEventOwner(c, event) {
		c.JSON(http.StatusNotFound, Types.ErrorResponse{Message: "event not found"})
		return
	}

	c.Header("Content-Type", "text/calendar; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="event-%d.ics"`, event.ID))
	c.String(http.StatusOK, eventCalendar(event))
}

func isEventOwner(c *gin.Context, event Types.EventResponse) bool {
	user, exists := c.Get(authenticatedUserKey)
	return exists && event.HostUserID == user.(Types.UserResponse).ID
}

func eventCalendar(event Types.EventResponse) string {
	endDateTime := event.DateTime.Add(time.Hour)
	if event.EndDateTime != nil {
		endDateTime = *event.EndDateTime
	}

	lines := []string{
		"BEGIN:VCALENDAR",
		"VERSION:2.0",
		"PRODID:-//Campus Connect//Campus Events//EN",
		"CALSCALE:GREGORIAN",
		"METHOD:PUBLISH",
		"BEGIN:VEVENT",
		"UID:" + fmt.Sprintf("event-%d@campus-connect", event.ID),
		"DTSTAMP:" + time.Now().UTC().Format("20060102T150405Z"),
		"DTSTART:" + event.DateTime.UTC().Format("20060102T150405Z"),
		"DTEND:" + endDateTime.UTC().Format("20060102T150405Z"),
		"SUMMARY:" + escapeCalendarText(event.Name),
		"DESCRIPTION:" + escapeCalendarText(event.Description),
		"LOCATION:" + escapeCalendarText(event.Location),
	}
	if event.WebsiteURL != "" {
		lines = append(lines, "URL:"+strings.ReplaceAll(strings.ReplaceAll(event.WebsiteURL, "\r", ""), "\n", ""))
	}
	lines = append(lines, "END:VEVENT", "END:VCALENDAR")

	return strings.Join(lines, "\r\n") + "\r\n"
}

func escapeCalendarText(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, ";", "\\;")
	value = strings.ReplaceAll(value, ",", "\\,")
	value = strings.ReplaceAll(value, "\r\n", "\\n")
	value = strings.ReplaceAll(value, "\r", "\\n")
	return strings.ReplaceAll(value, "\n", "\\n")
}
