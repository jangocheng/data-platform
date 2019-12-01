package redis

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)


type ClientRedis struct {
	cluster     	bool
	client 			*redis.Client
	clusterClient 	*redis.ClusterClient
}

type ConnectConfig struct {
	IpPortList      []string
	Password 		string
	PoolSize        int
	MinIdexCon      int
	DB              int
}


func NewRedisClient(conf ConnectConfig) (*ClientRedis, error) {
	if len(conf.IpPortList) == 0 {
		return nil, errors.New("fail to create redis client, the ip port list can not be none")
	}
	if len(conf.IpPortList) != 1 {
		client := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:  conf.IpPortList,
			Password:  conf.Password,
			PoolSize:   conf.PoolSize,
			MinIdleConns: conf.MinIdexCon,
		})
		_, err := client.Do("select", conf.DB).Result()
		if err != nil {
			return nil, errors.Wrap(err, "fail to create redis cluster client")
		}
		_, err = client.Ping().Result()
		if err != nil {
			return nil, errors.Wrap(err, "fail to create redis client")
		}
		return &ClientRedis{cluster: true, clusterClient:client, client: nil}, nil
	} else {
		client := redis.NewClient(&redis.Options{
			Addr: conf.IpPortList[0],
			Password: conf.Password,
			PoolSize:  conf.PoolSize,
			MinIdleConns:  conf.MinIdexCon,
			DB:   conf.DB,
		})
		_, err := client.Ping().Result()
		if err != nil {
			return nil, errors.Wrap(err, "fail to create redis client")
		}
		return &ClientRedis{cluster:false, client:client, clusterClient:nil}, nil
	}
}

func (c *ClientRedis) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	if c.client != nil {
		return c.clusterClient.Close()
	}
	return nil
}
