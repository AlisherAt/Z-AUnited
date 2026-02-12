package cache

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func InitRedis(addr, password string, db int) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

func SetRedis(key string, value interface{}, ttl time.Duration) error {
	if redisClient == nil {
		return nil
	}
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return redisClient.Set(context.Background(), key, b, ttl).Err()
}

func GetRedis(key string, dest interface{}) (bool, error) {
	if redisClient == nil {
		return false, nil
	}
	val, err := redisClient.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return json.Unmarshal([]byte(val), dest) == nil, nil
}

type Item struct {
	Value      interface{}
	Expiration int64
}

type Cache struct {
	mu    sync.RWMutex
	store map[string]Item
}

func New() *Cache {
	return &Cache{store: make(map[string]Item)}
}

func (c *Cache) Set(key string, v interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = Item{Value: v, Expiration: time.Now().Add(ttl).Unix()}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	it, ok := c.store[key]
	if !ok || time.Now().Unix() > it.Expiration {
		return nil, false
	}
	return it.Value, true
}
