package spider

import (
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja"
	"github.com/eslizn/proxies"
	"html"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Option struct {
	UserAgent   string
	ContentType string
	Type        string
	Country     string
	ChanSize    int
}

type Spider struct {
	*Option
	proxies chan *proxies.Transport
}

func (s *Spider) Discover(ctx context.Context) <-chan *proxies.Transport {
	go func() {
		body, err := s.request(ctx)
		if err != nil {
			return
		}
		defer body.Close()
		err = s.parse(body, func(tran *proxies.Transport) error {
			s.proxies <- tran
			return nil
		})
		if err != nil {
			return
		}
	}()
	return s.proxies
}

func (s *Spider) request(ctx context.Context) (io.ReadCloser, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://spys.one/free-proxy-list/%s/", s.Country), bytes.NewBuffer(nil))
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", s.UserAgent)
	request.Header.Set("Content-Type", s.ContentType)
	client := &http.Client{}
	if len(os.Getenv("HTTP_PROXY")) > 0 {
		client.Transport = &http.Transport{Proxy: func(r *http.Request) (*url.URL, error) {
			return url.Parse(os.Getenv("HTTP_PROXY"))
		}}
	}
	result, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}

func (s *Spider) parse(r io.Reader, cb func(*proxies.Transport) error) error {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return err
	}
	vm := goja.New()
	htmlStr, err := doc.Find("table:nth-of-type(1)").Next().Html()
	if err != nil {
		return err
	}
	_, err = vm.RunString(html.UnescapeString(htmlStr))
	if err != nil {
		return err
	}
	section := doc.Find("font.spy14")
	section.Each(func(i int, s *goquery.Selection) {
		script := s.ChildrenFiltered("script")
		if script.Length() != 1 {
			return
		}
		script.Remove()
		portScript := script.Text()
		portScript = portScript[strings.Index(portScript, "+")+1 : len(portScript)-1]
		portScript = strings.ReplaceAll(portScript, ")", ").toString()")
		port, err := vm.RunString(portScript)
		if err != nil {
			return
		}
		err = cb(proxies.NewTransport(url.URL{
			Scheme: strings.ToLower(strings.Split(s.Parent().Next().Text(), " ")[0]),
			Host:   fmt.Sprintf("%s:%s", s.Text(), port.String()),
		}))
		if err != nil {
			//
			return
		}
	})
	return nil
}

func New(opt *Option) *Spider {
	if opt == nil {
		opt = &Option{}
	}
	if len(opt.UserAgent) < 1 {
		opt.UserAgent = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36 MicroMessenger/6.5.2.501 NetType/WIFI WindowsWechat QBCore/3.43.1021.400 QQBrowser/9.0.2524.400"
	}
	if len(opt.ContentType) < 1 {
		opt.ContentType = "text/html; charset=utf-8"
	}
	if opt.ChanSize == 0 {
		opt.ChanSize = 32
	}

	return &Spider{
		Option:  opt,
		proxies: make(chan *proxies.Transport, opt.ChanSize),
	}
}
