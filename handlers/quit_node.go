package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func QuitNode(c *gin.Context) {
	if LocalNode == nil {
		sendErrorResponse(c, http.StatusBadRequest, "NODE_NOT_EXISTS_ERROR", errors.New("node not created: Please create a node first"))
		return
	}

	LocalNode.Quit()
	LocalNode = nil // reset the variable

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "quit node succeeded",
	})
}
