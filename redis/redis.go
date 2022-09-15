package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	red "github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/mapping"

	"github.com/sunshibao/go-utils/logs"
)

const (
	ClusterType = "cluster"
	NodeType    = "node"
	Nil         = red.Nil

	blockingQueryTimeout = 5 * time.Second
	readWriteTimeout     = 2 * time.Second

	slowThreshold = time.Millisecond * 100
)

var ErrNilNode = errors.New("nil redis node")

type (
	Pair struct {
		Key   string
		Score int64
	}

	// thread-safe
	JingRedis struct {
		Addr string
		Type string
		Pass string
		brk  breaker.Breaker
	}

	RedisNode interface {
		red.Cmdable
	}

	Pipeliner = red.Pipeliner

	// Z represents sorted set member.
	Z = red.Z

	IntCmd   = red.IntCmd
	FloatCmd = red.FloatCmd
)

func NewJingRedis(redisAddr, redisType string, redisPass ...string) *JingRedis {
	var pass string
	for _, v := range redisPass {
		pass = v
	}

	return &JingRedis{
		Addr: redisAddr,
		Type: redisType,
		Pass: pass,
		brk:  breaker.NewBreaker(),
	}
}

func (s *JingRedis) Breaker() breaker.Breaker {
	return s.brk
}

func (s *JingRedis) GetRedisConn() (*red.Client, error) {
	conn, err := getRedis(s)
	if err != nil {
		return nil, err
	}
	return conn.(*red.Client), nil
}

// Use passed in redis connection to execute blocking queries
// Doesn't benefit from pooling redis connections of blocking queries
func (s *JingRedis) Blpop(redisNode RedisNode, key string) (string, error) {
	if redisNode == nil {
		return "", ErrNilNode
	}

	vals, err := redisNode.BLPop(blockingQueryTimeout, key).Result()
	if err != nil {
		return "", err
	}

	if len(vals) < 2 {
		return "", fmt.Errorf("no value on key: %s", key)
	} else {
		return vals[1], nil
	}
}

func (s *JingRedis) BlpopEx(redisNode RedisNode, key string) (string, bool, error) {
	if redisNode == nil {
		return "", false, ErrNilNode
	}

	vals, err := redisNode.BLPop(blockingQueryTimeout, key).Result()
	if err != nil {
		return "", false, err
	}

	if len(vals) < 2 {
		return "", false, fmt.Errorf("no value on key: %s", key)
	} else {
		return vals[1], true, nil
	}
}

func (s *JingRedis) GetRedisLock(key string) *RedisLock {
	return NewRedisLock(s, key)
}

func (s *JingRedis) DoWithRedisLocker(key string, fs ...func() error) error {
	if len(fs) == 0 {
		return nil
	}
	if locker := s.GetRedisLock(key); locker.TryLock() {
		defer locker.Unlock()
		for _, f := range fs {
			if err := f(); err != nil {
				return err
			}
		}
		return nil
	}
	return fmt.Errorf("DoWithRedisLocker(): Acquire redis locker %s failed", key)
}

func (s *JingRedis) MustDoWithRedisLocker(key string, fs ...func() error) error {
	if len(fs) > 0 {
		locker := s.GetRedisLock(key)
		locker.Lock()
		defer locker.Unlock()
		for _, f := range fs {
			if err := f(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *JingRedis) DoWith(fn ...func(node RedisNode) error) error {
	return s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}
		for _, f := range fn {
			err = f(conn)
			if err != nil {
				return err
			}
		}
		return nil
	}, acceptable)
}

func (s *JingRedis) Publish(channel string, message interface{}) error {
	return s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}
		return conn.Publish(channel, message).Err()
	}, acceptable)
}

func (s *JingRedis) PublishJSON(channel string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return s.Publish(channel, data)
}

type Subscriber interface {
	Channel() string
	Chan() <-chan string
	Subscribe() (string, bool)
	SubscribeFunc(func(string) error, ...int)
	SubscribeJSON(interface{}) error
	SubscribeJSONFunc(interface{}, func(interface{}) error, ...int)
}

