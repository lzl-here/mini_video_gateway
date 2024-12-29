package repo

import (
	"github.com/go-redis/redis"
)

type CacheRepo struct {
	cache *redis.Client
}

func NewCacheRepo(c *redis.Client) *CacheRepo {
	return &CacheRepo{cache: c}
}
