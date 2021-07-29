// Author: Qingshan Luo <edoger@qq.com>
package kafka

import (
	"fmt"
	"sync"

	"github.com/Shopify/sarama"

	"github.com/sunshibao/go-utils/logs"
)

type Client struct {
	c   sarama.Client
	cmu sync.RWMutex
	pmu sync.RWMutex
	cs  map[string]*Consumer
	ps  map[string]*Producer
}

func MustNew(config *Config) *Client {
	client, err := New(config)
	if err != nil {
		panic(err)
	}
	return client
}

func New(config *Config) (*Client, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	if config.Version != "" {
		v, err := sarama.ParseKafkaVersion(config.Version)
		if err != nil {
			return nil, err
		}
		cfg.Version = v
	}
	client, err := sarama.NewClient(config.Addrs, cfg)
	if err != nil {
		return nil, err
	}
	r := &Client{c: client, cs: make(map[string]*Consumer), ps: make(map[string]*Producer)}
	if len(config.ConsumerGroups) > 0 {
		for _, group := range config.ConsumerGroups {
			if r.cs[group] != nil {
				return nil, fmt.Errorf("duplicate kafka consumer group %s", group)
			}
			consumer, err := r.newConsumer(group)
			if err != nil {
				return nil, err
			}
			r.cs[group] = consumer
		}
	}
	if len(config.ProducerTopics) > 0 {
		for _, topic := range config.ProducerTopics {
			if r.ps[topic] != nil {
				return nil, fmt.Errorf("duplicate kafka producer topic %s", topic)
			}
			producer, err := r.newProducer(topic)
			if err != nil {
				return nil, err
			}
			r.ps[topic] = producer
		}
	}
	return r, nil
}

func (c *Client) Consumer(group string) *Consumer {
	c.cmu.RLock()
	consumer, found := c.cs[group]
	c.cmu.RUnlock()
	if found {
		return consumer
	}
	c.cmu.Lock()
	defer c.cmu.Unlock()
	consumer, found = c.cs[group]
	if found {
		return consumer
	}
	var err error
	consumer, err = c.newConsumer(group)
	if err != nil {
		logs.Errorf("Create kafka consumer group %s error: %s", group, err)
	} else {
		c.cs[group] = consumer
	}
	return consumer
}

func (c *Client) newConsumer(group string) (*Consumer, error) {
	o, err := sarama.NewConsumerGroupFromClient(group, c.c)
	if err != nil {
		return nil, err
	}
	return &Consumer{c: o, group: group}, nil
}

func (c *Client) Producer(topic string) *Producer {
	c.pmu.RLock()
	producer, found := c.ps[topic]
	c.pmu.RUnlock()
	if found {
		return producer
	}
	c.pmu.Lock()
	defer c.pmu.Unlock()
	producer, found = c.ps[topic]
	if found {
		return producer
	}
	var err error
	producer, err = c.newProducer(topic)
	if err != nil {
		logs.Errorf("Create kafka producer topic %s error: %s", topic, err)
	} else {
		c.ps[topic] = producer
	}
	return producer
}

func (c *Client) newProducer(topic string) (*Producer, error) {
	o, err := sarama.NewSyncProducerFromClient(c.c)
	if err != nil {
		return nil, err
	}
	return &Producer{p: o, topic: topic}, nil
}

func (c *Client) Consume(group string, topics []string, handler func(string, []byte) error) error {
	consumer := c.Consumer(group)
	if consumer == nil {
		return fmt.Errorf("can not create consumer with group %s", group)
	}
	return consumer.Consume(topics, handler)
}

func (c *Client) Publish(topic string, data []byte) error {
	producer := c.Producer(topic)
	if producer == nil {
		return fmt.Errorf("can not create producer with topic %s", topic)
	}
	return producer.Publish(data)
}

func (c *Client) PublishFrom(topic string, f func() ([]byte, error)) error {
	producer := c.Producer(topic)
	if producer == nil {
		return fmt.Errorf("can not create producer with topic %s", topic)
	}
	return producer.PublishFrom(f)
}

func (c *Client) PublishString(topic string, data string) error {
	producer := c.Producer(topic)
	if producer == nil {
		return fmt.Errorf("can not create producer with topic %s", topic)
	}
	return producer.PublishString(data)
}

func (c *Client) PublishStringFrom(topic string, f func() (string, error)) error {
	producer := c.Producer(topic)
	if producer == nil {
		return fmt.Errorf("can not create producer with topic %s", topic)
	}
	return producer.PublishStringFrom(f)
}

func (c *Client) PublishJSON(topic string, value interface{}) error {
	producer := c.Producer(topic)
	if producer == nil {
		return fmt.Errorf("can not create producer with topic %s", topic)
	}
	return producer.PublishJSON(value)
}

func (c *Client) PublishJSONFrom(topic string, f func() (interface{}, error)) error {
	producer := c.Producer(topic)
	if producer == nil {
		return fmt.Errorf("can not create producer with topic %s", topic)
	}
	return producer.PublishJSONFrom(f)
}
