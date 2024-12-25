package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetNodeState(c *gin.Context) {
	if LocalNode == nil {
		sendErrorResponse(c, http.StatusBadRequest, "NODE_NOT_EXISTS_ERROR", errors.New("node not created: Please create a node first"))
		return
	}

	nodeState := LocalNode.GetState()

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"node_state": nodeState,
		},
	})
}
