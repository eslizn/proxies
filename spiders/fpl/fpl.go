package fpl

import (
	"bufio"
	"context"
	"fmt"
	"github.com/eslizn/proxies"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//free proxy list: https://www.proxy-list.download/api/v1

type FreeProxyList url.URL

func (fpl *FreeProxyList) Discover(ctx context.Context, channel chan string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.NewTicker(3 * time.Second).C:
			list, err := fpl.fetch(ctx)
			if err != nil {
				continue
				//return err
			}
			for k := range list {
				go func(addr string) {
					if proxies.Available(ctx, addr) == nil {
						channel <- list[k]
					}
				}(list[k])
			}
		}
	}
}

func (fpl *FreeProxyList) fetch(ctx context.Context) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", (*url.URL)(fpl).String(), nil)
	if err != nil {
		return nil, err
	}
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	reader := bufio.NewReader(rsp.Body)
	var (
		line  = ""
		lines = make([]string, 0, 128)
	)
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		lines = append(lines, fmt.Sprintf("%s://%s", strings.ToLower(fpl.Scheme), line))
	}
	return lines, err
}

func New(typ, country, anon string) *FreeProxyList {
	query := &url.Values{
		"type": []string{typ},
	}
	if len(country) > 0 {
		query.Set("country", country)
	}
	if len(anon) > 0 {
		query.Set("anon", anon)
	}
	return (*FreeProxyList)(&url.URL{
		Scheme:   typ,
		Host:     "www.proxy-list.download",
		Path:     "/api/v1/get",
		RawQuery: query.Encode(),
	})
}
