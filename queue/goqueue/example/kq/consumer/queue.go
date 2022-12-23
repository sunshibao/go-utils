package main

import (
	"fmt"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/conf"
)

func main() {
	var c kq.KqConf
	conf.MustLoad("/Users/sunshibao/Desktop/www/go/src/go-utils/queue/goqueue/example/kq/consumer/config.yaml", &c)

	q := kq.MustNewQueue(c, kq.WithHandle(func(k, v string) error {
		fmt.Printf("=> %s\n", v)
		return nil
	}))
	defer q.Stop()
	q.Start()
}
