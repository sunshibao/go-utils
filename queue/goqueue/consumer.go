// Author: Qingshan Luo <edoger@qq.com>
package goqueue

import (
	"github.com/zeromicro/go-queue/dq"
)

func NewConsumer(conf dq.DqConf) dq.Consumer {
	return dq.NewConsumer(conf)
}
