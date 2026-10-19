package Controller

import (
	"campus_connect/Model"
	"sync"

	"github.com/gin-gonic/gin"
)

type API struct {
	mu     sync.RWMutex
	campus *Model.Campus
}

func NewAPI(campus *Model.Campus) *API {
	return &API{
		campus: campus,
	}
}

func (api *API) readLock() gin.HandlerFunc {
	return func(c *gin.Context) {
		api.mu.RLock()
		defer api.mu.RUnlock()
		c.Next()
	}
}

func (api *API) writeLock() gin.HandlerFunc {
	return func(c *gin.Context) {
		api.mu.Lock()
		defer api.mu.Unlock()
		c.Next()
	}
}
