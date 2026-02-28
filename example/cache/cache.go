// Copyright © 2026 Soner Astan astansoner@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
