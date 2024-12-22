package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func StoreFile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func GetFile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
