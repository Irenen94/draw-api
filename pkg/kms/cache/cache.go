package cache

import (
	"encoding/json"
	"github.com/coocood/freecache"
)

const (
	cacheSize = 1024 * 1024
)

var (
	kmsCache = freecache.NewCache(cacheSize)
)

func Set(key string, value interface{}, expireSeconds int) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return kmsCache.Set([]byte(key), v, expireSeconds)
}

func Get(key string, data interface{}) error {
	v, err := kmsCache.Get([]byte(key))
	if err != nil {
		return err
	}
	return json.Unmarshal(v, data)
}
