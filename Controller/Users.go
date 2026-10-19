package Controller

import (
	"campus_connect/Model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (api *API) createUser(c *gin.Context) {
	var request CreateUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	newUser := api.campus.GetUserByName(request.Name)
	status := http.StatusOK

	if newUser == nil {
		newUser = api.campus.CreateUser(request.Name)
		status = http.StatusCreated
	}

	c.JSON(status, toUserResponse(newUser))
}

func (api *API) getUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "user id must be an integer",
		})
		return
	}

	user := api.campus.GetUser(userID)
	if user == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Message: "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, toUserResponse(user))
}

func toUserResponse(user *Model.User) UserResponse {
	return UserResponse{
		ID:   user.GetID(),
		Name: user.GetName(),
	}
}
