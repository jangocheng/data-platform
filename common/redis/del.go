package redis

import (
	"fmt"
	"github.com/pkg/errors"
)

func (c *ClientRedis) DelKey(k string) error {
	var err error
	if c.cluster {
		_, err = c.clusterClient.Del(k).Result()
	} else {
		_, err = c.client.Del(k).Result()
	}
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("fail to delete key: %s from redis", k))
	}
	return nil
}
