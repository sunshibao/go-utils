package redis

import "errors"

var (
	ErrEmptyHost = errors.New("empty redis host")
	ErrEmptyType = errors.New("empty redis type")
	ErrEmptyKey  = errors.New("empty redis key")
)

type (
	JingRedisConf struct {
		Host string
		Type string `json:",default=node,options=node|cluster"`
		Pass string `json:",optional"`
	}

	JingRedisKeyConf struct {
		JingRedisConf
		Key string `json:",optional"`
	}
)

func (rc JingRedisConf) NewJingRedis() *JingRedis {
	return NewJingRedis(rc.Host, rc.Type, rc.Pass)
}

func (rc JingRedisConf) Validate() error {
	if len(rc.Host) == 0 {
		return ErrEmptyHost
	}

	if len(rc.Type) == 0 {
		return ErrEmptyType
	}

	return nil
}

func (rkc JingRedisKeyConf) Validate() error {
	if err := rkc.JingRedisConf.Validate(); err != nil {
		return err
	}

	if len(rkc.Key) == 0 {
		return ErrEmptyKey
	}

	return nil
}
