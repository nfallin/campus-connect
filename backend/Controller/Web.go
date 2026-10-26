package Controller

import (
	"campus_connect/Model"

	"github.com/gin-gonic/gin"
)

// starts the router
func Serve() error {
	router := newRouter()
	return router.Run("0.0.0.0:8080")
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

	readRoutes.GET("/events", api.getEvents)
	writeRoutes.POST("/events", api.createEvent)
	writeRoutes.PATCH("/events/:id", api.updateEvent)
	readRoutes.GET("/events/:id", api.getEvent)

	writeRoutes.POST("/users", api.createUser)
	readRoutes.GET("/users/:userId", api.getUser)

	writeRoutes.POST("/users/:userId/rsvps/:eventId", api.rsvpAdd)
	writeRoutes.DELETE("/users/:userId/rsvps/:eventId", api.rsvpRemove)
	readRoutes.GET("/users/:userId/rsvps", api.getRSVPs)
}
