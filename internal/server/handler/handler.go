package handler

import (
	"api-validator/internal/domains"
	"api-validator/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	validatorService *service.ValidatorService
}

func NewHandler(vs *service.ValidatorService) *Handler {
	return &Handler{
		validatorService: vs,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	api := router.Group("/validator")

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	api.POST("/validate", h.Validate)

	return router
}

func (h *Handler) Validate(c *gin.Context) {
	var req domains.ValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.validatorService.ValidateProduct(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
