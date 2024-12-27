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

	if LocalNode == nil {
		sendNotExistErrorResponse(c)
		return
	}

	// Generate the file identifier
	// Keep in mind that you shouldn't invoke it when the node is not created
	fileIdentifier := tools.GenerateIdentifier(fileHeader.Filename)

	file, err := fileHeader.Open()
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError,
			"OPEN_FILE_ERROR", err,
			gin.H{"file_identifier": fileIdentifier},
		)
		return
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError,
			"READ_FILE_ERROR", err,
			gin.H{"file_identifier": fileIdentifier},
		)
		return
	}

	targetNode, err := LocalNode.GetInfo().FindSuccessorIter(fileIdentifier)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError,
			"FIND_ERROR", err,
			gin.H{"file_identifier": fileIdentifier},
		)
		return
	}

	if config.NodeConfig.AESBool {
		fileContent, err = aes.EncryptAES(fileContent, config.NodeConfig.AESKey)
		if err != nil {
			sendErrorResponse(c, http.StatusInternalServerError,
				"ENCRYPT_ERROR", err,
				gin.H{
					"file_identifier": fileIdentifier,
					"target_node":     targetNode,
				},
			)
			return
		}
	}

	reply, err := targetNode.StoreFile(fileHeader.Filename, fileContent)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError,
			"ACCESS_ERROR", err,
			gin.H{
				"file_identifier": fileIdentifier,
				"target_node":     targetNode,
			},
		)
		return
	}
	if !reply.Success {
		sendErrorResponse(c, http.StatusInternalServerError,
			"STORE_DENIED_ERROR", errors.New("target node reply: it can't store the file"),
			gin.H{
				"file_identifier": fileIdentifier,
				"target_node":     targetNode,
			},
		)
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
