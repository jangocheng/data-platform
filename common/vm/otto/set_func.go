package otto

import "github.com/pkg/errors"

func (c *ClientOtto) SetFunc(name string, function interface{}) error {
	c.lock.Lock()
	err := c.vm.Set(name, function)
	if err != nil {
		err = errors.Wrap(err, "fail to set function to otto client")
	}
	c.lock.Unlock()
	return err
}
