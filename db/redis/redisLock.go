package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

//分布式锁
//- 互斥——同时刻只能有一个持有者；
//- 可重入——同一个持有者可多次进入；
//- 阻塞——多个锁请求者时，未获得锁的阻塞等待；
//- 无死锁——持有锁的客户端崩溃（crashed)或者网络被分裂（gets partitioned)，锁仍然可以被获取；
//- 容错——只要大部分节点正常，仍然可以获取和释放锁；

type redisClient redis.Client

func connRedisCluster(address []string, password string) *redis.ClusterClient {
	conf := redis.ClusterOptions{
		Addrs:    address,
		Password: password,
	}
	return redis.NewClusterClient(&conf)
}

func connRedisSingle(addr, password string) *redis.Client {
	conf := redis.Options{
		Addr:     addr,
		Password: password,
	}
	return redis.NewClient(&conf)
}

func (r *redisClient) lock(value string) (error, bool) {
	ret := r.SetNX("hello", value, time.Second*10)
	if err := ret.Err(); err != nil {
		fmt.Printf("set value %s error: %v\n", value, err)
		return err, false
	}
	return nil, ret.Val()
}

func (r *redisClient) unlock() bool {
	ret := r.Del("hello")
	if err := ret.Err(); err != nil {
		fmt.Println("unlock error: ", err)
		return false
	}
	return true
}

func (r *redisClient) retryLock() bool {
	ok := false
	for !ok {
		err, t := r.getTTL()
		if err != nil {
			return false
		}
		if t > 0 {
			fmt.Printf("锁被抢占, %f 秒后重试...\n", (t / 10).Seconds())
			time.Sleep(t / 10)
		}
		err, ok = r.lock("Jan")
		if err != nil {
			return false
		}
	}
	return ok
}

func (r *redisClient) getLock() (error, string) {
	ret := r.Get("hello")
	if err := ret.Err(); err != nil {
		fmt.Println("get lock error: ", err)
		return err, ""
	}
	rt, _ := ret.Bytes()
	return nil, string(rt)
}

// 获取锁的过期剩余时间
func (r *redisClient) getTTL() (error, time.Duration) {
	ret := r.TTL("hello")
	if err := ret.Err(); err != nil {
		fmt.Println("get TTL error: ", err)
		return err, 0
	}
	return nil, ret.Val()
}

func (r *redisClient) threadLock(threadId string) {
	for {
		err, _ := r.getLock()
		if err != nil && err.Error() == "redis: nil" {
			// 没有获取到值，说明目前没有人持有锁
			fmt.Printf("线程 %s 开始加锁\n", threadId)
			err, ok := r.lock("Jan")
			if err != nil {
				return
			}
			if !ok {
				if !r.retryLock() {
					fmt.Printf("线程 %s 加锁失败\n", threadId)
					return
				}
			}
			fmt.Printf("线程 %s 已加锁\n", threadId)
			// 加锁后执行相应操作
			time.Sleep(5 * time.Second)
			// 释放锁
			r.unlock()
			fmt.Printf("线程 %s 已释放锁\n", threadId)
			return
		} else if err != nil {
			return
		}
		err, t := r.getTTL()
		if err != nil {
			return
		}
		if t > 0 {
			fmt.Printf("线程 %s 锁被占用, %f 秒后重试\n", threadId, (t / 10).Seconds())
			time.Sleep(t / 10)
		}
	}
}

func main() {
	var r redisClient
	address := "127.0.0.1:6379"
	cl := connRedisSingle(address, "")
	defer cl.Close()
	r = redisClient(*cl)
	// 线程1获取锁
	go r.threadLock("1")
	// time.Sleep(10 * time.Millisecond)
	// 线程2获取锁
	go r.threadLock("2")
	select {}
}
