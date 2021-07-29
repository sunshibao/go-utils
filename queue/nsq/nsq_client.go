package nsq

// ----------------------------------------
//  NSQ 队列客户端
// ----------------------------------------
type NSQClient struct {
	name     string
	config   *NSQConfig
	consumer *NSQConsumer
	producer *NSQProducer
}

// 创建一个 NSQ 队列客户端
func CreateNSQClient(name string, config *NSQConfig) *NSQClient {
	return &NSQClient{
		name:     name,
		config:   config.Copy(),
		consumer: CreateNSQConsumer(config),
		producer: CreateNSQProducer(config),
	}
}

// 获取 NSQ 客户端名称
func (client *NSQClient) Name() string {
	return client.name
}

// 获取 NSQ 配置项
func (client *NSQClient) Config() *NSQConfig {
	return client.config.Copy()
}

// 获取 NSQ 消费者
func (client *NSQClient) Consumer() *NSQConsumer {
	return client.consumer
}

// 获取 NSQ 生产者
func (client *NSQClient) Producer() *NSQProducer {
	return client.producer
}
