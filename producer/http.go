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
	router.POST("/generate-transfers", c.GenerateTransfers)
}

func (c *Controller) GenerateTransfers(ctx *gin.Context) {
	type Body struct {
		Control bool `json:"control"`
	}
	body := Body{}
	if err := ctx.ShouldBindJSON(body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.txProducer.GenerateTransfers(body.Control)
	ctx.JSON(http.StatusOK, gin.H{"message": "Sent"})
}
