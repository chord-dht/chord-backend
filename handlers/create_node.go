package handlers

import (
	"chord-backend/config"
	"net/http"

	cfs "github.com/chord-dht/chord-core/cachefilesystem"
	"github.com/gin-gonic/gin"
)

func CreateNode(c *gin.Context) {
	if LocalNode != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"error":   "Node already exists",
			"details": "Please quit the existing node first",
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

	cfg := config.JsontToConfig(json)

	if err := config.ValidateAndSetConfig(cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"error":   "Failed to validate config",
			"details": err.Error(),
		})
		return
	}

	config.NodeConfig = cfg

	var err error = nil
	LocalNode, err = NewNodeWithConfig(config.NodeConfig, cfs.CacheStorageFactory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"error":   "Failed to create node",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "new node successed",
	})
}
