package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func sendErrorResponse(c *gin.Context, statusCode int, errorCode string, err error, data ...interface{}) {
	response := gin.H{
		"status":        "error",
		"error_code":    errorCode,
		"error_message": err.Error(),
	}

	if len(data) > 0 {
		response["data"] = data[0]
	}

	c.JSON(statusCode, response)
}

func sendBindJSONErrorResponse(c *gin.Context, bindErr error) {
	sendErrorResponse(c, http.StatusBadRequest, "BIND_JSON_ERROR", bindErr)
}

func sendParseJSONErrorResponse(c *gin.Context, parseErr error) {
	sendErrorResponse(c, http.StatusBadRequest, "PARSE_JSON_ERROR", parseErr)
}

var ErrNodeNotExist = errors.New("node not created: Please create a node first")
var ErrNodeExist = errors.New("node already exists: Please quit the existing node first")

func sendNotExistErrorResponse(c *gin.Context, data ...interface{}) {
	sendErrorResponse(c, http.StatusBadRequest, "NODE_NOT_EXISTS_ERROR", ErrNodeNotExist, data)
}

func sendExistErrorResponse(c *gin.Context, data ...interface{}) {
	sendErrorResponse(c, http.StatusBadRequest, "NODE_EXISTS_ERROR", ErrNodeExist, data)
}
