package registry

import (
	"context"
)

type NewsItem struct {
	Title string      `json:"title"`
	Href string       `json:"href"`
	SpiderHref string `json:"spider_href"`
	CreateTime int64  `json:"create_time"`
}

type Storer interface {
	Store(ctx context.Context, news ...NewsItem) error
}

type Searcher interface {
	Search(ctx context.Context, keywords ...string) ([]NewsItem, error)
}

type Pager interface {
	Page(ctx context.Context, offset int32, limit int32) (news []NewsItem, hasMore bool, err error)
}

