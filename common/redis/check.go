package redis

import (
	"fmt"
	"github.com/pkg/errors"
)

func (c *ClientRedis) CheckKey(k string) (bool, error) {
	var exist int64
	var err error
	if c.cluster {
		exist, err = c.clusterClient.Exists(k).Result()
	} else {
		exist, err = c.client.Exists(k).Result()
	}
	if err != nil {
		return false, errors.Wrap(err, fmt.Sprintf("fail to check key: %s from redis", k))
	}
	if exist != 0 {
		return true, nil
	}
	return false, nil
}
