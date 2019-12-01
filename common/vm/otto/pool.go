package otto

import (
	"github.com/pkg/errors"
	"platform/common/vm"
	"time"
)

type ClientOttoPool struct {
	vmFree    	chan vm.ClientVM
	vmList      []vm.ClientVM
	killchan    chan bool
	size        int

}

func NewClientOttoPool(size int) (vm.ClientVM, error) {
	clientPool := &ClientOttoPool{vmFree: make(chan vm.ClientVM, size), killchan: make(chan bool, size), size: size}
	for i:= 0; i<size; i++ {
		vmClient, err := NewClientOtto()
		if err != nil {
			return nil, errors.Wrap(err, "fail to create otto client pool")
		}
		clientPool.vmFree <- vmClient
		clientPool.vmList = append(clientPool.vmList, vmClient)
	}
	return clientPool, nil
}

func (c *ClientOttoPool) Close() {
	defer close(c.vmFree)
	defer close(c.killchan)
	for i:= 0; i<c.size; i++ {
		c.killchan <- true
	}
}


func (c *ClientOttoPool) SetFunc(name string, function interface{}) error {
	for _, vmClient := range c.vmList {
		err := vmClient.SetFunc(name, function)
		if err != nil {
			return err
		}
	}
	return nil
}


func (c *ClientOttoPool) Run(js string, args ...interface{}) (vm.Value, error) {
	timeout := time.After(time.Second * 5)
	select {
	case client := <- c.vmFree:
		defer func() {
			if err := recover(); err != nil {
				c.vmFree <- client
			}
		}()
		value, err := client.Run(js, args...)
		if err != nil {
			c.vmFree <- client
			return nil, errors.Wrap(err, "fail to run js code in vm client")
		}
		c.vmFree <- client
		return value, nil
	case <- timeout:
		return nil, errors.New("get js code vm timeout")
	case <- c.killchan:
		return nil, errors.New("receive kill signal")
	}
}
