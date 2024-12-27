package handlers

import "github.com/gin-gonic/gin"

func bindJSON(c *gin.Context) (map[string]interface{}, error) {
	_json := make(map[string]interface{})
	if bindErr := c.BindJSON(&_json); bindErr != nil {
		return nil, bindErr
	}
	return _json, nil
}
