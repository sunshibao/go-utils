package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"strings"
	"time"
)

func main() {
	redisDB := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	// 开启一个TxPipeline事务
	pipe := redisDB.TxPipeline()
	defer pipe.Close()

	// 这里是放命令
	pipe.SetNX("Freeze:DL201544a00007:a00002:a1:300", "200", 30*time.Second)
	pipe.SetNX("Freeze:DL201544a00008:a00002:a1:400", "400", 30*time.Second)
	pipe.SetNX("Freeze:DL201544a00009:a00002:a1:500", "500", 30*time.Second)
	pipe.SetNX("Freeze:DL201544a0000991123:a000091123:a1:592100", "5923100", 30*time.Second)
	// 通过Exec函数提交redis事务
	r, err := pipe.Exec()
	if err != nil {
		// 取消提交
		pipe.Discard()
	}
	// 这里调用exec执行刚刚加的命令，redis的事务和mysql事务是不一样的，一般情况下，这里的err出错是在编译期间出错，运行期间是很少出错的
	// mysql的事务是具有原子性，一致性，隔离性 ，持久性

	// redis事务三阶段：
	//
	// 开启：以MULTI开始一个事务
	// 入队：将多个命令入队到事务中，接到这些命令并不会立即执行，而是放到等待执行的事务队列里面
	// 执行：由EXEC命令触发事务
	// redis事务三大特性：
	//
	// 单独的隔离操作：事务中的所有命令都会序列化、按顺序地执行。事务在执行的过程中，不会被其他客户端发送来的命令请求所打断。
	// 没有隔离级别的概念：队列中的命令没有提交之前都不会实际的被执行，因为事务提交前任何指令都不会被实际执行，也就不存在”事务内的查询要看到事务里的更新，在事务外查询不能看到”这个让人万分头痛的问题
	// 不保证原子性：redis同一个事务中如果有一条命令执行失败，其后的命令仍然会被执行，没有回滚
	//
	// 所以，如果出现，第一条数据处理是true，其他的处理是false，那需要把其中的为true的数据，撤销其操作，这里也只能手动去撤销

	var resultmap []map[string]string
	resultmap = make([]map[string]string, 0)
	for _, v := range r {
		params := fmt.Sprintf("%s", v)
		res := strings.Split(params, " ")
		fmt.Println("key=", res[1])
		// 处理结果
		fmt.Println("res=", res[6])

		if res[6] == "true" {
			var model map[string]string
			model = make(map[string]string, 0)
			model["key"] = res[1]
			model["result"] = res[6]
			resultmap = append(resultmap, model)
		}
	}
	// 这堆代码是为了事务处理结果不一致导致的问题
	if len(r) != len(resultmap) {
		for _, vb := range resultmap {
			redisDB.Del(vb["key"]).Result()
		}
	}
}
