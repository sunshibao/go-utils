package main

import (
	"fmt"
	"github.com/beanstalkd/go-beanstalk"
	"time"
)

func main() {
	c, err := beanstalk.Dial("tcp", "127.0.0.1:11300")
	if err != nil {
		return
	}
	id, err := c.Put([]byte("hello"), 1, 0, 120*time.Second)
	if err != nil {
		return
	}
	fmt.Println(id)
}
