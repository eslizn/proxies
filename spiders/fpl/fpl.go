package fpl

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/eslizn/proxies"
	"net/http"
	"net/url"
	"time"
)

//free proxy list: https://www.free-proxy-list.com/

type FreeProxyList url.URL

func (fpl *FreeProxyList) Discover(ctx context.Context, channel chan string) error {
	page := 1
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.NewTicker(3 * time.Second).C:
			list, err := fpl.fetchWithPage(ctx, page)
			if err != nil {
				continue
				//return err
			}
			page++
			for k := range list {
				go func(addr string) {
					if proxies.Available(ctx, addr) == nil {
						channel <- list[k]
					}
				}(list[k])
			}
			if len(list) < 10 {
				page = 1
				time.Sleep(3 * time.Minute)
			}
		}
	}
}

func (fpl *FreeProxyList) fetchWithPage(ctx context.Context, page int) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", (&url.URL{
		Scheme:   fpl.Scheme,
		Host:     fpl.Host,
		Path:     fpl.Path,
		RawQuery: (&url.Values{"page": []string{fmt.Sprint(page)}}).Encode(),
	}).String(), nil)
	if err != nil {
		return nil, err
	}
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	dom, err := goquery.NewDocumentFromReader(rsp.Body)
	if err != nil {
		return nil, err
	}
	return dom.Find(".proxy-list > tbody > tr").Map(func(index int, selection *goquery.Selection) string {
		return selection.Find("td").Eq(8).Text() + "://" + selection.Find("td").Eq(0).Find("a").AttrOr("alt", "")
	}), nil
}

func New() *FreeProxyList {
	return (*FreeProxyList)(&url.URL{
		Scheme: "https",
		Host:   "www.free-proxy-list.com",
		Path:   "/",
	})
}
