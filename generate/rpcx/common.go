package rpcx

import "github.com/sunshibao/go-utils/generate"

func init() {
	generate.RegisterLayouter("rpcxclient", &layclient{})
	generate.RegisterLayouter("rpcx", &lay{})
}
