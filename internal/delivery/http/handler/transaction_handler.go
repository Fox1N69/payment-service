package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Fox1N69/iq-testtask/internal/domain/service"
	"github.com/gin-gonic/gin"
)

type TransactionHandler interface {
	Replenish(c *gin.Context)
	Transfer(c *gin.Context)
	LastTransactions(c *gin.Context)
}

type transactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(service service.TransactionService) TransactionHandler {
	return &transactionHandler{
		service: service,
	}
}

func (h *transactionHandler) Replenish(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	amountStr := c.Param("amount")
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	err = h.service.Replenish(context.Background(), userID, amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to replenish"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "replenished successfully"})
}

func (h *transactionHandler) Transfer(c *gin.Context) {
	fromUserIDStr := c.Query("from_user_id")
	fromUserID, err := strconv.ParseInt(fromUserIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid from user ID"})
		return
	}

	toUserIDStr := c.Query("to_user_id")
	toUserID, err := strconv.ParseInt(toUserIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid to user ID"})
		return
	}

	amountStr := c.Query("amount")
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	err = h.service.Transfer(context.Background(), fromUserID, toUserID, amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to transfer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "transferred successfully"})
}

func (h *transactionHandler) LastTransactions(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	limitStr := c.Param("limit")
	limit := int64(10)
	if limitStr != "" {
		limit, err = strconv.ParseInt(limitStr, 10, 8)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
			return
		}
	}

	transactions, err := h.service.LastTransactions(context.Background(), userID, int8(limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch last transactions"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
