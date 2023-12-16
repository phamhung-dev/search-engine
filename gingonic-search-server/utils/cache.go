package utils

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
)

func GenerateCacheKey(table string, options ...interface{}) string {
	content := []byte{}

	for _, option := range options {
		data, err := json.Marshal(option)
		if err != nil {
			continue
		}
		content = append(content, data...)
	}

	cacheKey := fmt.Sprintf("%s:%x", table, md5.Sum(content))

	return cacheKey
}
