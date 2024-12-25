package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/chord-dht/chord-backend/aes"
	"github.com/chord-dht/chord-backend/config"
	"github.com/chord-dht/chord-core/tools"
	"github.com/gin-gonic/gin"
)

func StoreFile(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "FORM_FILE_ERROR", err)
		return
	}

	fileIdentifier := tools.GenerateIdentifier(fileHeader.Filename)

	if LocalNode == nil {
		sendErrorResponse(c, http.StatusBadRequest, "NODE_NOT_EXISTS_ERROR", errors.New("node not created: Please create a node first"))
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, "OPEN_FILE_ERROR", err)
		return
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, "READ_FILE_ERROR", err)
		return
	}

	targetNode, err := LocalNode.GetInfo().FindSuccessorIter(fileIdentifier)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, "FIND_ERROR", err)
		return
	}

	if config.NodeConfig.AESBool {
		fileContent, err = aes.EncryptAES(fileContent, config.NodeConfig.AESKey)
		if err != nil {
			sendErrorResponse(c, http.StatusInternalServerError, "ENCRYPT_ERROR", err)
			return
		}
	}

	reply, err := targetNode.StoreFile(fileHeader.Filename, fileContent)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, "ACCESS_ERROR", err)
		return
	}
	if !reply.Success {
		sendErrorResponse(c, http.StatusInternalServerError, "STORE_DENIED_ERROR", errors.New("target node reply: it can't store the file"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"file_identifier": fileIdentifier,
			"target_node":     targetNode,
		},
	})
}
