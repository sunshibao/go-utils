// Author: Qingshan Luo <edoger@qq.com>
package goqueue

import (
    "github.com/tal-tech/go-queue/dq"
)

func NewProducer(conf dq.DqConf) dq.Producer {
    return dq.NewProducer(conf.Beanstalks)
}
