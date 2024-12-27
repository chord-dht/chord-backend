package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NodeStatus(c *gin.Context) {
	if LocalNode == nil {
		c.JSON(http.StatusOK, gin.H{
			"exists": true,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"exists": false,
		})
	}
}
