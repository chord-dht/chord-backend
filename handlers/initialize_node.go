package handlers

import (
	"errors"
	"net/http"

	"github.com/chord-dht/chord-backend/config"

	"github.com/gin-gonic/gin"
)

func InitializeNode(c *gin.Context) {
	if LocalNode == nil {
		sendErrorResponse(c, http.StatusBadRequest, "NODE_NOT_EXISTS_ERROR", errors.New("node not created: Please create a node first"))
		return
	}

	if err := LocalNode.Initialize(
		config.NodeConfig.Mode,
		config.NodeConfig.JoinAddress,
		config.NodeConfig.JoinPort,
	); err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, "INITIALIZE_ERROR", errors.New("failed to initialize node: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "initialize node succeeded",
	})
}
