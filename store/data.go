package store

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var c *cache.Cache

func init() {
	c = cache.New(5*time.Minute, 10*time.Minute)
}

func SetData(key string, val interface{}) {
	c.Set(key, val, cache.DefaultExpiration)
}

func GetData(key string) (interface{}, bool) {
	return c.Get(key)
}
