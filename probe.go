package proxies

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	ErrProbeFail = errors.New("probe fail")
)

var (
	probes       []string
	probesLocker = &sync.RWMutex{}
)

// Available check proxy is available
func Available(ctx context.Context, addr string) error {
	probesLocker.RLock()
	req, err := http.NewRequestWithContext(ctx, "GET", probes[int(time.Now().Unix())%len(probes)], bytes.NewBuffer(nil))
	probesLocker.RUnlock()
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "curl/7.29.0")
	client := &http.Client{Transport: &http.Transport{
		Proxy: func(*http.Request) (*url.URL, error) {
			return url.Parse(addr)
		},
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}
	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	if rsp.StatusCode != http.StatusOK {
		return ErrProbeFail
	}
	return nil
}

// RegisterProbes register probe server list
func RegisterProbes(list ...string) {
	probesLocker.Lock()
	probes = list
	probesLocker.Unlock()
}

func init() {
	RegisterProbes([]string{
		"https://ifconfig.co/ip",
		"http://myip.ipip.net/",
		"https://ifconfig.me/",
		"https://ipinfo.io",
	}...)
}
