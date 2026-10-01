package utils

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var ctx = context.Background()

func GetRedisConnection() (*redis.Client, error) {
	var err = godotenv.Load()
	if err != nil {
		log.Fatal("Unable to Load Database Credentials from env: ", err.Error())
		return nil, err
	}
	rdb = redis.NewClient(&redis.Options{
		Addr:         os.Getenv("REDIS_URL"),
		Username:     os.Getenv("REDIS_USERNAME"),
		Password:     os.Getenv("REDIS_PASSWORD"),
		DB:           0,
		PoolSize:     50,
		MinIdleConns: 10,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Unable to connect to Redis: ", err.Error())
		return nil, err
	}
	log.Info("Connected to Redis. Server response: %s\n", pong)
	return rdb, nil
}

func SetListData[T any](key string, result []any, ttl time.Duration) bool {
	var list = FormatListResult[T](result)
	json_data, err := json.Marshal(list)
	if err != nil {
		log.Error(err.Error())
		return false
	}
	if err := rdb.Set(ctx, key, json_data, ttl).Err(); err != nil {
		return false
	}
	return true
}

func SetSingleData[T any](key string, result []any, ttl time.Duration) bool {
	var data = FormatSingleResult[T](result)
	json_data, err := json.Marshal(data)
	if err != nil {
		log.Error(err.Error())
		return false
	}
	if err := rdb.Set(ctx, key, json_data, ttl).Err(); err != nil {
		return false
	}
	return true
}

func SetStringData(key, value string, ttl time.Duration) bool {
	if err := rdb.Set(ctx, key, value, ttl).Err(); err != nil {
		return false
	}
	return true
}

func GetData(key string) (string, error) {
	data, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Error("Key does not exist in Redis. Cache Miss ", err.Error())
		return "", err
	} else if err != nil {
		log.Error("Unable to fetch data from Redis: ", err.Error())
		return "", err
	} else {
		return data, nil
	}
}

func DeleteData(key string) bool {
	delete_count, err := rdb.Unlink(ctx, key).Result()
	if err != nil {
		log.Error("Unable to delete key: ", err.Error())
		return false
	}
	if delete_count == 0 {
		log.Info("No matching key found for the key: ", key)
	} else {
		log.Info("Item Deleted from Redis")
	}
	return true
}

func DeleteKeyPattern(pattern string) bool {
	var cursor uint64
	var delete_keys []string

	for {
		var keys []string
		var err error
		keys, cursor, err = rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			log.Error("Error in scanning keys ", err.Error())
			return false
		}
		delete_keys = append(delete_keys, keys...)
		if cursor == 0 {
			break
		}
	}

	if len(delete_keys) > 0 {
		delete_count, err := rdb.Unlink(ctx, delete_keys...).Result()
		if err != nil {
			log.Error("Error in unlinking keys ", err.Error())
			return false
		}
		log.Info("Deleted Keys: ", delete_count)
	} else {
		log.Info("No matching key found")
	}

	return true
}
