package main

import (
	"chord-backend/config"
	"fmt"
	"net/http"
	"time"

	cfs "github.com/chord-dht/chord-core/cachefilesystem"
	"github.com/chord-dht/chord-core/node"
	st "github.com/chord-dht/chord-core/storage"
	"github.com/gin-gonic/gin"
)

var LocalNode *node.Node = nil

func main() {
	r := gin.Default()

	r.POST("/new", func(c *gin.Context) {
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
	})

	r.GET("/initialize", func(c *gin.Context) {
		if LocalNode == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"error":   "Node not created",
				"details": "Please create a node first",
			})
			return
		}

		LocalNode.Initialize(
			config.NodeConfig.Mode,
			config.NodeConfig.JoinAddress,
			config.NodeConfig.JoinPort,
		)

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "initialize node successed",
		})
	})

	r.GET("/quit", func(c *gin.Context) {
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
	})

	r.GET("/printstate", func(c *gin.Context) {
		if LocalNode == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"error":   "Node not created",
				"details": "Please create a node first",
			})
			return
		}

		nodeState := LocalNode.GetState()
		nodeState.PrintState()

		c.JSON(http.StatusOK, gin.H{
			"status":    "success",
			"nodestate": nodeState,
		})
	})

	r.POST("/storefile", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/getfile", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8080")
}

// NewNodeWithConfig uses the configuration to create a new node.
func NewNodeWithConfig(
	cfg *config.Config,
	storageFactory func(string) (st.Storage, error),
) (*node.Node, error) {
	chordNode, err := node.NewNode(
		cfg.IdentifierLength,
		cfg.SuccessorsLength,
		cfg.IpAddress,
		cfg.Port,
		storageFactory,
		cfg.StorageDir,
		cfg.BackupDir,
		time.Duration(cfg.StabilizeTime),
		time.Duration(cfg.FixFingersTime),
		time.Duration(cfg.CheckPredecessorTime),
		cfg.TLSBool,
		cfg.ServerTLSConfig,
		cfg.ClientTLSConfig,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating node: %w", err)
	}

	return chordNode, nil
}
