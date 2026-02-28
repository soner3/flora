package cache

import (
	"fmt"

	"github.com/soner3/flora"
)

type Client interface {
	Get(key string) string
}

type RedisClient struct {
	url string
}

func (r *RedisClient) Get(key string) string {
	return "redis: " + r.url
}

type MemoryClient struct{}

func (m *MemoryClient) Get(key string) string {
	return "memory"
}

type CacheConfig struct {
	flora.Configuration
}

func (c *CacheConfig) ProvideRedisClient() (*RedisClient, func(), error) {
	fmt.Println("--> [CacheConfig] ProvideRedisClient (Singleton)")
	return &RedisClient{url: "localhost:6379"}, func() {}, nil
}

// flora:scope=prototype
func (c *CacheConfig) ProvideMemoryClient() (*MemoryClient, func(), error) {
	fmt.Println("--> [CacheConfig] ProvideMemoryClient (Prototype)")
	return &MemoryClient{}, func() {}, nil
}

// flora:primary
func (c *CacheConfig) ProvideDefaultClient(redis *RedisClient) (Client, func(), error) {
	return redis, func() {}, nil
}
