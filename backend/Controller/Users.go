package Controller

import (
	"campus_connect/Model"
	"campus_connect/Types"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const authenticatedUserKey = "authenticatedUser"

func (api *API) registerUser(c *gin.Context) {
	var request Types.RegisterUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	_, err := api.repository.RegisterUser(request)
	if errors.Is(err, Model.ErrDuplicateUsername) {
		c.JSON(http.StatusConflict, Types.ErrorResponse{Message: err.Error()})
		return
	}
	if errors.Is(err, Model.ErrDuplicateEmail) {
		c.JSON(http.StatusConflict, Types.ErrorResponse{Message: err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, Types.ErrorResponse{Message: err.Error()})
		return
	}

	auth, err := api.repository.LoginUser(Types.LoginUserRequest{
		Identifier: request.Username,
		Password:   request.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, Types.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, auth)
}

func (api *API) loginUser(c *gin.Context) {
	var request Types.LoginUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{Message: err.Error()})
		return
	}

	auth, err := api.repository.LoginUser(request)
	if errors.Is(err, Model.ErrInvalidCredentials) {
		c.JSON(http.StatusUnauthorized, Types.ErrorResponse{Message: err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, Types.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, auth)
}

func (api *API) logoutUser(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func (api *API) getCurrentUser(c *gin.Context) {
	c.JSON(http.StatusOK, c.MustGet(authenticatedUserKey))
}

func (api *API) requireAuth(c *gin.Context) {
	user, err := api.repository.AuthenticateToken(bearerToken(c))
	if errors.Is(err, Model.ErrUnauthorized) {
		c.JSON(http.StatusUnauthorized, Types.ErrorResponse{Message: err.Error()})
		c.Abort()
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, Types.ErrorResponse{Message: err.Error()})
		c.Abort()
		return
	}
	c.Set(authenticatedUserKey, user)
	c.Next()
}

func bearerToken(c *gin.Context) string {
	const prefix = "Bearer "
	header := c.GetHeader("Authorization")
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}

func authenticatedUser(c *gin.Context) Types.UserResponse {
	user, _ := c.Get(authenticatedUserKey)
	return user.(Types.UserResponse)
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
