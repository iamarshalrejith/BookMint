package redis

import (
	"context"
	"log"

	goredis "github.com/redis/go-redis/v9"
)

func NewClient(addr string) *goredis.Client {
	rdb := goredis.NewClient(&goredis.Options{Addr: addr}) // "I want to communicate with the Redis server running at "I want to communicate with the Redis server running at localhost
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Redis Ping: %v", err)
	}
	log.Printf("Connected to redis at %s",addr)

	return rdb
}
