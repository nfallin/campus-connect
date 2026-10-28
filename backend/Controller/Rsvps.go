package Controller

import (
	"campus_connect/Model"
	"campus_connect/Types"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func parseRSVPIDs(c *gin.Context) (int, int, bool) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{Message: "user id must be an integer"})
		return 0, 0, false
	}
	eventID, err := strconv.Atoi(c.Param("eventId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{Message: "event id must be an integer"})
		return 0, 0, false
	}
	return userID, eventID, true
}

func (api *API) rsvpAdd(c *gin.Context) {
	userID, eventID, ok := parseRSVPIDs(c)
	if !ok {
		return
	}
	err := api.repository.AddRSVP(userID, eventID)
	if errors.Is(err, Model.ErrNotFound) {
		c.JSON(http.StatusNotFound, Types.ErrorResponse{Message: "user or event not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, Types.ErrorResponse{Message: err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (api *API) rsvpRemove(c *gin.Context) {
	userID, eventID, ok := parseRSVPIDs(c)
	if !ok {
		return
	}
	err := api.repository.RemoveRSVP(userID, eventID)
	if errors.Is(err, Model.ErrNotFound) {
		c.JSON(http.StatusNotFound, Types.ErrorResponse{Message: "user or event not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, Types.ErrorResponse{Message: err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (api *API) getRSVPs(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{Message: "user id must be an integer"})
		return
	}
	events, err := api.repository.GetUserRSVPs(userID)
	if errors.Is(err, Model.ErrNotFound) {
		c.JSON(http.StatusNotFound, Types.ErrorResponse{Message: "user not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, Types.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"events": events})
}
