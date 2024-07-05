package redis_cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-instaloader/utils/fwRedis"
	"go-instaloader/utils/rlog"
	"reflect"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type Key interface{}
type Value interface{}
type LoaderFunc func(Key) (Value, error)
type Option func(c *RedisCache)

func NewRedisCacheWithOption(keyPrefix string, loader LoaderFunc, options ...Option) *RedisCache {
	c := &RedisCache{
		loader:    loader,
		keyPrefix: keyPrefix,
	}
	for _, opt := range options {
		opt(c)
	}

	return c
}

func WithExpireAfterWrite(expire time.Duration) Option {
	if expire < 0 {
		expire = 0
	}

	return func(c *RedisCache) {
		c.expire = expire
	}
}

type RedisCache struct {
	loader    LoaderFunc
	keyPrefix string
	expire    time.Duration
}

func (c *RedisCache) GetIfPresent(id int64) (interface{}, bool) {
	if id <= 0 {
		return nil, false
	}

	key := c.getKey(id)
	ctx := context.Background()
	getCmd := fwRedis.RedisQueue().Get(ctx, key)
	// The data read by redis is nil or read error
	if err := getCmd.Err(); err != nil && err != redis.Nil {
		return nil, false
	}

	jsonStr := getCmd.Val()
	if jsonStr != "" {
		obj := reflect.New(reflect.TypeOf(c.loader)).Interface()
		if err := JsonUnmarshal(jsonStr, obj); err == nil {
			return obj, true // Get the correct cache data
		} else {
			// Record stack
			rlog.Error(err) // The obtained data is deserialized incorrectly. After outputting the log, continue to reload the subsequent logic.
		}
	}

	return nil, false
}
func (c *RedisCache) Put(id int64, input interface{}) error {
	if id <= 0 {
		return errors.New("id cannot be less than or equal to 0")
	}

	key := c.getKey(id)
	return c.set(key, input)
}

func (c *RedisCache) getKey(id int64) string {
	key := fmt.Sprintf("%s:%v", c.keyPrefix, id)
	return key
}

func (c *RedisCache) getKeyString(s string) string {
	key := fmt.Sprintf("%s:%s", c.keyPrefix, s)
	return key
}

func (c *RedisCache) set(key string, input interface{}) error {
	jStr, err := FormatJson(input)
	if err != nil {
		return err
	}
	ctx := context.Background()
	cmd := fwRedis.RedisQueue().Set(ctx, key, jStr, c.expire)
	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}

func (c *RedisCache) Get(id int64, out interface{}) (interface{}, error) {
	if id <= 0 {
		return nil, errors.New("id cannot be less than or equal to 0")
	}

	key := c.getKey(id)
	ctx := context.Background()
	getCmd := fwRedis.RedisQueue().Get(ctx, key)
	// The data read by redis is nil or read error
	if err := getCmd.Err(); err != nil && err != redis.Nil {
		return nil, err
	}

	jsonStr := getCmd.Val()
	if jsonStr != "" {
		if err := JsonUnmarshal(jsonStr, out); err == nil {
			return out, nil // Get the correct cache data
		} else {
			rlog.Error(err, string(debug.Stack())) // The obtained data is deserialized incorrectly. After outputting the log, continue to reload the subsequent logic.
		}
	}

	// Read data from mysql again
	if c.loader == nil {
		return nil, errors.New("loader is nil")
	}
	val, err := c.loader(id)
	if err != nil {
		return nil, err
	}

	if val != nil {
		if err := c.set(key, val); err != nil {
			// The cache fails and an error is output. mysql hard resistance
			rlog.Error("caching failed：%v", val)
		}
	}
	return val, err

}

func (c *RedisCache) GetString(s string, out interface{}) (interface{}, error) {
	if s == "" {
		return nil, errors.New("key cannot be empty")
	}

	key := c.getKeyString(s)
	ctx := context.Background()
	getCmd := fwRedis.RedisQueue().Get(ctx, key)
	// The data read by redis is nil or read error
	if err := getCmd.Err(); err != nil && err != redis.Nil {
		return nil, err
	}

	jsonStr := getCmd.Val()
	if jsonStr != "" {
		if err := JsonUnmarshal(jsonStr, out); err == nil {
			return out, nil // Get the correct cache data
		} else {
			rlog.Error(err, string(debug.Stack())) // The obtained data is deserialized incorrectly. After outputting the log, continue to reload the subsequent logic.
		}
	}

	// Read data from mysql again
	if c.loader == nil {
		return nil, errors.New("loader is nil")
	}
	val, err := c.loader(s)
	if err != nil {
		return nil, err
	}

	if val != nil {
		if err := c.set(key, val); err != nil {
			// The cache fails and an error is output. mysql hard resistance
			rlog.Error("caching failed：%v", val)
		}
	}
	return val, err

}

