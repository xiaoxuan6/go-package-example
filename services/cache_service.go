package services

import (
	cache2 "github.com/patrickmn/go-cache"
	"package-example/common"
)

type cacheService struct {
}

var CacheService = newCacheService()

func newCacheService() *cacheService {
	return new(cacheService)
}

func (c cacheService) Set(key, val string) {
	common.Cache.Set(key, val, cache2.DefaultExpiration)
}

func (c cacheService) Get(key string) string {
	val, found := common.Cache.Get(key)
	if !found {
		return ""
	}

	return val.(string)
}

func (c cacheService) Del(key string)  {
	common.Cache.Delete(key)
}

// Flush 清楚所有的缓存
func (c cacheService) Flush()  {
	common.Cache.Flush()
}
