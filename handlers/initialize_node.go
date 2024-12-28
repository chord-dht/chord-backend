package handlers

import (
	"errors"
	"net/http"

	"github.com/chord-dht/chord-backend/config"

	"github.com/gin-gonic/gin"
)

func InitializeNode(c *gin.Context) {
	if LocalNode == nil {
		sendNotExistErrorResponse(c)
		return
	}

	if err := LocalNode.Initialize(
		config.NodeConfig.Mode,
		config.NodeConfig.JoinAddress,
		config.NodeConfig.JoinPort,
	); err != nil {
		sendErrorResponse(c, http.StatusInternalServerError,
			"INITIALIZE_ERROR", errors.New("failed to initialize node: "+err.Error()),
		)
		LocalNode = nil
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "initialize node succeeded",
	})
}
