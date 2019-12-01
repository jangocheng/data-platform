package redis

import (
	"fmt"
	"github.com/pkg/errors"
	"time"
)

func (c *ClientRedis) Set (k, v string) error {
	var err error
	if c.cluster {
		err = c.clusterClient.Set(k, v, 0).Err()
	} else {
		err = c.client.Set(k, v, 0).Err()
	}
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("fail to set key: %s to redis", k))
	}
	return nil
}

func (c *ClientRedis)  SetKeyExpire(k, v string, ex time.Duration) error {
	var err error
	if c.cluster {
		err = c.clusterClient.Set(k, v, ex).Err()
	} else {
		err = c.client.Set(k, v, ex).Err()
	}
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("fail to set expire key: %s to redis", k))
	}
	return nil
}

func (c *ClientRedis)  LPush(k string, vs...interface{}) error {
	var err error
	if c.cluster {
		err = c.clusterClient.LPush(k, vs).Err()
	} else {
		err = c.client.LPush(k, vs).Err()
	}
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("fail to lpush data to redis key: %s to redis", k))
	}
	return nil
}

func (c *ClientRedis)  Rpush(k string, vs...interface{}) error {
	var err error
	if c.cluster {
		err = c.clusterClient.RPush(k, vs).Err()
	} else {
		err = c.client.RPush(k, vs).Err()
	}
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("fail to rpush data to redis key: %s to redis", k))
	}
	return nil
}







