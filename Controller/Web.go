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
func Serve() error {
	router := newRouter()
	return router.Run("localhost:8080")
}

func newRouter() *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	api := NewAPI(Model.NewCampus())
	defineRoutes(router, api)

	return router
}

// define get/post routes here
func defineRoutes(r *gin.Engine, api *API) {
	// define read/write routes - preventing race conditions when multiple API calls are accessing the same data
	// this will stop being necessary when we implement a database
	readRoutes := r.Group("")
	readRoutes.Use(api.readLock())

	writeRoutes := r.Group("")
	writeRoutes.Use(api.writeLock())

	writeRoutes.POST("/events", api.createEvent)
	readRoutes.GET("/events/:id", api.getEvent)

	writeRoutes.POST("/users", api.createUser)
	readRoutes.GET("/users/:id", api.getUser)

	writeRoutes.POST("/users/:userId/rsvps/:eventId/add", api.rsvpAdd)
	writeRoutes.POST("/users/:userId/rsvps/:eventId/remove", api.rsvpRemove)
	readRoutes.POST("/users/:userId/rsvps", api.getRSVPs)
}
