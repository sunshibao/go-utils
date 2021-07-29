// Author: Sunshibao <664588619@qq.com>
package configsdk

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/common/logger"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type ACMLoaderOptions struct {
	Endpoint  string
	Namespace string
	Group     string
	AccessKey string
	SecretKey string
}

func MustNewACMLoader(opt ACMLoaderOptions) Loader {
	loader, err := NewACMLoader(opt)
	if err != nil {
		panic(err)
	}
	return loader
}

func NewACMLoader(opt ACMLoaderOptions) (Loader, error) {
	cfg := constant.ClientConfig{
		Endpoint:            opt.Endpoint + ":8080",
		NamespaceId:         opt.Namespace,
		AccessKey:           opt.AccessKey,
		SecretKey:           opt.SecretKey,
		TimeoutMs:           3000,
		BeatInterval:        5000,
		NotLoadCacheAtStart: true,
	}

	client, err := clients.CreateConfigClient(map[string]interface{}{
		constant.KEY_CLIENT_CONFIG: cfg,
	})
	if err != nil {
		return nil, err
	}
	logger.SetLogger(GetLogger())
	return &acmLoader{opt.Group, client}, nil
}

type acmLoader struct {
	group  string
	client config_client.IConfigClient
}

func (o *acmLoader) Load(target string) ([]byte, error) {
	content, err := o.client.GetConfig(vo.ConfigParam{DataId: target, Group: o.group})
	if err != nil {
		// 当 ACM 获取不到配置时返回一个内容为 "config not found" 的错误
		if err.Error() == "config not found" {
			return nil, ErrNotFound
		}
		return nil, err
	}
	// 当 ACM 获取不到配置时会将空字符串写入内部缓存，第二次获取时如果配置还不存在，
	// 则不会返回错误，直接返回空字符串，这里需要转换为 ErrNotFound
	if content == "" {
		return nil, ErrNotFound
	}
	return []byte(content), nil
}