func (c *RedisCache) GetByKey(key string) (interface{}, error) {
	ctx := context.Background()
	getCmd := fwRedis.RedisQueue().Get(ctx, key)
	// The data read by redis is nil or read error
	if err := getCmd.Err(); err != nil && err != redis.Nil {
		return nil, err
	}

	jsonStr := getCmd.Val()
	out := interface{}(nil)
	if jsonStr != "" {
		if err := JsonUnmarshal(jsonStr, &out); err == nil {
			return out, nil // Get the correct cache data
		} else {
			rlog.Error(err) // The obtained data is deserialized incorrectly. After outputting the log, continue to reload the subsequent logic.
		}
	}

	// Read data from mysql again
	if c.loader == nil {
		return nil, errors.New("loader is nil")
	}
	val, err := c.loader(key)
	if err != nil {
		return nil, err
	}

	if val != nil {
		if err := c.set(key, val); err != nil {
			// The cache fails and an error is output. mysql hard resistance
			rlog.Error("caching failed：%v", val)
		}
	}

	return val, err

}
func (c *RedisCache) Del(key string) error {
	ctx := context.Background()
	_, err := fwRedis.RedisQueue().Del(ctx, key).Result()
	return err
}
func (c *RedisCache) Invalidate(id int64) error {
	key := c.getKey(id)
	ctx := context.Background()
	_, err := fwRedis.RedisQueue().Del(ctx, key).Result()
	return err
}

func (c *RedisCache) InvalidateString(username string) error {
	key := c.getKeyString(username)
	ctx := context.Background()
	_, err := fwRedis.RedisQueue().Del(ctx, key).Result()
	return err
}

func (c *RedisCache) InvalidateAll() error {
	key := fmt.Sprintf("%s:*", c.keyPrefix)
	ctx := context.Background()
	_, err := fwRedis.RedisQueue().Del(ctx, key).Result()
	return err
}

// Obtain data in batches [MGet does not support Cluster, and the logic needs to be modified during upgrade]
func (c *RedisCache) GetArray(ids []int64, outType reflect.Type) (outList []interface{}, failListIds []int64, err error) {
	var keys []string

	for _, id := range ids {
		keys = append(keys, c.getKey(id))
	}
	ctx := context.Background()
	values, err := fwRedis.RedisQueue().MGet(ctx, keys...).Result()
	if err != nil {
		return nil, nil, err
	}

	for i, val := range values {
		id := ids[i]
		if val != nil {
			jsonStr := val.(string)
			if jsonStr != "" {
				obj := reflect.New(outType).Interface()
				if err := JsonUnmarshal(jsonStr, obj); err == nil {
					outList = append(outList, obj)
					continue
				}
			}
		}

		obj := reflect.New(outType).Interface()
		obj, err := c.Get(id, obj)
		if err != nil {
			failListIds = append(failListIds, id)
			continue
		} else {
			outList = append(outList, obj)
			continue
		}
	}

	return outList, failListIds, nil
}

func (c *RedisCache) Exists(id int64) bool {
	key := c.getKey(id)
	ctx := context.Background()
	ret := fwRedis.RedisQueue().Exists(ctx, key).Val()
	return ret == 1
}

func JsonUnmarshal(str string, data interface{}) error {
	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Ptr {
		return errors.New("parameter structure must be a pointer")
	}

	dec := json.NewDecoder(strings.NewReader(str))
	dec.DisallowUnknownFields()
	if err := dec.Decode(data); err != nil {
		msg := fmt.Sprintf("type:[%s] err:[%+v] json:%+v ", strings.ToLower(t.String()), err, str)
		return errors.New(msg)
	}

	return nil
}

func FormatJson(obj interface{}) (str string, err error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return
	}
	str = string(data)
	return
}
