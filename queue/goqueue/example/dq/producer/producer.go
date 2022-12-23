package main

import (
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
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
	//延时发送
	for i := 1000; i < 1005; i++ {
		_, err := producer.Delay([]byte("延时发送："+gconv.String(i)), time.Second*5)
		if err != nil {
			fmt.Println(err)
		}
	}
	//定时发送，设置时间是一定要用ParseInLocation 不然容易出现时区问题
	parseTime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2022-12-23 18:07:05", time.Local)
	_, err := producer.At([]byte("定时发送"), parseTime)
	if err != nil {
		fmt.Println(err)
	}

}
