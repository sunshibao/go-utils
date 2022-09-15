// Author: Qingshan Luo <edoger@qq.com>
package goqueue

import (
	"github.com/zeromicro/go-queue/dq"
)

func NewProducer(conf dq.DqConf) dq.Producer {
	return dq.NewProducer(conf.Beanstalks)
}
