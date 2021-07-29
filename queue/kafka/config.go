// Author: Qingshan Luo <edoger@qq.com>
package kafka

type Config struct {
	Addrs          []string
	ConsumerGroups []string
	ProducerTopics []string
	Version        string
}
