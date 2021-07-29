// Author: Qingshan Luo <edoger@qq.com>
package goqueue

import (
    "github.com/tal-tech/go-queue/dq"
)

func NewConsumer(conf dq.DqConf) dq.Consumer {
    return dq.NewConsumer(conf)
}