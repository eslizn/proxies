package spider

import (
	"context"
	"github.com/eslizn/proxies"
	"testing"
	"time"
)

func TestSpider_Discover(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	driver := New(nil)
	for {
		select {
		case <-ctx.Done():
			t.Logf("context done\n")
			return
		case val := <-driver.Discover(ctx):
			t.Logf("discover: %s\n", val.String())
		}
	}
}

func TestSpider_Parse(t *testing.T) {
	s := New(nil)
	reader, err := s.request(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	defer reader.Close()
	err = s.parse(reader, func(tran *proxies.Transport) error {
		t.Logf("parse: %s\n", tran.String())
		return nil
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestTransports(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	trans := proxies.NewTransports(ctx, New(nil))
	defer trans.Close()
	for i := 0; i < 5; i++ {
		tran, err := trans.Assign()
		if err != nil {
			t.Error(err)
			return
		}
		t.Logf("get transport: %s\n", tran.String())
	}
}
