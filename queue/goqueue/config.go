// Author: Qingshan Luo <edoger@qq.com>
package goqueue

import (
	"github.com/tal-tech/go-zero/core/stores/redis"
)

/**
GoQueue:
  - Name: default
    Conf:
      Beanstalks:
        - Endpoint: 127.0.0.1:11300
          Tube: tube
        - Endpoint: 127.0.0.1:11301
          Tube: tube
      Redis:
        Host: 127.0.0.1:6379
        Type: node
        Pass: ''
*/
type Configs []Config

/* type Config struct {
	Name string
	Conf dq.DqConf
} */

type Config struct {
	Endpoints []string
	Redis     redis.RedisConf
	Topics    []Topic
}

type Topic struct {
	Name string
	Tube string
}
