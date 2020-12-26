package proxies

import (
	"bytes"
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func NewTransport(url url.URL) *Transport {
	return &Transport{
		URL: url,
		Transport: http.Transport{
			Proxy:           http.ProxyURL(&url),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

type Transport struct {
	sync.RWMutex
	url.URL
	http.Transport
	Stats
}

func (t *Transport) Available(ctx context.Context, url string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, bytes.NewBuffer(nil))
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "curl/7.29.0")
	client := &http.Client{Transport: t}
	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	if rsp.StatusCode != http.StatusOK {
		return ErrProbeFail
	}
	t.Lock()
	defer t.Unlock()
	now := time.Now()
	t.Verified = &now
	return nil
}

type Transports struct {
	ctx        context.Context
	cancel     context.CancelFunc
	expired    time.Duration
	spiders    Spiders
	transports chan *Transport
	urls       []string
}

func (t *Transports) process() {
	for {
		select {
		case <-t.ctx.Done():
			return
		case tran := <-t.spiders.Discover(t.ctx):
			go func() {
				if tran.Available(t.ctx, t.urls[int(time.Now().Unix())%len(t.urls)]) == nil {
					t.transports <- tran
				}
			}()
		}
	}
}

func (t *Transports) Assign() (*Transport, error) {
	select {
	case <-t.ctx.Done():
		return nil, context.Canceled
	case tran := <-t.transports:
		return tran, nil
	}
}

func (t *Transports) Close() {
	if t.cancel == nil {
		return
	}
	t.cancel()
	close(t.transports)
	t.cancel = nil
}

func NewTransports(ctx context.Context, spiders ...Spider) *Transports {
	ctx, cancel := context.WithCancel(ctx)
	trans := &Transports{
		ctx:        ctx,
		cancel:     cancel,
		expired:    5 * time.Minute,
		spiders:    spiders,
		transports: make(chan *Transport, 32),
		urls: []string{
			"https://ifconfig.co/ip",
			"http://myip.ipip.net/",
			"https://ifconfig.me/",
			"https://ipinfo.io",
		},
	}
	go trans.process()
	return trans
}
