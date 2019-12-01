package utils

import (
	"os"
	"os/signal"
	"syscall"
)

type Killer struct {
	killNow bool
	signals []os.Signal
}

func (k *Killer) init() {
	c := make(chan os.Signal)
	signal.Notify(c, k.signals...)
	go func() {
		for s := range c {
			switch s {
			case k.signals[0], k.signals[1]:
				k.killNow = true
				close(c)
			default:
				break
			}
		}
	}()
}

func NewKiller() *Killer {
	killer := Killer{}
	killer.killNow = false
	killer.signals = []os.Signal{
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGKILL,
	}
	killer.init()
	return &killer
}

func (k *Killer) KillNow() bool {
	return k.killNow
}
