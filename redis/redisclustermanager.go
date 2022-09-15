package redis

import (
	"io"

	red "github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/core/syncx"
)

var clusterManager = syncx.NewResourceManager()

func GetCluster(server, pass string) (*red.ClusterClient, error) {
	return getCluster(server, pass)
}

func getCluster(server, pass string) (*red.ClusterClient, error) {
	val, err := clusterManager.GetResource(server, func() (io.Closer, error) {
		store := red.NewClusterClient(&red.ClusterOptions{
			Addrs:        []string{server},
			Password:     pass,
			MaxRetries:   maxRetries,
			MinIdleConns: idleConns,
		})
		store.WrapProcess(process)

		return store, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*red.ClusterClient), nil
}
