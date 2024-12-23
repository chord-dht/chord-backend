package handlers

import (
	"chord-backend/aes"
	"chord-backend/config"
	"io"
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

func StoreFile(c *gin.Context) {
	if LocalNode == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Node not created",
			"details": "Please create a node first",
		})
		return
	}

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
	targetNode, err := LocalNode.GetInfo().FindSuccessorIter(fileIdentifier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to find successor",
			"details": err.Error(),
		})
		return
	}

	if config.NodeConfig.AESBool {
		fileContent, err = aes.EncryptAES(fileContent, config.NodeConfig.AESKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to encrypt the file content",
				"details": err.Error(),
			})
			return
		}
	}

	reply, err := targetNode.StoreFile(file.Filename, fileContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to get the reply from node",
			"details": err.Error(),
		})
		return
	}
	if !reply.Success {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Node reply: it can't store the file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":          "success",
		"file_identifier": fileIdentifier,
	})
}

func GetFile(c *gin.Context) {
	if LocalNode == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Node not created",
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

	reply, err := targetNode.GetFile(filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"error":   "Failed to get the reply from node",
			"details": err.Error(),
		})
		return
	}
	if !reply.Success {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Node reply: it can't get the file",
		})
		return
	}

	fileContent := reply.FileContent

	if config.NodeConfig.AESBool {
		fileContent, err = aes.DecryptAES(fileContent, config.NodeConfig.AESKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"error":   "Failed to decrypt the file content",
				"details": err.Error(),
			})
			return
		}
	}

	tempFilePath := filepath.Join(tempDir, filename)
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to create temp file",
			"details": err.Error(),
		})
		return
	}
	defer tempFile.Close()

	if _, err := tempFile.Write(fileContent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to write to temp file",
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

func DownloadFile(c *gin.Context) {
	if LocalNode == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Node not created",
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

	tempFilePath := filepath.Join(tempDir, filename)
	defer os.Remove(tempFilePath)

	// Set the filename in the header
	c.Header("Content-Disposition", "attachment; filename="+filename)
	// Send the file
	c.File(tempFilePath)
}
