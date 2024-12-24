package handlers

import (
	"net/http"

	"github.com/chord-dht/chord-backend/config"

	"github.com/gin-gonic/gin"
)

func InitializeNode(c *gin.Context) {
	if LocalNode == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"error":   "Node not created",
			"details": "Please create a node first",
		})
		return
	}

	err := LocalNode.Initialize(
		config.NodeConfig.Mode,
		config.NodeConfig.JoinAddress,
		config.NodeConfig.JoinPort,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"error":   "Failed to initialize node",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "initialize node successed",
	})
}
