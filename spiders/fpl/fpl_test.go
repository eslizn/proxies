package fpl

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestFPLRequest(t *testing.T) {
	req, err := http.NewRequest("GET", (&url.URL{
		Scheme: "https",
		Host:   "www.free-proxy-list.com",
		Path:   "/",
		RawQuery: (&url.Values{
			"page": []string{"1"},
		}).Encode(),
	}).String(), nil)
	if err != nil {
		t.Error(err)
		return
	}
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
		return
	}
	defer rsp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(rsp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(doc)
}

func TestFPLFetch(t *testing.T) {
	fpl := New()
	list, err := fpl.fetchWithPage(context.Background(), 1)
	if err != nil {
		t.Error(err)
		return
	}
	for k := range list {
		t.Log(list[k])
	}
}

func TestFPLDiscover(t *testing.T) {
	fpl := New()
	channel := make(chan string, 0)
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	go fpl.Discover(ctx, channel)
	for {
		select {
		case <-ctx.Done():
			t.Log(ctx.Err())
			return
		case str := <-channel:
			t.Log(str)
		}
	}
}
