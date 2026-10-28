package Controller

import (
	"campus_connect/Database"

	"github.com/gin-gonic/gin"
)

// starts the router
func Serve(repository *Database.Repository) error {
	router := newRouter(repository)
	return router.Run("0.0.0.0:8080")
}

func newRouter(repository *Database.Repository) *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	api := NewAPI(repository)
	defineRoutes(router, api)

	return router
}

// define get/post routes here
func defineRoutes(r *gin.Engine, api *API) {
	r.GET("/events", api.getEvents)
	r.POST("/events", api.createEvent)
	r.PATCH("/events/:id", api.updateEvent)
	r.GET("/events/:id", api.getEvent)

	r.POST("/users", api.createUser)
	r.GET("/users/:userId", api.getUser)

	r.POST("/users/:userId/rsvps/:eventId", api.rsvpAdd)
	r.DELETE("/users/:userId/rsvps/:eventId", api.rsvpRemove)
	r.GET("/users/:userId/rsvps", api.getRSVPs)
}
