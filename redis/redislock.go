package redis

import (
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"

	red "github.com/go-redis/redis"

	"github.com/sunshibao/go-utils/logs"
)

const (
	letters     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lockCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
    return "OK"
else
    return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
end`
	delCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`
	randomLen       = 16
	tolerance       = 500 // milliseconds
	millisPerSecond = 1000
)

type RedisLock struct {
	store   *JingRedis
	seconds uint32
	key     string
	id      string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewRedisLock(store *JingRedis, key string) *RedisLock {
	return &RedisLock{
		store:   store,
		seconds: 30,
		key:     key,
		id:      randomStr(randomLen),
	}
}

func (rl *RedisLock) Acquire() (bool, error) {
	seconds := atomic.LoadUint32(&rl.seconds)
	resp, err := rl.store.Eval(lockCommand, []string{rl.key}, []string{
		rl.id, strconv.Itoa(int(seconds)*millisPerSecond + tolerance)})
	if err == red.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	} else if resp == nil {
		return false, nil
	}

	reply, ok := resp.(string)
	if ok && reply == "OK" {
		return true, nil
	} else {
		return false, nil
	}
}

func (rl *RedisLock) Release() (bool, error) {
	resp, err := rl.store.Eval(delCommand, []string{rl.key}, []string{rl.id})
	if err != nil {
		return false, err
	}

	if reply, ok := resp.(int64); !ok {
		return false, nil
	} else {
		return reply == 1, nil
	}
}

func (rl *RedisLock) TryLock() bool {
	for i := 0; i < 10; i++ {
		ok, err := rl.Acquire()
		if ok {
			return true
		}
		if err != nil {
			logs.Errorf("RedisLock.Acquire failed (%d/10): %s", i+1, err)
		}
		time.Sleep(time.Millisecond * 5)
	}
	return false
}

func (rl *RedisLock) Lock() {
	for {
		if rl.TryLock() {
			return
		}
	}
}

func (rl *RedisLock) Unlock() {
	for i := 0; i < 10; i++ {
		if _, err := rl.Release(); err == nil {
			return
		} else {
			logs.Errorf("RedisLock.Release failed (%d/10): %s", i+1, err)
			time.Sleep(time.Millisecond * 5)
		}
	}
}

func (rl *RedisLock) SetExpire(seconds int) {
	atomic.StoreUint32(&rl.seconds, uint32(seconds))
}

func randomStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
