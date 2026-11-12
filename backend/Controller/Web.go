package Controller

import (
	"campus_connect/Model"

	"github.com/gin-gonic/gin"
)

// starts the router
func Serve(repository *Model.Repository) error {
	router := newRouter(repository)
	return router.Run("0.0.0.0:8080")
}

func newRouter(repository *Model.Repository) *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	api := NewAPI(repository)
	defineRoutes(router, api)

	return router
}

// define get/post routes here
func defineRoutes(r *gin.Engine, api *API) {
	r.GET("/events", api.getEvents)
	r.POST("/events", api.requireAuth, api.createEvent)
	r.PATCH("/events/:id", api.requireAuth, api.updateEvent)
	r.GET("/events/:id", api.getEvent)

	r.POST("/users/register", api.registerUser)
	r.POST("/users/login", api.loginUser)
	r.GET("/users/me", api.requireAuth, api.getCurrentUser)
	r.GET("/users/:userId", api.getUser)

	r.POST("/users/:userId/rsvps/:eventId", api.requireAuth, api.rsvpAdd)
	r.DELETE("/users/:userId/rsvps/:eventId", api.requireAuth, api.rsvpRemove)
	r.GET("/users/:userId/rsvps", api.requireAuth, api.getRSVPs)
}
