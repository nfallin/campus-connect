package Controller

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)

// *****************
// * API Functions *
// *****************

// curl command to test:
// curl -i -X GET localhost:8080/
func defaultGet(c *gin.Context) {
	var defaultResponse = []string{
		"This is the default API response",
	}

	// return a JSON
	c.IndentedJSON(http.StatusOK, defaultResponse)
}

// curl command to test:
// curl -i -X POST -H "Content-Type: application/json" -d "{\"message\": \"Hello World!\"}" localhost:8080/post
func defaultPost(c *gin.Context) {
	// 1. declare type designed to hold requests made with this function
	var requestBody defaultPostBody

	// 2. check for errors and return an error message
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. handle request using types and return a success message
	fmt.Println(requestBody.Message)
	c.JSON(http.StatusOK, gin.H{"status": "message received succsessfully"})
}
