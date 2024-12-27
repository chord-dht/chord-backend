package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetNodeState(c *gin.Context) {
	if LocalNode == nil {
		sendNotExistErrorResponse(c)
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
