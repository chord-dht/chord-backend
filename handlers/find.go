package handlers

import (
	"net/http"

	"github.com/chord-dht/chord-core/tools"
	"github.com/gin-gonic/gin"
)

func Lookup(c *gin.Context) {
	if LocalNode == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"error":   "Node not created",
			"details": "Please create a node first",
		})
		return
	}

	json := make(map[string]interface{})
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"error":   "Failed to bind JSON",
			"details": err.Error(),
		})
		return
	}

	var filename string
	if val, ok := json["filename"].(string); ok {
		filename = val
	}

	fileIdentifier := tools.GenerateIdentifier(filename)

	targetNode, err := LocalNode.GetInfo().FindSuccessorIter(fileIdentifier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"error":   "Failed to find successor",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":          "success",
		"file_identifier": fileIdentifier,
		"node":            targetNode,
	})
}
