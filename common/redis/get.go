package redis

import (
	"fmt"
	"github.com/pkg/errors"
	"time"
)

func (c *ClientRedis) GetStringValue(k string) (string, error) {
	var err error
	var value string
	if c.cluster {
		value, err = c.clusterClient.Get(k).Result()
	} else {
		value, err = c.client.Get(k).Result()
	}
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("fail to get string data for key: %s from redis", k))
	}
	return value, nil
}


func (c *ClientRedis) LPop(k string) (string, error) {
	var err error
	var value string
	if c.cluster {
		value, err = c.clusterClient.LPop(k).Result()
	} else {
		value, err = c.client.LPop(k).Result()
	}
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("fail to lpop data for key: %s from redis", k))
	}
	return value, nil
}

func (c *ClientRedis) BLPop(k string, timeout time.Duration) (string, error) {
	var err error
	var value []string
	if c.cluster {
		value, err = c.clusterClient.BLPop(timeout, k).Result()
	} else {
		value, err = c.client.BLPop(timeout, k).Result()
	}
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("fail to blpop data for key: %s from redis", k))
	}
	return value[0], nil
}

func (c *ClientRedis) RPop(k string) (string, error) {
	var err error
	var value string
	if c.cluster {
		value, err = c.clusterClient.RPop(k).Result()
	} else {
		value, err = c.client.RPop(k).Result()
	}
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("fail to rpop data for key: %s from redis", k))
	}
	return value, nil
}

func (c *ClientRedis) BRPop(k string, timeout time.Duration) (string, error) {
	var err error
	var value []string
	if c.cluster {
		value, err = c.clusterClient.BRPop(timeout, k).Result()
	} else {
		value, err = c.client.BRPop(timeout, k).Result()
	}
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("fail to brpop data for key: %s from redis", k))
	}
	return value[0], nil
}


