package Controller

import (
	"campus_connect/Model"
	"campus_connect/Types"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (api *API) createUser(c *gin.Context) {
	var request Types.CreateUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	user, created, err := api.repository.CreateUser(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Types.ErrorResponse{Message: err.Error()})
		return
	}

	status := http.StatusOK
	if created {
		status = http.StatusCreated
	}

	c.JSON(status, gin.H{
		"user":      user,
		"isNewUser": created,
	})
}

func (api *API) getUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userId"))

	if err != nil {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{
			Message: "user id must be an integer",
		})
		return
	}

	user, err := api.repository.GetUser(userID)
	if errors.Is(err, Model.ErrInvalidID) {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{Message: err.Error()})
		return
	}
	if errors.Is(err, Model.ErrNotFound) {
		c.JSON(http.StatusNotFound, Types.ErrorResponse{
			Message: "user not found",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, Types.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
