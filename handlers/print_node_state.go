package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PrintNodeState(c *gin.Context) {
	if LocalNode == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"error":   "Node not created",
			"details": "Please create a node first",
		})
		return
	}

	nodeState := LocalNode.GetState()
	nodeState.PrintState()

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"nodestate": nodeState,
	})
}
