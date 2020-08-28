package spider

import (
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/lycblank/gonews/pkg/registry"
	"github.com/parnurzeal/gorequest"
	"math/rand"
	"time"
)

type NewsSpiderOptions struct {
	Addr string
	Path string

	Interval time.Duration
	RandRange time.Duration
}

type NewsSpiderOption func(opts *NewsSpiderOptions)

func WithAddrAndPath(addr string, path string) NewsSpiderOption {
	return func(opts *NewsSpiderOptions) {
		opts.Addr = addr
		opts.Path = path
	}
}

func WithSpiderInterval(interval time.Duration, randRange time.Duration) NewsSpiderOption {
	return func(opts *NewsSpiderOptions) {
		opts.Interval = interval
		opts.RandRange = randRange
	}
}

type NewsListItem struct {
	Title string `json:"title"`
	Href string  `json:"href"`
}

type NewsSpider struct {
	Addr string
	Path string

	Interval time.Duration
	RandRange time.Duration

	store registry.Storer
}

func NewNewsSpider(store registry.Storer, options ...NewsSpiderOption) *NewsSpider {
	opts := &NewsSpiderOptions{
		Addr: "https://gocn.vip",
		Path: "/topics/node18",
		Interval:time.Hour,
		RandRange:10*time.Minute,
	}
	for _, option := range options {
		option(opts)
	}

	sp := &NewsSpider{
		Addr:opts.Addr,
		Path:opts.Path,
		Interval:opts.Interval,
		RandRange:opts.RandRange,
		store:store,
	}

	return sp
}

func (ns *NewsSpider) Run() {
	if err := ns.process(1); err != nil {
		fmt.Println(err)
	}
	duration := ns.getSleepDuration()
	tm := time.NewTimer(duration)
	for {
		select {
		case <-tm.C:
			if err := ns.process(1); err != nil {
				fmt.Println(err)
			}
			duration = ns.getSleepDuration()
			tm.Reset(duration)
		}
	}
}

func (ns *NewsSpider) getSleepDuration() time.Duration {
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return ns.Interval + time.Duration(rd.Int63n(int64(ns.RandRange)))
}

func (ns *NewsSpider) process(page int) error {
	ctx := context.Background()
	newsList, err := ns.GetNewsList(ctx, page)
	if err != nil {
		return err
	}
	for i,cnt:=0,len(newsList); i< cnt;i++{
		items, _ := ns.GetNewsContent(ctx, newsList[i].Href)
		if err := ns.store.Store(ctx, items...); err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func (ns *NewsSpider) GetNewsList(ctx context.Context, page int) ([]NewsListItem, error) {
	listUrl := fmt.Sprintf("%s%s?page=%d", ns.Addr, ns.Path, page)
	body, err := ns.httpGet(ctx, listUrl)
	if err != nil {
		return nil, err
	}
	return ns.parseNewsList(ctx, body)
}

func (ns *NewsSpider) parseNewsList(ctx context.Context, datas []byte) ([]NewsListItem, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(datas))
	if err != nil {
		return nil, err
	}
	items := make([]NewsListItem, 0, 16)
	doc.Find("#main .row .infos .title a").Each(func(_ int, s *goquery.Selection){
		if title, ok := s.Attr("title"); ok {
			href, _ := s.Attr("href")
			items = append(items, NewsListItem{
				Title:title,
				Href:ns.Addr+href,
			})
		}
	})
	return items, nil
}

func (ns *NewsSpider) GetNewsContent(ctx context.Context, href string) ([]registry.NewsItem, error) {
	body, err := ns.httpGet(ctx, href)
	if err != nil {
		return nil, err
	}
	return ns.parseNewsItemList(ctx, body, href)
}

func (ns *NewsSpider) parseNewsItemList(ctx context.Context, datas []byte, spiderHref string) ([]registry.NewsItem, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(datas))
	if err != nil {
		return nil, err
	}
	items := make([]registry.NewsItem, 0, 16)
	doc.Find("#main .row ol li").Each(func(_ int, s *goquery.Selection){
		title := s.Text()
		var href string
		if s = s.Find("a"); s!=nil{
			href = s.AttrOr("href", "")
		}

		if href != "" {
			items = append(items, registry.NewsItem{
				Title:      title,
				Href:       href,
				SpiderHref: spiderHref,
				CreateTime: time.Now().Unix(),
			})
		}
	})
	return items, nil
}

func (ns *NewsSpider) httpGet(ctx context.Context, addr string) ([]byte, error) {
	time.Sleep(time.Second + time.Duration(rand.Int63n(100))*time.Millisecond)
	resp, body, errs := gorequest.New().Get(addr).EndBytes()
	if len(errs) > 0 {
		return nil, errs[0]
	}
	defer resp.Body.Close()
	return body, nil
}