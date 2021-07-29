package nsq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/nsqio/go-nsq"
)

// NSQProducer
type NSQProducer struct {
	mu       *sync.Mutex   // 操作互斥锁
	config   *NSQConfig    // 配置项
	producer *nsq.Producer // 生产者实例
	address  string        // 正在使用的服务地址
}

// CreateNSQProducer 创建一个 NSQ 队列生产者实例
func CreateNSQProducer(config *NSQConfig) *NSQProducer {
	producer := &NSQProducer{
		mu:     &sync.Mutex{},
		config: config.Copy(),
	}
	if err := producer.Connect(); err != nil {
		panic(err)
	}
	return producer
}

// Connect 连接当前生产者到远程服务器
func (er *NSQProducer) Connect() error {
	er.mu.Lock()
	defer er.mu.Unlock()
	if er.producer != nil {
		return nil
	}
	cfg := nsq.NewConfig()
	// 断开重连设置成 500 毫秒（默认 60 秒太长了）
	cfg.LookupdPollInterval = time.Millisecond * 500
	// 从多个地址中随机选取一个使用
	addresses := strings.Split(er.config.Address, ",")
	if count := len(addresses); count > 1 {
		er.address = strings.TrimSpace(addresses[rand.Intn(count)])
	} else {
		er.address = strings.TrimSpace(addresses[0])
	}
	// 创建生产者
	if p, err := nsq.NewProducer(er.address, cfg); err != nil {
		return fmt.Errorf("create nsq producer error: [%s] %s", er.address, err)
	} else {
		// 可达性检查
		if err := p.Ping(); err != nil {
			return fmt.Errorf("ping nsq producer error: [%s] %s", er.address, err)
		}
		er.producer = p
		return nil
	}
}

// DeferredPublishString 推送一个延时消息
func (er *NSQProducer) DeferredPublishString(delay time.Duration, message string) error {
	return er.DeferredPublish(delay, []byte(message))
}

// DeferredPublishJSON 推送一个json延时消息
func (er *NSQProducer) DeferredPublishJSON(delay time.Duration, message interface{}) error {
	if bs, err := json.Marshal(message); err != nil {
		return fmt.Errorf("json encode nsq message error: [%s] %s", er.address, err)
	} else {
		return er.DeferredPublish(delay, bs)
	}
}

// DeferredPublishStringWithTopic 推送一个string类型延时消息到指定topic
func (er *NSQProducer) DeferredPublishStringWithTopic(topic string, delay time.Duration, message string) error {
	return er.DeferredPublishWithTopic(topic, delay, []byte(message))
}

// DeferredPublishJSONWithTopic 推送一个json类型延时消息到指定topic
func (er *NSQProducer) DeferredPublishJSONWithTopic(topic string, delay time.Duration, message interface{}) error {
	if bs, err := json.Marshal(message); err != nil {
		return fmt.Errorf("json encode nsq message error: [%s] %s", er.address, err)
	} else {
		return er.DeferredPublishWithTopic(topic, delay, bs)
	}
}

// PublishString 立即推送一个字符串消息到队列中
func (er *NSQProducer) PublishString(message string) error {
	return er.Publish(bytes.NewBufferString(message).Bytes())
}

// PublishStringWithTopic 立即推送一个字符串消息到队列中
func (er *NSQProducer) PublishStringWithTopic(topic string, message string) error {
	return er.PublishWithTopic(topic, bytes.NewBufferString(message).Bytes())
}

// PublishJSON 立即推送一个 JSON 消息到队列中
func (er *NSQProducer) PublishJSON(message interface{}) error {
	if bs, err := json.Marshal(message); err != nil {
		return fmt.Errorf("json encode nsq message error: [%s] %s", er.address, err)
	} else {
		return er.Publish(bs)
	}
}

// PublishJSONWithTopic 立即推送一个 JSON 消息到队列中
func (er *NSQProducer) PublishJSONWithTopic(topic string, message interface{}) error {
	if bs, err := json.Marshal(message); err != nil {
		return fmt.Errorf("json encode nsq message error: [%s] %s", er.address, err)
	} else {
		return er.PublishWithTopic(topic, bs)
	}
}

// PublishWithTopic 立即推送一个消息到指定topic
func (er *NSQProducer) PublishWithTopic(topic string, message []byte) error {
	if er.producer == nil {
		return fmt.Errorf("nsq producer has stopped or not connected")
	}
	// 不能推送空消息
	if len(message) == 0 {
		return nil
	}
	// 推送消息
	if err := er.producer.Publish(topic, message); err != nil {
		return fmt.Errorf("publish message to nsq error: [%s] %s", er.address, err)
	}
	return nil
}

// Publish 立即推送一个消息到队列中
func (er *NSQProducer) Publish(message []byte) error {
	if er.producer == nil {
		return fmt.Errorf("nsq producer has stopped or not connected")
	}
	// 不能推送空消息
	if len(message) == 0 {
		return nil
	}
	// 推送消息
	if err := er.producer.Publish(er.config.Topic, message); err != nil {
		return fmt.Errorf("publish message to nsq error: [%s] %s", er.address, err)
	}
	return nil
}

// DeferredPublish 立即推送一个消息到队列中
func (er *NSQProducer) DeferredPublish(delay time.Duration, message []byte) error {
	if er.producer == nil {
		return fmt.Errorf("nsq producer has stopped or not connected")
	}
	// 不能推送空消息
	if len(message) == 0 {
		return nil
	}
	// 推送消息
	if err := er.producer.DeferredPublish(er.config.Topic, delay, message); err != nil {
		return fmt.Errorf("publish message to nsq error: [%s] %s", er.address, err)
	}
	return nil
}

// DeferredPublishWithTopic 立即推送一个消息到指定topic
func (er *NSQProducer) DeferredPublishWithTopic(topic string, delay time.Duration, message []byte) error {
	if er.producer == nil {
		return fmt.Errorf("nsq producer has stopped or not connected")
	}
	// 不能推送空消息
	if len(message) == 0 {
		return nil
	}
	// 推送消息
	if err := er.producer.DeferredPublish(topic, delay, message); err != nil {
		return fmt.Errorf("publish message to nsq error: [%s] %s", er.address, err)
	}
	return nil
}

// Close 关闭当前消费者
func (er *NSQProducer) Close() {
	er.mu.Lock()
	if er.producer != nil {
		er.producer.Stop()
		er.producer = nil
		er.address = ""
	}
	er.mu.Unlock()
}
