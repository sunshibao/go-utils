// Author: Qingshan Luo <edoger@qq.com>
package kafka

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"

	"github.com/sunshibao/go-utils/logs"
)

type Consumer struct {
	c     sarama.ConsumerGroup
	group string
}

func (c *Consumer) Consume(topics []string, handler func(string, []byte) error) error {
	return c.c.Consume(context.Background(), topics, newConsumerHandlerWrapper(handler))
}

func (c *Consumer) Close() error {
	return c.c.Close()
}

func newConsumerHandlerWrapper(handler func(string, []byte) error) sarama.ConsumerGroupHandler {
	return &consumerHandlerWrapper{handler: handler}
}

type consumerHandlerWrapper struct {
	handler func(topic string, data []byte) error
}

func (w *consumerHandlerWrapper) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (w *consumerHandlerWrapper) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (w *consumerHandlerWrapper) ConsumeClaim(s sarama.ConsumerGroupSession, g sarama.ConsumerGroupClaim) error {
	for message := range g.Messages() {
		s.MarkMessage(message, "")
		if err := w.call(message.Topic, message.Value); err != nil {
			logs.Errorf("Kafka consume topic %s error: %s, with message %s", message.Topic, err, string(message.Value))
		}
	}
	return nil
}

func (w *consumerHandlerWrapper) call(topic string, data []byte) (err error) {
	defer func() {
		if v := recover(); v != nil {
			switch o := v.(type) {
			case string:
				err = fmt.Errorf("panic: %s", o)
			case error:
				err = fmt.Errorf("panic: %w", o)
			default:
				err = fmt.Errorf("panic: %v", v)
			}
		}
	}()
	err = w.handler(topic, data)
	return
}