type jingSubscriber struct {
	channel string
	sub     *red.PubSub
	ch      chan string
}

func (o *jingSubscriber) Channel() string {
	return o.channel
}

func (o *jingSubscriber) loop() {
	defer func() { close(o.ch) }()
	ch := o.sub.ChannelSize(1000)
	for {
		msg, ok := <-ch
		if !ok {
			return
		}
		if msg != nil {
			o.ch <- msg.Payload
		}
	}
}

func (o *jingSubscriber) Chan() <-chan string {
	return o.ch
}

func (o *jingSubscriber) Subscribe() (data string, ok bool) {
	data, ok = <-o.Chan()
	return
}

func (o *jingSubscriber) SubscribeFunc(f func(string) error, nn ...int) {
	n := 1
	if len(nn) > 0 && nn[0] > n {
		n = nn[0]
	}
	for i := 0; i < n; i++ {
		go o.call(f)
	}
}

func (o *jingSubscriber) SubscribeJSON(v interface{}) error {
	data, ok := o.Subscribe()
	if !ok {
		return Nil // redis connection closed
	}
	return json.Unmarshal([]byte(data), v)
}

func (o *jingSubscriber) SubscribeJSONFunc(arg interface{}, f func(interface{}) error, nn ...int) {
	rt := reflect.TypeOf(arg)
	if kind := rt.Kind(); kind != reflect.Ptr || !reflect.New(rt.Elem()).CanInterface() {
		panic("redis subscriber: invalid arguments pointer")
	}
	n := 1
	if len(nn) > 0 && nn[0] > n {
		n = nn[0]
	}
	for i := 0; i < n; i++ {
		go o.json(rt, f)
	}
}

func (o *jingSubscriber) call(f func(string) error) {
	for {
		data, ok := <-o.Chan()
		if !ok {
			return
		}
		if err := f(data); err != nil {
			logs.Errorf("Subscribe redis channel %s error: %s", o.Channel(), err)
		}
	}
}

func (o *jingSubscriber) json(t reflect.Type, f func(interface{}) error) {
	for {
		data, ok := <-o.Chan()
		if !ok {
			return
		}
		v := reflect.New(t.Elem()).Interface()
		if err := json.Unmarshal([]byte(data), v); err != nil {
			logs.Errorf("Subscribe redis channel %s error: %s", o.Channel(), err)
			continue
		}
		if err := f(v); err != nil {
			logs.Errorf("Subscribe redis channel %s error: %s", o.Channel(), err)
		}
	}
}

func (s *JingRedis) Subscriber(channel string) (Subscriber, error) {
	var sub *red.PubSub
	err := s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}
		switch c := conn.(type) {
		case *red.Client:
			sub = c.Subscribe()
		case *red.ClusterClient:
			sub = c.Subscribe()
		default:
			err = fmt.Errorf("unknown redis driver %T", conn)
		}
		return err
	}, acceptable)
	if err != nil {
		return nil, err
	}
	r := &jingSubscriber{channel: channel, sub: sub, ch: make(chan string, 5000)}
	go r.loop()
	return r, nil
}

func (s *JingRedis) Del(keys ...string) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.Del(keys...).Result(); err != nil {
			return err
		} else {
			val = int(v)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Delete(keys []string) (err error) {
	_, err = s.Del(keys...)
	return
}

func (s *JingRedis) Eval(script string, keys []string, args ...interface{}) (val interface{}, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.Eval(script, keys, args...).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) Exists(key string) (val bool, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.Exists(key).Result(); err != nil {
			return err
		} else {
			val = v == 1
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Expire(key string, seconds int) error {
	return s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		return conn.Expire(key, time.Duration(seconds)*time.Second).Err()
	}, acceptable)
}

func (s *JingRedis) Expireat(key string, expireTime int64) error {
	return s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		return conn.ExpireAt(key, time.Unix(expireTime, 0)).Err()
	}, acceptable)
}

