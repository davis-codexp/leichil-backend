package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var ctx = context.Background()

const (
	deviceIdsByPhone = "devices:phone:%s"
	phoneByDeviceId  = "phone:device:%s"
	allDeviceIds     = "all:devices"
)

/*
local existing = redis.call('GET', KEYS[2])
if existing and existing ~= ARGV[1] then
	return redis.error_reply("Device already linked to another phone: " .. existing)
end
*/

var add_device_script = redis.NewScript(`
redis.call('SADD', KEYS[1], ARGV[2])
redis.call('SET', KEYS[2], ARGV[1])
redis.call('SADD', KEYS[3], ARGV[2])
return 1
`)

var remove_device_script = redis.NewScript(`
redis.call('SREM', KEYS[1], ARGV[1])
redis.call('DEL', KEYS[2])
redis.call('SREM', KEYS[3], ARGV[1])
return 1
`)

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

func SetHashData(hash, key, value string) bool {
	if err := rdb.HSet(ctx, hash, key, value).Err(); err != nil {
		return false
	}
	return true
}

func GetAllHashData(hash string) (map[string]string, error) {
	var list = make(map[string]string)
	var iter *redis.ScanIterator = rdb.HScan(ctx, hash, 0, "*", 1000).Iterator()
	for iter.Next(ctx) {
		var field = iter.Val()
		if !iter.Next(ctx) {
			break
		}
		var value = iter.Val()
		list[field] = value
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func GetSingleHashData(hash, key string) (string, error) {
	data, err := rdb.HGet(ctx, hash, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("No matching data found")
	} else if err != nil {
		return "", err
	}
	return data, nil
}

func AddDevice(phone_number, device_id string) bool {
	var keys = []string{
		fmt.Sprintf(deviceIdsByPhone, phone_number),
		fmt.Sprintf(phoneByDeviceId, device_id),
		allDeviceIds,
	}
	if err := add_device_script.Run(ctx, rdb, keys, phone_number, device_id).Err(); err != nil {
		log.Error("Add device failed: ", err.Error())

		return false
	}
	return true
}

func RemoveDevice(phone_number, device_id string) bool {
	var keys = []string{
		fmt.Sprintf(deviceIdsByPhone, phone_number),
		fmt.Sprintf(phoneByDeviceId, device_id),
		allDeviceIds,
	}
	if err := remove_device_script.Run(ctx, rdb, keys, device_id).Err(); err != nil {
		log.Error("Remove device failed: ", err.Error())
		return false
	}
	return true
}

func GetDevicesByPhone(phone_number string) ([]string, error) {
	var key = fmt.Sprintf(deviceIdsByPhone, phone_number)
	return rdb.SMembers(ctx, key).Result()
}

func GetPhoneByDevice(device_id string) (string, error) {
	var key = fmt.Sprintf(phoneByDeviceId, device_id)
	return rdb.Get(ctx, key).Result()
}

func GetAllDeviceIds() ([]string, error) {
	return rdb.SMembers(ctx, allDeviceIds).Result()
}
