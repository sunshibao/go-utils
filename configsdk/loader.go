// Author: Sunshibao <664588619@qq.com>
package configsdk

import (
	"errors"
)

// ErrNotFound 表示配置项找不到，当加载流程遇到这个错误时，将加载权交给下一个加载器
var ErrNotFound = errors.New("config: item not found")

type Loader interface {
	Load(target string) ([]byte, error)
}

type LoaderFunc func(target string) ([]byte, error)

func (f LoaderFunc) Load(target string) ([]byte, error) {
	return f(target)
}
