package nsq

import (
	"fmt"
	"testing"
)

// 测试 NSQClient 用配置
var configForNSQClientTests = &NSQConfig{
	Topic:              "foo",
	Channel:            "foo",
	ConsumerCount:      100,
	ConsumerConcurrent: 200,
	Address:            "127.0.0.1:4160",
	Lookup:             true,
}

func TestCreateNSQClient(t *testing.T) {
	client := CreateNSQClient("foo", configForNSQClientTests)
	if client == nil {
		t.Fatalf("CreateNSQClient() return nil.")
	}
}

func TestNSQClient_Name(t *testing.T) {
	client := CreateNSQClient("foo", configForNSQClientTests)
	if client.Name() != "foo" {
		t.Fatalf("NSQClient.Name() return value does not match the set value.")
	}
}

func TestNSQClient_Config(t *testing.T) {
	client := CreateNSQClient("foo", configForNSQClientTests)
	if client.Config() == nil {
		t.Fatalf("NSQClient.Config() return nil.")
	}
	if fmt.Sprintf("%v", client.Config()) != fmt.Sprintf("%v", configForNSQClientTests) {
		t.Fatalf("NSQClient.Config() return value does not match the set value.")
	}
}

func TestNSQClient_Consumer(t *testing.T) {
	client := CreateNSQClient("foo", configForNSQClientTests)
	if client.Consumer() == nil {
		t.Fatalf("NSQClient.Consumer() return nil.")
	}
}

func TestNSQClient_Producer(t *testing.T) {
	client := CreateNSQClient("foo", configForNSQClientTests)
	if client.Producer() == nil {
		t.Fatalf("NSQClient.Producer() return nil.")
	}
}
