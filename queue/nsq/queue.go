package nsq

// ----------------------------------------
//  队列管理器
// ----------------------------------------
type Queue struct {
	nm *NSQManager
}

// 创建队列管理器
func CreateQueue() *Queue {
	return &Queue{
		nm: CreateNSQManager(),
	}
}

// 注册一个 NSQ 客户端
func (queue *Queue) RegisterNSQClient(name string, config *NSQConfig) error {
	queue.nm.Set(CreateNSQClient(name, config))
	return nil
}

// 注册一组 NSQ 客户端
func (queue *Queue) RegisterNSQClients(configs map[string]*NSQConfig) error {
	if len(configs) > 0 {
		for name, config := range configs {
			if err := queue.RegisterNSQClient(name, config); err != nil {
				return err
			}
		}
	}
	return nil
}

// 获取一个指定名称的 NSQ 客户端实例
func (queue *Queue) NSQ(name string) *NSQClient {
	return queue.nm.Get(name)
}
