package nsq

import (
	"sync"
)

// ----------------------------------------
//  NSQ 队列客户端管理器
// ----------------------------------------
type NSQManager struct {
	rw      *sync.RWMutex
	clients map[string]*NSQClient
}

// 创建 NSQ 队列客户端管理器
func CreateNSQManager() *NSQManager {
	return &NSQManager{
		rw:      &sync.RWMutex{},
		clients: make(map[string]*NSQClient),
	}
}

// 获取指定名称的 NSQ 客户端实例，如果客户端不存在，返回 nil
func (manager *NSQManager) Get(name string) *NSQClient {
	manager.rw.RLock()
	defer manager.rw.RUnlock()
	if client, exists := manager.clients[name]; exists {
		return client
	}
	return nil
}

// 添加一个 NSQ 客户端，如果已经存在同名的客户端了，将覆盖旧的客户端
func (manager *NSQManager) Set(client *NSQClient) {
	manager.rw.Lock()
	manager.clients[client.Name()] = client
	manager.rw.Unlock()
}
