package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	txProducer *TransactionService
}

func NewController(txProducer *TransactionService) *Controller {
	return &Controller{
		txProducer: txProducer,
	}
}

func (c *Controller) RegisterRoutes(router *gin.Engine) {
	router.POST("/transfer", c.Transfer)
}

func (c *Controller) Transfer(ctx *gin.Context) {

	// TODO: Send to Kafka producer.

	ctx.JSON(http.StatusOK, gin.H{"message": "Sent"})
}
