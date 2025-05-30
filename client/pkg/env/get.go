package env

import (
	"os"
	"strconv"
)

// GetEnv 获取环境变量
func GetEnv[T any](key string, defaultValue T) T {
	var result T
	strValue, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	switch any(result).(type) {
	case string:
		return any(strValue).(T)
	case int:
		if v, err := strconv.Atoi(strValue); err == nil {
			return any(v).(T)
		}
	case bool:
		if v, err := strconv.ParseBool(strValue); err == nil {
			return any(v).(T)
		}
	case float64:
		if v, err := strconv.ParseFloat(strValue, 64); err == nil {
			return any(v).(T)
		}
	}

	return defaultValue
}
