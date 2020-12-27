package proxies

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"sync"
)

var (
	channel = make(chan string, 32)
	cache   = sync.Map{}
)

// Get get trans
func Get() (*http.Transport, error) {
	return GetWithContext(context.Background())
}

// GetWithContext get trans with context
func GetWithContext(ctx context.Context) (*http.Transport, error) {
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	case val := <-channel:
		trans, has := cache.Load(val)
		if !has {
			trans = &http.Transport{
				Proxy: func(*http.Request) (*url.URL, error) {
					return url.Parse(val)
				},
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			cache.Store(val, trans)
		}
		return trans.(*http.Transport), nil
	}
}

// Put put trans
func Put(trans *http.Transport) error {
	return PutWithContext(context.Background(), trans)
}

// PutWithContext put trans with context
func PutWithContext(ctx context.Context, trans *http.Transport) error {
	val, err := trans.Proxy(nil)
	if err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return context.Canceled
	case channel <- val.String():
		return nil
	}
}
