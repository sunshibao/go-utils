// Author: Qingshan Luo <edoger@qq.com>
package kafka

import (
	"encoding/json"

	"github.com/Shopify/sarama"
)

type Producer struct {
	p     sarama.SyncProducer
	topic string
}

func (p *Producer) publish(message *sarama.ProducerMessage) error {
	if _, _, err := p.p.SendMessage(message); err != nil {
		return err
	}
	return nil
}

func (p *Producer) Publish(data []byte) error {
	message := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(data),
	}
	return p.publish(message)
}

func (p *Producer) PublishFrom(f func() ([]byte, error)) error {
	data, err := f()
	if err != nil {
		return err
	}
	return p.Publish(data)
}

func (p *Producer) PublishString(data string) error {
	message := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(data),
	}
	return p.publish(message)
}

func (p *Producer) PublishStringFrom(f func() (string, error)) error {
	data, err := f()
	if err != nil {
		return err
	}
	return p.PublishString(data)
}

func (p *Producer) PublishJSON(value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	message := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(data),
	}
	return p.publish(message)
}

func (p *Producer) PublishJSONFrom(f func() (interface{}, error)) error {
	data, err := f()
	if err != nil {
		return err
	}
	return p.PublishJSON(data)
}

func (p *Producer) Close() error {
	return p.p.Close()
}
