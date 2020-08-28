package registry

import (
	"context"
	"encoding/json"
	"fmt"
)

var _ Registry = &stdoutRegistry{}

type stdoutRegistry struct {

}

func NewStdoutRegistry() Registry {
	return &stdoutRegistry{}
}

func (sr *stdoutRegistry) Store(ctx context.Context, news ...NewsItem) error {
	datas, _ := json.Marshal(news)
	fmt.Println(string(datas))
	return nil
}

func (sr *stdoutRegistry) Search(ctx context.Context, keywords ...string) ([]NewsItem, error) {
	return nil, nil
}

func (sr *stdoutRegistry) Page(ctx context.Context, offset int32, limit int32) (news []NewsItem, hasMore bool, err error) {
	return nil, false, nil
}