package handlers

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"github.com/chord-dht/chord-backend/config"
	"github.com/chord-dht/chord-backend/json"

	"github.com/chord-dht/chord-backend/aes"

	"github.com/chord-dht/chord-core/tools"
	"github.com/gin-gonic/gin"
)

var tempDir string = "temp"

func init() {
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		os.Mkdir(tempDir, 0755)
	}
}

func createTempFile(tempDir, filename string, fileContent []byte) error {
	tempFilePath := filepath.Join(tempDir, filename)
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	if _, err := tempFile.Write(fileContent); err != nil {
		return err
	}

	return nil
}

func GetFile(c *gin.Context) {
	filenameJson, bindErr := bindJSON(c)
	if bindErr != nil {
		sendBindJSONErrorResponse(c, bindErr)
		return
	}

	filename, parseErr := json.GetStringFromJson(filenameJson, "filename")
	if parseErr != nil {
		sendParseJSONErrorResponse(c, parseErr)
		return
	}

	if LocalNode == nil {
		sendNotExistErrorResponse(c)
		return
	}

	// Generate the file identifier
	// Keep in mind that you shouldn't invoke it when the node is not created
	fileIdentifier := tools.GenerateIdentifier(filename)

	targetNode, err := LocalNode.GetInfo().FindSuccessorIter(fileIdentifier)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError,
			"FIND_ERROR", err,
			gin.H{"file_identifier": fileIdentifier},
		)
		return
	}

	reply, err := targetNode.GetFile(filename)
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
			"NON_FILE_ERROR", errors.New("target node reply: it doesn't have the file"),
			gin.H{
				"file_identifier": fileIdentifier,
				"target_node":     targetNode,
			},
		)
		return
	}

	fileContent := reply.FileContent

	if config.NodeConfig.AESBool {
		fileContent, err = aes.DecryptAES(fileContent, config.NodeConfig.AESKey)
		if err != nil {
			sendErrorResponse(c, http.StatusInternalServerError,
				"DECRYPT_ERROR", err,
				gin.H{
					"file_identifier": fileIdentifier,
					"target_node":     targetNode,
				},
			)
			return
		}
	}

	if err := createTempFile(tempDir, filename, fileContent); err != nil {
		sendErrorResponse(c, http.StatusInternalServerError,
			"TEMP_ERROR", err,
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

func DownloadFile(c *gin.Context) {
	filenameJson, bindErr := bindJSON(c)
	if bindErr != nil {
		sendBindJSONErrorResponse(c, bindErr)
		return
	}

	filename, parseErr := json.GetStringFromJson(filenameJson, "filename")
	if parseErr != nil {
		sendParseJSONErrorResponse(c, parseErr)
		return
	}

	tempFilePath := filepath.Join(tempDir, filename)

	// Check if the file exists
	if _, err := os.Stat(tempFilePath); os.IsNotExist(err) {
		sendErrorResponse(c, http.StatusNotFound, "FILE_NOT_FIND_ERROR", errors.New("file not found"))
		return
	}
	defer os.Remove(tempFilePath)

	// Set the filename in the header
	c.Header("Content-Disposition", "attachment; filename="+filename)
	// Send the file
	c.File(tempFilePath)
}
