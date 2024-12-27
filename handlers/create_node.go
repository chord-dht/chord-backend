package handlers

import (
	"net/http"

	"github.com/chord-dht/chord-backend/config"

	cfs "github.com/chord-dht/chord-core/cachefilesystem"
	"github.com/gin-gonic/gin"
)

func CreateNode(c *gin.Context) {
	if LocalNode != nil {
		sendExistErrorResponse(c)
		return
	}

	cfgJson, bindErr := bindJSON(c)
	if bindErr != nil {
		sendBindJSONErrorResponse(c, bindErr)
		return
	}

	cfg, parseErr := config.JsonToConfig(cfgJson)
	if parseErr != nil {
		sendParseJSONErrorResponse(c, parseErr)
		return
	}

	if validErr := config.ValidateAndSetConfig(cfg); validErr != nil {
		sendErrorResponse(c, http.StatusBadRequest, "VALIDATE_CONFIG_ERROR", validErr)
		return
	}

	config.NodeConfig = cfg

	newNode, newErr := NewNodeWithConfig(config.NodeConfig, cfs.CacheStorageFactory)
	if newErr != nil {
		sendErrorResponse(c, http.StatusInternalServerError, "CREATE_NODE_ERROR", newErr)
		return
	}

	LocalNode = newNode

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "new node succeeded",
	})
}
