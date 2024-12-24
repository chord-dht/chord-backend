package handlers

import (
	"chord-backend/aes"
	"chord-backend/config"
	"io"
	"net/http"

	"github.com/chord-dht/chord-core/tools"
	"github.com/gin-gonic/gin"
)

func StoreFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Failed to store file",
			"details": err.Error(),
		})
		return
	}

	_file, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to open file",
			"details": err.Error(),
		})
		return
	}
	defer _file.Close()

	fileContent, err := io.ReadAll(_file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to read file content",
			"details": err.Error(),
		})
		return
	}

	fileIdentifier := tools.GenerateIdentifier(file.Filename)

	if LocalNode == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":          "error",
			"message":         "Node not created",
			"details":         "Please create a node first",
			"file_identifier": fileIdentifier,
		})
		return
	}

	targetNode, err := LocalNode.GetInfo().FindSuccessorIter(fileIdentifier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":          "error",
			"message":         "Failed to find successor",
			"details":         err.Error(),
			"file_identifier": fileIdentifier,
		})
		return
	}

	if config.NodeConfig.AESBool {
		fileContent, err = aes.EncryptAES(fileContent, config.NodeConfig.AESKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":          "error",
				"message":         "Failed to encrypt the file content",
				"details":         err.Error(),
				"file_identifier": fileIdentifier,
				"target_node":     targetNode,
			})
			return
		}
	}

	reply, err := targetNode.StoreFile(file.Filename, fileContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":          "error",
			"message":         "Failed to get the reply from target node",
			"details":         err.Error(),
			"file_identifier": fileIdentifier,
			"target_node":     targetNode,
		})
		return
	}
	if !reply.Success {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":          "error",
			"message":         "Target node reply: it can't store the file",
			"file_identifier": fileIdentifier,
			"target_node":     targetNode,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":          "success",
		"file_identifier": fileIdentifier,
		"target_node":     targetNode,
	})
}
