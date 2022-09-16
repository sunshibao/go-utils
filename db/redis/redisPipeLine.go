package redis

import (
	"fmt"
	"github.com/go-redis/redis"
)

//RedisPipeLine 管道的使用 批量执行，不保证原子性
//使用pipeline组装的命令个数不能太多，不然数据量过大，增加客户端的等待时间，还可能造成网络阻塞，可以将大量命令的拆分多个小的pipeline命令完成。
//高频命令场景下，应避免使用管道，因为需要先将全部执行命令放入管道，会耗时。另外，需要使用返回值的情况也不建议使用， 同步所有管道返回结果也是个耗时的过程！管道无法提供原子性/事务保障。
func RedisPipeLine() {
	redisDB := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	var (
		zreCmd     *redis.Cmd
		ZscCmd     *redis.Cmd
		myScore    int64
		myGouLiang int64
	)
	if _, err := redisDB.Pipelined(func(pipe redis.Pipeliner) error {
		zreCmd = pipe.Do("Incrby", "aaaaa", "1")
		ZscCmd = pipe.Do("Incrby", "bbbbb", "2")
		return nil
	}); err != nil && err != redis.Nil {
		return
	}

	if zreCmd.Err() == nil {
		myScore, _ = zreCmd.Int64()
		myScore += 1
	}
	if ZscCmd.Err() == nil {
		myGouLiang, _ = ZscCmd.Int64()
	}

	fmt.Println(myGouLiang)
	fmt.Println(myScore)
}
