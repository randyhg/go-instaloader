package caches

import (
	"github.com/goburrow/cache"
	"strconv"
)

func key2Int64(key cache.Key) int64 {
	if s, ok := key.(string); ok {
		v, _ := strconv.ParseInt(s, 10, 64)
		return v
	}
	return key.(int64)
}

func key2String(key cache.Key) string {
	if s, ok := key.(string); ok {
		return s
	}
	return key.(string)
}