func (s *JingRedis) Get(key string) (val string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if val, err = conn.Get(key).Result(); err == red.Nil {
			return nil
		} else if err != nil {
			return err
		} else {
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) GetBit(key string, offset int64) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.GetBit(key, offset).Result(); err != nil {
			return err
		} else {
			val = int(v)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Hdel(key, field string) (val bool, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.HDel(key, field).Result(); err != nil {
			return err
		} else {
			val = v == 1
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Hexists(key, field string) (val bool, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.HExists(key, field).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) Hget(key, field string) (val string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.HGet(key, field).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) Hgetall(key string) (val map[string]string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.HGetAll(key).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) Hincrby(key, field string, increment int) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.HIncrBy(key, field, int64(increment)).Result(); err != nil {
			return err
		} else {
			val = int(v)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Hkeys(key string) (val []string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.HKeys(key).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) Hlen(key string) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.HLen(key).Result(); err != nil {
			return err
		} else {
			val = int(v)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Hmget(key string, fields ...string) (val []string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.HMGet(key, fields...).Result(); err != nil {
			return err
		} else {
			val = toStrings(v)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Hset(key, field, value string) error {
	return s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		return conn.HSet(key, field, value).Err()
	}, acceptable)
}

func (s *JingRedis) Hsetnx(key, field, value string) (val bool, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.HSetNX(key, field, value).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) Hmset(key string, fieldsAndValues map[string]interface{}) error {
	return s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}
		return conn.HMSet(key, fieldsAndValues).Err()
	}, acceptable)
}

func (s *JingRedis) Hvals(key string) (val []string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.HVals(key).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) Incr(key string) (val int64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.Incr(key).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) Incrby(key string, increment int64) (val int64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.IncrBy(key, int64(increment)).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) Keys(pattern string) (val []string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.Keys(pattern).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) Llen(key string) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.LLen(key).Result(); err != nil {
			return err
		} else {
			val = int(v)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Lpop(key string) (val string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.LPop(key).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) Lpush(key string, values ...interface{}) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.LPush(key, values...).Result(); err != nil {
			return err
		} else {
			val = int(v)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Lrange(key string, start int, stop int) (val []string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.LRange(key, int64(start), int64(stop)).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) Lrem(key string, count int, value string) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.LRem(key, int64(count), value).Result(); err != nil {
			return err
		} else {
			val = int(v)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) LIndex(key string, index int64) (val string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err1 := getRedis(s)
		if err1 != nil {
			return err1
		}
		result, err2 := conn.LIndex(key, index).Result()
		if err2 != nil {
			if err2 == Nil {
				return nil
			}
			return err2
		}
		val = result
		return nil
	}, acceptable)
	return
}

func (s *JingRedis) Mget(keys ...string) (val []string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.MGet(keys...).Result(); err != nil {
			return err
		} else {
			val = toStrings(v)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Persist(key string) (val bool, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.Persist(key).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) Pfadd(key string, values ...interface{}) (val bool, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.PFAdd(key, values...).Result(); err != nil {
			return err
		} else {
			val = v == 1
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Pfcount(key string) (val int64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.PFCount(key).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) Pfmerge(dest string, keys ...string) error {
	return s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		_, err = conn.PFMerge(dest, keys...).Result()
		return err
	}, acceptable)
}

func (s *JingRedis) Ping() (val bool) {
	// ignore error, error means false
	_ = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			val = false
			return nil
		}

		if v, err := conn.Ping().Result(); err != nil {
			val = false
			return nil
		} else {
			val = v == "PONG"
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Pipelined(fn func(Pipeliner) error) (err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		_, err = conn.Pipelined(fn)
		return err

	}, acceptable)

	return
}

func (s *JingRedis) DoWithPipeliner(fs ...func(Pipeliner) error) error {
	if len(fs) == 0 {
		return nil
	}
	return s.Pipelined(func(p Pipeliner) error {
		for i, j := 0, len(fs); i < j; i++ {
			if err := fs[i](p); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *JingRedis) Rpush(key string, values ...interface{}) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.RPush(key, values...).Result(); err != nil {
			return err
		} else {
			val = int(v)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Sadd(key string, values ...interface{}) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.SAdd(key, values...).Result(); err != nil {
			return err
		} else {
			val = int(v)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Scan(cursor uint64, match string, count int64) (keys []string, cur uint64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		keys, cur, err = conn.Scan(cursor, match, count).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) SetBit(key string, offset int64, value int) error {
	return s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		_, err = conn.SetBit(key, offset, value).Result()
		return err
	}, acceptable)
}

func (s *JingRedis) Sscan(key string, cursor uint64, match string, count int64) (keys []string, cur uint64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		keys, cur, err = conn.SScan(key, cursor, match, count).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) Scard(key string) (val int64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.SCard(key).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) Set(key string, value string) error {
	return s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		return conn.Set(key, value, 0).Err()
	}, acceptable)
}

func (s *JingRedis) Setex(key, value string, seconds int) error {
	return s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		return conn.Set(key, value, time.Duration(seconds)*time.Second).Err()
	}, acceptable)
}

func (s *JingRedis) Setnx(key, value string) (val bool, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.SetNX(key, value, 0).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) SetnxEx(key, value string, seconds int) (val bool, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.SetNX(key, value, time.Duration(seconds)*time.Second).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) Sismember(key string, value interface{}) (val bool, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}
		val, err = conn.SIsMember(key, value).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) Srem(key string, values ...interface{}) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.SRem(key, values...).Result(); err != nil {
			return err
		} else {
			val = int(v)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Smembers(key string) (val []string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.SMembers(key).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) Spop(key string) (val string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.SPop(key).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) Srandmember(key string, count int) (val []string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.SRandMemberN(key, int64(count)).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) Sunion(keys ...string) (val []string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.SUnion(keys...).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) Sunionstore(destination string, keys ...string) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.SUnionStore(destination, keys...).Result(); err != nil {
			return err
		} else {
			val = int(v)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Sdiff(keys ...string) (val []string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.SDiff(keys...).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) Sdiffstore(destination string, keys ...string) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.SDiffStore(destination, keys...).Result(); err != nil {
			return err
		} else {
			val = int(v)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Ttl(key string) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if duration, err := conn.TTL(key).Result(); err != nil {
			return err
		} else {
			val = int(duration / time.Second)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Zadd(key string, score int64, value string) (val bool, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.ZAdd(key, red.Z{
			Score:  float64(score),
			Member: value,
		}).Result(); err != nil {
			return err
		} else {
			val = v == 1
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Zadds(key string, ps ...Pair) (val int64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		var zs []red.Z
		for _, p := range ps {
			z := red.Z{Score: float64(p.Score), Member: p.Key}
			zs = append(zs, z)
		}

		if v, err := conn.ZAdd(key, zs...).Result(); err != nil {
			return err
		} else {
			val = v
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Zcard(key string) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.ZCard(key).Result(); err != nil {
			return err
		} else {
			val = int(v)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Zcount(key string, start, stop int64) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.ZCount(key, strconv.FormatInt(start, 10),
			strconv.FormatInt(stop, 10)).Result(); err != nil {
			return err
		} else {
			val = int(v)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Zincrby(key string, increment int64, field string) (val int64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.ZIncrBy(key, float64(increment), field).Result(); err != nil {
			return err
		} else {
			val = int64(v)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Zscore(key string, value string) (val int64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.ZScore(key, value).Result(); err != nil {
			return err
		} else {
			val = int64(v)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Zrank(key, field string) (val int64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.ZRank(key, field).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) Zrem(key string, values ...interface{}) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.ZRem(key, values...).Result(); err != nil {
			return err
		} else {
			val = int(v)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Zremrangebyscore(key string, start, stop int64) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.ZRemRangeByScore(key, strconv.FormatInt(start, 10),
			strconv.FormatInt(stop, 10)).Result(); err != nil {
			return err
		} else {
			val = int(v)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Zremrangebyrank(key string, start, stop int64) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.ZRemRangeByRank(key, start, stop).Result(); err != nil {
			return err
		} else {
			val = int(v)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Zrange(key string, start, stop int64) (val []string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.ZRange(key, start, stop).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) ZrangeWithScores(key string, start, stop int64) (val []Pair, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.ZRangeWithScores(key, start, stop).Result(); err != nil {
			return err
		} else {
			val = toPairs(v)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) ZRevRangeWithScores(key string, start, stop int64) (val []Pair, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.ZRevRangeWithScores(key, start, stop).Result(); err != nil {
			return err
		} else {
			val = toPairs(v)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) ZrangebyscoreWithScores(key string, start, stop int64) (val []Pair, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.ZRangeByScoreWithScores(key, red.ZRangeBy{
			Min: strconv.FormatInt(start, 10),
			Max: strconv.FormatInt(stop, 10),
		}).Result(); err != nil {
			return err
		} else {
			val = toPairs(v)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) ZrangebyscoreWithScoresAndLimit(key string, start, stop int64, page, size int) (
	val []Pair, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		if size <= 0 {
			return nil
		}

		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.ZRangeByScoreWithScores(key, red.ZRangeBy{
			Min:    strconv.FormatInt(start, 10),
			Max:    strconv.FormatInt(stop, 10),
			Offset: int64(page * size),
			Count:  int64(size),
		}).Result(); err != nil {
			return err
		} else {
			val = toPairs(v)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) Zrevrange(key string, start, stop int64) (val []string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.ZRevRange(key, start, stop).Result()
		return err
	}, acceptable)

	return
}

func (s *JingRedis) ZrevrangebyscoreWithScores(key string, start, stop int64) (val []Pair, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.ZRevRangeByScoreWithScores(key, red.ZRangeBy{
			Min: strconv.FormatInt(start, 10),
			Max: strconv.FormatInt(stop, 10),
		}).Result(); err != nil {
			return err
		} else {
			val = toPairs(v)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) ZrevrangebyscoreWithScoresAndLimit(key string, start, stop int64, page, size int) (
	val []Pair, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		if size <= 0 {
			return nil
		}

		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if v, err := conn.ZRevRangeByScoreWithScores(key, red.ZRangeBy{
			Min:    strconv.FormatInt(start, 10),
			Max:    strconv.FormatInt(stop, 10),
			Offset: int64(page * size),
			Count:  int64(size),
		}).Result(); err != nil {
			return err
		} else {
			val = toPairs(v)
			return nil
		}
	}, acceptable)

	return
}

func (s *JingRedis) String() string {
	return s.Addr
}

func (s *JingRedis) scriptLoad(script string) (string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return "", err
	}

	return conn.ScriptLoad(script).Result()
}

func acceptable(err error) bool {
	return err == nil || err == red.Nil
}

func GetRedis(r *JingRedis) (RedisNode, error) {
	return getRedis(r)
}

func getRedis(r *JingRedis) (RedisNode, error) {
	switch r.Type {
	case ClusterType:
		return getCluster(r.Addr, r.Pass)
	case NodeType:
		return getClient(r.Addr, r.Pass)
	default:
		return nil, fmt.Errorf("redis type '%s' is not supported", r.Type)
	}
}

func toPairs(vals []red.Z) []Pair {
	pairs := make([]Pair, len(vals))
	for i, val := range vals {
		switch member := val.Member.(type) {
		case string:
			pairs[i] = Pair{
				Key:   member,
				Score: int64(val.Score),
			}
		default:
			pairs[i] = Pair{
				Key:   mapping.Repr(val.Member),
				Score: int64(val.Score),
			}
		}
	}
	return pairs
}

func toStrings(vals []interface{}) []string {
	ret := make([]string, len(vals))
	for i, val := range vals {
		if val == nil {
			ret[i] = ""
		} else {
			switch val := val.(type) {
			case string:
				ret[i] = val
			default:
				ret[i] = mapping.Repr(val)
			}
		}
	}
	return ret
}
