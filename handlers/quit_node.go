package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func QuitNode(c *gin.Context) {
	if LocalNode == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"error":   "Node not created",
			"details": "Please create a node first",
		})
		return
	}

	LocalNode.Quit()
	LocalNode = nil // reset the variable

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "quit node successed",
	})
}
