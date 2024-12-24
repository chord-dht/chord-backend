package handlers

import (
	"chord-backend/aes"
	"chord-backend/config"
	"net/http"
	"os"
	"path/filepath"

	"github.com/chord-dht/chord-core/tools"
	"github.com/gin-gonic/gin"
)

var tempDir string = "temp"

func init() {
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		os.Mkdir(tempDir, 0755)
	}
}

func GetFile(c *gin.Context) {
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
			"error":           "Failed to find successor",
			"details":         err.Error(),
			"file_identifier": fileIdentifier,
		})
		return
	}

	reply, err := targetNode.GetFile(filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":          "error",
			"error":           "Failed to get the reply from node",
			"details":         err.Error(),
			"file_identifier": fileIdentifier,
			"target_node":     targetNode,
		})
		return
	}
	if !reply.Success {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":          "error",
			"error":           "Target node reply: it can't get the file",
			"file_identifier": fileIdentifier,
			"target_node":     targetNode,
		})
		return
	}

	fileContent := reply.FileContent

	if config.NodeConfig.AESBool {
		fileContent, err = aes.DecryptAES(fileContent, config.NodeConfig.AESKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":          "error",
				"error":           "Failed to decrypt the file content",
				"details":         err.Error(),
				"file_identifier": fileIdentifier,
				"target_node":     targetNode,
			})
			return
		}
	}

	tempFilePath := filepath.Join(tempDir, filename)
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":          "error",
			"message":         "Failed to create temp file",
			"details":         err.Error(),
			"file_identifier": fileIdentifier,
			"target_node":     targetNode,
		})
		return
	}
	defer tempFile.Close()

	if _, err := tempFile.Write(fileContent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":          "error",
			"message":         "Failed to write to temp file",
			"details":         err.Error(),
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

func DownloadFile(c *gin.Context) {
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

	if LocalNode == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Node not created",
			"details": "Please create a node first",
		})
		return
	}

	tempFilePath := filepath.Join(tempDir, filename)

	// Check if the file exists
	if _, err := os.Stat(tempFilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "File not found",
			"details": err.Error(),
		})
		return
	}
	defer os.Remove(tempFilePath)

	// Set the filename in the header
	c.Header("Content-Disposition", "attachment; filename="+filename)
	// Send the file
	c.File(tempFilePath)
}
