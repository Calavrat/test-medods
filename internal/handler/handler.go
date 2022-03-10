package handler

import (
	"github.com/Calavrat/TestMedods/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.GET("/sing-in/:id", h.SingIn)
		auth.POST("/refresh", h.Refresh)
	}

	return router
}
