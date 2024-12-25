package handlers

import "github.com/gin-gonic/gin"

func bindJSON(c *gin.Context) (map[string]interface{}, error) {
	_json := make(map[string]interface{})
	if bindErr := c.BindJSON(&_json); bindErr != nil {
		return nil, bindErr
	}
	return _json, nil
}

func sendErrorResponse(c *gin.Context, statusCode int, errorCode string, err error) {
	c.JSON(statusCode, gin.H{
		"status":        "error",
		"error_code":    errorCode,
		"error_message": err.Error(),
	})
}
