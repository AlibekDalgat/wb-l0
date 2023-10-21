package handlers

import (
	"github.com/gin-gonic/gin"
	"wb-l0/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	api := router.Group("/api")
	{
		api.GET("/:id", h.getOrder)
	}

	return router
}
