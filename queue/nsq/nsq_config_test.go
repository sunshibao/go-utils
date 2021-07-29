package nsq

import (
	"fmt"
	"testing"
)

func TestNSQConfig_Copy(t *testing.T) {
	config := &NSQConfig{
		Topic:              "foo",
		Channel:            "foo",
		ConsumerCount:      100,
		ConsumerConcurrent: 200,
		Address:            "127.0.0.1:4160",
		Lookup:             true,
	}
	copied := config.Copy()
	if copied == nil {
		t.Fatalf("NSQConfig.Copy() return nil.")
	}
	// 测试复制后的值是否与原值相等
	if fmt.Sprintf("%v", config) != fmt.Sprintf("%v", copied) {
		t.Fatalf("The copied NSQConfig value is not equal to the original value.")
	}
	// 测试修改复制后的副本是否影响原值
	copied.Topic = "bar"
	if config.Topic == copied.Topic {
		t.Fatalf("Modifying the copied NSQConfig value affects the original value.")
	}
}
