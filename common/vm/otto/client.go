package otto

import (
	"github.com/robertkrimen/otto"
	vmInterface "platform/common/vm"
	"sync"
)

type ClientOtto struct {
	vm    	*otto.Otto
	lock    sync.Mutex
}

func NewClientOtto() (vmInterface.ClientVM, error) {
	vm := otto.New()
	client := &ClientOtto{vm: vm}
	return client, nil
}


func (c *ClientOtto) Close() {
	c.vm = nil
}