package fpl

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestFPLRequest(t *testing.T) {
	req, err := http.NewRequest("GET", (&url.URL{
		Scheme: "https",
		Host:   "www.proxy-list.download",
		Path:   "/api/v1/get",
		RawQuery: (&url.Values{
			"type": []string{"http"},
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
	//buff, err := ioutil.ReadAll(rsp.Body)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//t.Log(strings.Split(string(buff), "\n"))
	r := bufio.NewReader(rsp.Body)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				t.Error(err)
			}
			return
		}
		t.Log(line)
	}
}

func TestFPLFetch(t *testing.T) {
	fpl := New("http", "", "")
	list, err := fpl.fetch(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	for k := range list {
		t.Log(list[k])
	}
}

func TestFPLDiscover(t *testing.T) {
	fpl := New("http", "", "")
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
