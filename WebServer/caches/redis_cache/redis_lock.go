package redis_cache

import (
	"context"
	"fmt"
	"go-instaloader/utils/fwRedis"
	"go-instaloader/utils/rlog"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	default_timeout               = 60
	default_redis_pool_max_idle   = 10
	default_redis_pool_max_active = 20
)

type Lock struct {
	resource string
	token    string
	timeout  int
}

func (lock *Lock) tryLock() (ok bool, err error) {
	ctx := context.Background()
	cmd := fwRedis.RedisQueue().SetNX(ctx, lock.key(), lock.token, time.Duration(lock.timeout)*time.Second)
	err = cmd.Err()
	if err == redis.Nil {
		// The lock was not successful, it already exists.
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (lock *Lock) UnlockDeferDefault() (err error) {
	time.Sleep(time.Duration(default_timeout) * time.Second)
	ctx := context.Background()
	err = fwRedis.RedisQueue().Del(ctx, lock.key()).Err()
	return
}

func (lock *Lock) Unlock() (err error) {
	ctx := context.Background()
	err = fwRedis.RedisQueue().Del(ctx, lock.key()).Err()
	return
}

func (lock *Lock) key() string {
	return fmt.Sprintf("riskcontrol:redislock:%s", lock.resource)
}

func (lock *Lock) AddTimeout(ex_time int64) (ok bool, err error) {
	ctx := context.Background()
	cmd := fwRedis.RedisQueue().TTL(ctx, lock.key())
	err = cmd.Err()
	if err != nil {
		rlog.Error("redis get failed:", err)
	}

	TTL := cmd.Val()
	if TTL > 0 {
		setCmd := fwRedis.RedisQueue().Set(ctx, lock.key(), lock.token, TTL+time.Duration(ex_time)*time.Second)
		err = setCmd.Err()
		if err == redis.Nil {
			return false, nil
		}
		if err != nil {
			return false, err
		}
	}
	return false, nil
}

func TryLock(resource string, token string) (lock *Lock, ok bool, err error) {
	return TryLockWithTimeout(resource, token, default_timeout)
}

func TryLockWithTimeout(resource string, token string, timeout int) (lock *Lock, ok bool, err error) {
	lock = &Lock{resource, token, timeout}

	ok, err = lock.tryLock()

	if !ok || err != nil {
		lock = nil
	}

	return
}
