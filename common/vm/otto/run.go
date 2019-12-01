package otto

import (
	"github.com/pkg/errors"
	"platform/common/vm"
)

func (c *ClientOtto) Run(js string, args ...interface{}) (vm.Value, error) {
	js = "(function(){"  + js +  "})();"
	c.lock.Lock()
	var err error
	if len(args) != 0 {
		err = c.vm.Set("input", args[0])
	}
	if err != nil {
		c.lock.Unlock()
		return nil, errors.Wrap(err, "fail to set args to js code vm client")
	}
	defer func() {
		if err := recover(); err != nil {
			c.lock.Unlock()
		}
	}()
	value, err := c.vm.Run(js)
	c.lock.Unlock()
	if err != nil {
		return nil, errors.Wrap(err, "fail to run js code in vm client")
	}
	return &ValueOtto{value:value}, nil
}
