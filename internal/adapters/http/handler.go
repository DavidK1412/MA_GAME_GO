package http

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
	router *gin.Engine
}

func NewHandler() *Handler {
	router := gin.Default()

	h := &Handler{
		router: router,
	}

	h.setupRoutes()

	return h
}

func (h *Handler) setupRoutes() {
	h.router.GET("/ping", h.handlePing)
}

func (h *Handler) handlePing(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (h *Handler) GetRouter() *gin.Engine {
	return h.router
}
