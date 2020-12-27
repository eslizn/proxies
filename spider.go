package proxies

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrSpiderAlreadyRegister = errors.New("spider already register")
	ErrSpiderNotRegister     = errors.New("spider not register")
)

var (
	spiders       = sync.Map{}
	spidersCancel = sync.Map{}
)

type Spider interface {
	Discover(context.Context, chan string) error
}

// Register register and start spider
func Register(name string, spider Spider) error {
	_, ok := spiders.LoadOrStore(name, spider)
	if ok {
		return ErrSpiderAlreadyRegister
	}
	ctx, cancel := context.WithCancel(context.Background())
	go spider.Discover(ctx, channel)
	spidersCancel.Store(name, cancel)
	return nil
}

// Deregister deregister and stop spider
func Deregister(name string) error {
	val, ok := spidersCancel.Load(name)
	if !ok {
		return ErrSpiderNotRegister
	}
	spiders.Delete(name)
	spidersCancel.Delete(name)
	cancel, ok := val.(context.CancelFunc)
	if ok {
		cancel()
	}
	return nil
}
