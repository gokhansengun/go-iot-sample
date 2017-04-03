package iot

import (
	"github.com/go-martini/martini"
	"github.com/go-redis/redis"
)

// RedisSession is the struct to keep Sarama
type RedisSession struct {
	*redis.Client
	redisAddr string
}

// NewRedisSession connects to Redis
func NewRedisSession(addr string) *RedisSession {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()

	if err != nil {
		panic(err)
	}

	return &RedisSession{client, addr}
}

// SetKey sets a value with the specified key
func (redisSession *RedisSession) SetKey(key string, value interface{}) error {
	err := redisSession.Set(key, value, 0).Err()

	if err != nil {
		return err
	}

	return nil
}

// GetKey sets a value with the specified key
func (redisSession *RedisSession) GetKey(key string) (string, error) {
	val, err := redisSession.Get(key).Result()

	if err != nil {
		return "", err
	}

	return val, nil
}

// KeyExists checks existence of a key
func (redisSession *RedisSession) KeyExists(key string) bool {
	val, err := redisSession.Exists(key).Result()

	if err != nil {
		return false
	}

	return val > 0
}

// NewRedisHandler adds Redis to the Martini pipeline
func (redisSession *RedisSession) NewRedisHandler() martini.Handler {
	return func(context martini.Context) {
		context.Map(redisSession)
		context.Next()
	}
}
