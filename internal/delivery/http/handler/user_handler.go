package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/Fox1N69/iq-testtask/internal/domain/service"
	"github.com/Fox1N69/iq-testtask/pkg/logger"
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	UserByID(c *gin.Context)
}

type userHandler struct {
	service service.UserService
	log     logger.Logger
}

func NewUserHandler(service service.UserService) UserHandler {
	return &userHandler{
		service: service,
		log:     logger.GetLogger(),
	}
}

func (h *userHandler) UserByID(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.service.UserByID(c.Request.Context(), userID)
	if err != nil {
		h.log.Errorf("Failed to get user: %v", err)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, user)
}
