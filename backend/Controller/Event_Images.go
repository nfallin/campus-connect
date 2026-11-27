package Controller

import (
	"bytes"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"strconv"

	"campus_connect/Model"
	"campus_connect/Types"

	"github.com/gin-gonic/gin"
)

const (
	maxEventImageSize    = 5 * 1024 * 1024 // 5 MB
	maxImageRequestSize  = maxEventImageSize + 1024*1024
	maxEventImagePixels  = 25_000_000
)

func (api *API) uploadEventImage(c *gin.Context) {
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
		c.JSON(http.StatusForbidden, Types.ErrorResponse{Message: "only the event host can upload an image"})
		return
	}

	// Limit the total request size before parsing the uploaded file.
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxImageRequestSize)

	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{Message: "a valid image file is required"})
		return
	}

	if fileHeader.Size > maxEventImageSize {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{Message: "image must be 5 MB or smaller"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{Message: "could not read uploaded image"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, maxEventImageSize+1))
	if err != nil {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{Message: "could not read uploaded image"})
		return
	}

	if len(data) == 0 || len(data) > maxEventImageSize {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{Message: "image must be between 1 byte and 5 MB"})
		return
	}

	contentType := http.DetectContentType(data)
	if contentType != "image/jpeg" && contentType != "image/png" {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{Message: "only JPEG and PNG images are supported"})
		return
	}

	// Verify that the file is a decodable image.
	config, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || config.Width <= 0 || config.Height <= 0 ||
		int64(config.Width)*int64(config.Height) > maxEventImagePixels {
		c.JSON(http.StatusBadRequest, Types.ErrorResponse{Message: "image file is invalid or has dimensions that are too large"})
		return
	}

	if err := api.repository.SaveEventImage(eventID, data, contentType); err != nil {
		if errors.Is(err, Model.ErrNotFound) {
			c.JSON(http.StatusNotFound, Types.ErrorResponse{Message: "event not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, Types.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event image uploaded"})
}

func (api *API) getEventImage(c *gin.Context) {
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

	isOwner := false
	if user, exists := c.Get(authenticatedUserKey); exists {
		isOwner = event.HostUserID == user.(Types.UserResponse).ID
	}

	if event.Pending && !isOwner {
		c.JSON(http.StatusNotFound, Types.ErrorResponse{Message: "event not found"})
		return
	}

	data, contentType, err := api.repository.GetEventImage(eventID)
	if errors.Is(err, Model.ErrNotFound) {
		c.JSON(http.StatusNotFound, Types.ErrorResponse{Message: "event image not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, Types.ErrorResponse{Message: err.Error()})
		return
	}

	c.Header("X-Content-Type-Options", "nosniff")
	c.Data(http.StatusOK, contentType, data)
}

func (api *API) deleteEventImage(c *gin.Context) {
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
		c.JSON(http.StatusForbidden, Types.ErrorResponse{Message: "only the event host can delete an image"})
		return
	}

	if err := api.repository.DeleteEventImage(eventID); err != nil {
		if errors.Is(err, Model.ErrNotFound) {
			c.JSON(http.StatusNotFound, Types.ErrorResponse{Message: "event not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, Types.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event image deleted"})
}
