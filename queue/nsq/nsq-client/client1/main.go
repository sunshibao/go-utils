package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nsqio/go-nsq"
)

// NSQ Producer(生产者)示例演示
var Producer *nsq.Producer

type NsqProducer struct {
	nsqd_host string
	nsqd_port int
}

// 初始化NSQ生产者
func (NP NsqProducer) InitProducer() (*nsq.Producer, error) {
	// NewConfig返回一个新的默认nsq配置。
	config := nsq.NewConfig()
	// 组合nsqd服务连接地址。
	nsqdAddr := fmt.Sprintf("%s:%d", NP.nsqd_host, NP.nsqd_port)
	Producer, err := nsq.NewProducer(nsqdAddr, config)
	if err != nil {
		fmt.Printf("create producer failed, err:%v\n", err)
		return nil, err
	}
	return Producer, nil
}

func main() {
	// 1.实例化以及初始化
	nsqd := &NsqProducer{
		nsqd_host: "127.0.0.1",
		nsqd_port: 4150,
	}
	Producer, err := nsqd.InitProducer()
	if err != nil {
		fmt.Printf("Init producer failed, err:%v\n", err)
		return
	} else {
		log.Printf("Init Producer success!")
	}

	// 2.从标准输入读取
	fmt.Println("请输入你要向Topic_Demo消息队列传递的消息:")
	reader := bufio.NewReader(os.Stdin)
	for {
		data, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("read string from stdin failed, err:%v\n", err)
			continue
		}
		data = strings.TrimSpace(data)
		// 3.当前输入Q或者q时退出程序
		if strings.ToUpper(data) == "Q" {
			break
		}
		// 4.向 'Topic_Demo' publish 数据
		err = Producer.Publish("Topic_Demo", []byte(data))
		if err != nil {
			fmt.Printf("publish msg to nsq failed, err:%v\n", err)
			continue
		}
		// 5.向 'Topic_Demo' DeferredPublish 延时发送数据  目前nsq的延时消息只支持1小时以内。
		err = Producer.DeferredPublish("Topic_Demo", 3600*time.Second, []byte(data))
		if err != nil {
			fmt.Printf("publish msg to nsq failed, err:%v\n", err)
			continue
		}
	}
}
