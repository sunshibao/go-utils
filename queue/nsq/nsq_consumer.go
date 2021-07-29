package nsq

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/nsqio/go-nsq"
)

// ----------------------------------------
//  NSQ 队列消息（与 nsq.Message 一致）
// ----------------------------------------
type NSQMessage struct {
	*nsq.Message
}

// ----------------------------------------
//  NSQ 队列消息处理器
// ----------------------------------------
type NSQHandler interface {
	Handle(*NSQMessage) error
}

// ----------------------------------------
//  NSQ 队列消息处理器函数
// ----------------------------------------
type NSQHandlerFunc func(*NSQMessage) error

// 处理 NSQ 队列消息，这是对 NSQHandler 接口的实现
func (f NSQHandlerFunc) Handle(message *NSQMessage) error {
	return f(message)
}

// ----------------------------------------
//  NSQ 队列消息处理器桥
// ----------------------------------------
type NSQHandlerBridge struct {
	handler NSQHandler
}

// 对 nsq.Handler 的实现
func (bridge *NSQHandlerBridge) HandleMessage(message *nsq.Message) error {
	return bridge.handler.Handle(&NSQMessage{Message: message})
}

// ----------------------------------------
//  NSQ 队列消费者
// ----------------------------------------
type NSQConsumer struct {
	mu        *sync.Mutex     // 操作互斥锁
	config    *NSQConfig      // 配置项
	consumers []*nsq.Consumer // 消费者实例列表
}

// 创建一个 NSQ 队列消费者实例
func CreateNSQConsumer(config *NSQConfig) *NSQConsumer {
	return &NSQConsumer{
		mu:     &sync.Mutex{},
		config: config.Copy(),
	}
}

// 开始使用给定的消费者从队列中消费消息
func (er *NSQConsumer) Consume(handler NSQHandler) error {
	er.mu.Lock()
	defer er.mu.Unlock()
	// 禁止在未关闭前重复调用消费
	if len(er.consumers) != 0 {
		return fmt.Errorf("consumers are already running")
	}
	if handler == nil {
		return fmt.Errorf("nil nsq message handler")
	}
	// 检查消费者客户端数量和并发数配置，以确定是否使用默认值
	count, concurrent := er.config.ConsumerCount, er.config.ConsumerConcurrent
	if count <= 0 {
		count = 1
	}
	if concurrent <= 0 {
		concurrent = 1
	}
	cfg := nsq.NewConfig()
	// 断开重连设置成 500 毫秒（默认 60 秒太长了）
	cfg.LookupdPollInterval = time.Millisecond * 500
	// 设置消费者端允许等待中的消息最大数量为并发数的 2 倍（经验值）
	cfg.MaxInFlight = concurrent * 2
	// 设置 msg timeout
	if er.config.MsgTimeout != "" {
		msgTimeoutDuration, err := time.ParseDuration(er.config.MsgTimeout)
		if err != nil {
			return fmt.Errorf("msg timeout duration err: %s", er.config.MsgTimeout)
		}
		cfg.MsgTimeout = msgTimeoutDuration
	}
	// 依次创建消费者并开始消费
	addresses := strings.Split(er.config.Address, ",")
	for i, j := 0, len(addresses); i < j; i++ {
		addresses[i] = strings.TrimSpace(addresses[i])
		if addresses[i] == "" {
			return fmt.Errorf("the config address item value can not be empty")
		}
	}
	for i := 0; i < count; i++ {
		if c, err := nsq.NewConsumer(er.config.Topic, er.config.Channel, cfg); err != nil {
			return fmt.Errorf("create nsq consumer error: %s", err)
		} else {
			c.SetLogger(nil, nsq.LogLevelError) // TODO: Set logger ...
			c.AddConcurrentHandlers(&NSQHandlerBridge{handler: handler}, concurrent)
			if er.config.Lookup {
				if err := c.ConnectToNSQLookupds(addresses); err != nil {
					er.close()
					return fmt.Errorf("connect to nsq lookupd error: %s", err)
				}
			} else {
				if err := c.ConnectToNSQDs(addresses); err != nil {
					er.close()
					return fmt.Errorf("connect to nsqd error: %s", err)
				}
			}
			er.consumers = append(er.consumers, c)
		}
	}
	return nil
}

// 阻塞当前工作协程，等待 Close 方法被调用
func (er *NSQConsumer) Wait() {
	wg := &sync.WaitGroup{}
	for i, j := 0, len(er.consumers); i < j; i++ {
		wg.Add(1)
		// 这里通过监听每一个消费者的 StopChan 关闭状态来确定消费者是否已经退出
		go func(c *nsq.Consumer, g *sync.WaitGroup) {
			<-c.StopChan
			g.Done()
		}(er.consumers[i], wg)
	}
	wg.Wait()
}

// 关闭所有的消费者（无锁）
func (er *NSQConsumer) close() {
	for i, j := 0, len(er.consumers); i < j; i++ {
		er.consumers[i].Stop()
	}
}

// 关闭所有的消费者
func (er *NSQConsumer) Close() {
	er.mu.Lock()
	er.close()
	er.mu.Unlock()
}
