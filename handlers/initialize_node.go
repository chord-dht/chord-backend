package handlers

import (
	"chord-backend/config"
	"net/http"

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

	LocalNode.Initialize(
		config.NodeConfig.Mode,
		config.NodeConfig.JoinAddress,
		config.NodeConfig.JoinPort,
	)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "initialize node successed",
	})
}
