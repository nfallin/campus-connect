package Controller

import (
	"campus_connect/Model"

	"github.com/gin-gonic/gin"
)

// ***************************
// * POST request structures *
// ***************************
type defaultPostBody struct {
	Message string `json:"message"`
}

// starts the router
func Serve(campus *Model.Campus) error {
	router := newRouter(campus)
	return router.Run("localhost:8080")
}

func newRouter(campus *Model.Campus) *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	defineRoutes(router, campus)
	return router
}

// define get/post routes here
func defineRoutes(r *gin.Engine, c *Model.Campus) {
	// api routes
	r.GET("/", defaultGet)
	r.POST("/post", defaultPost)
}
