package handlers

import (
	"errors"
	"net/http"

	"github.com/chord-dht/chord-backend/config"

	cfs "github.com/chord-dht/chord-core/cachefilesystem"
	"github.com/gin-gonic/gin"
)

func CreateNode(c *gin.Context) {
	if LocalNode != nil {
		sendErrorResponse(c, http.StatusBadRequest, "NODE_EXISTS_ERROR", errors.New("node already exists: Please quit the existing node first"))
		return
	}

	cfgJson, bindErr := bindJSON(c)
	if bindErr != nil {
		sendErrorResponse(c, http.StatusBadRequest, "BIND_JSON_ERROR", bindErr)
		return
	}

	cfg, parseErr := config.JsonToConfig(cfgJson)
	if parseErr != nil {
		sendErrorResponse(c, http.StatusBadRequest, "PARSE_JSON_ERROR", parseErr)
		return
	}

	if validErr := config.ValidateAndSetConfig(cfg); validErr != nil {
		sendErrorResponse(c, http.StatusBadRequest, "VALIDATE_CONFIG_ERROR", validErr)
		return
	}

	config.NodeConfig = cfg

	var newErr error
	LocalNode, newErr = NewNodeWithConfig(config.NodeConfig, cfs.CacheStorageFactory)
	if newErr != nil {
		sendErrorResponse(c, http.StatusInternalServerError, "CREATE_NODE_ERROR", newErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "new node succeeded",
	})
}
