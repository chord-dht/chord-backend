package json

import "fmt"

func GetIntFromJson(json map[string]interface{}, key string) (int, error) {
	val, ok := json[key]
	if !ok {
		return 0, fmt.Errorf("%s must be specified", key)
	}
	if v, ok := val.(float64); ok {
		return int(v), nil
	}
	return 0, fmt.Errorf("%s must be a float64, got %T", key, val)
}

func GetStringFromJson(json map[string]interface{}, key string) (string, error) {
	val, ok := json[key]
	if !ok {
		return "", fmt.Errorf("%s must be specified", key)
	}
	if v, ok := val.(string); ok {
		return v, nil
	}
	return "", fmt.Errorf("%s must be a string, got %T", key, val)
}

func GetBoolFromJson(json map[string]interface{}, key string) (bool, error) {
	val, ok := json[key]
	if !ok {
		return false, fmt.Errorf("%s must be specified", key)
	}
	if v, ok := val.(bool); ok {
		return v, nil
	}
	return false, fmt.Errorf("%s must be a bool, got %T", key, val)
}
