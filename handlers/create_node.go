package handlers

import (
	"net/http"

	"github.com/chord-dht/chord-backend/config"

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
	if bindErr := c.BindJSON(&json); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"error":   "Failed to bind JSON",
			"details": bindErr.Error(),
		})
		return
	}

	cfg, parseErr := config.JsonToConfig(json)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"error":   "Failed to parse JSON",
			"details": parseErr.Error(),
		})
		return
	}

	if validErr := config.ValidateAndSetConfig(cfg); validErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"error":   "Failed to validate config",
			"details": validErr.Error(),
		})
		return
	}

	config.NodeConfig = cfg

	var newErr error = nil
	LocalNode, newErr = NewNodeWithConfig(config.NodeConfig, cfs.CacheStorageFactory)
	if newErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"error":   "Failed to create node",
			"details": newErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "new node succeeded",
	})
}
