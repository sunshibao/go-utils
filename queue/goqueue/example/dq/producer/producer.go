package main

import (
	"fmt"
	"time"

	"github.com/zeromicro/go-queue/dq"
)

func main() {
	producer := dq.NewProducer([]dq.Beanstalk{
		{
			Endpoint: "127.0.0.1:11300",
			Tube:     "tube",
		},
		{
			Endpoint: "127.0.0.1:11301",
			Tube:     "tube",
		},
	})
	for i := 1000; i < 1005; i++ {
		_, err := producer.Delay([]byte("222"), time.Second*5)
		if err != nil {
			fmt.Println(err)
		}
	}
}
