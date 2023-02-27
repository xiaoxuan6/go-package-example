package common

import (
	cache2 "github.com/patrickmn/go-cache"
	"time"
)

var Cache *cache2.Cache

func Init() {
	Cache = cache2.New(10*time.Minute, 10*time.Minute) // 第一个参数：第二个参数：
}
