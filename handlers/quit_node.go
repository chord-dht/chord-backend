package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func QuitNode(c *gin.Context) {
	if LocalNode == nil {
		sendNotExistErrorResponse(c)
		return
	}

	LocalNode.Quit()
	LocalNode = nil // reset the variable

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "quit node succeeded",
	})
}
