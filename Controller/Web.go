package Controller

import (
	"github.com/gin-gonic/gin"
)

// ***************************
// * POST request structures *
// ***************************
type defaultPostBody struct {
	Message string `json:"message"`
}

// starts the router
func Serve() {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	defineRoutes(router)
	router.Run("localhost:8080")
}

// define get/post routes here
func defineRoutes(r *gin.Engine) {
	// api routes
	r.GET("/", defaultGet)
	r.POST("/post", defaultPost)
}
