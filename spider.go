package proxies

import "context"

type Spider interface {
	Discover(context.Context) <-chan *Transport
}

type Spiders []Spider

func (list Spiders) Discover(ctx context.Context) <-chan *Transport {
	merge := make(chan *Transport)
	for _, c := range list {
		go func(c <-chan *Transport) {
			for v := range c {
				merge <- v
			}
		}(c.Discover(ctx))
	}
	return merge
}
