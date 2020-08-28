package registry

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/jinzhu/gorm"
	"io"
)

type MysqlNewsItem struct {
	ID int64 		  `json:"id" gorm:"column:id;AUTO_INCREMENT;PRIMARY_KEY"`
	Title string      `json:"title" gorm:"column:title;type:varchar(128)"`
	Href string       `json:"href" gorm:"column:href;type:varchar(512)"`
	SpiderHref string `json:"spider_href" gorm:"column:spider_href;type:varchar(512)"`
	Sign string 	  `json:"sign" gorm:"column:sign;type:varchar(64);UNIQUE"`
	CreateTime int64  `json:"create_time" gorm:"column:create_time"`
}

func (MysqlNewsItem) TableName() string {
	return "news"
}

type mysqlRegistry struct {
	db *gorm.DB
}

func NewMsqlRegistry(db *gorm.DB) Registry {
	mr := &mysqlRegistry{
		db:db,
	}
	mr.db.AutoMigrate(MysqlNewsItem{})
	return mr
}

func (mr *mysqlRegistry) Store(ctx context.Context, news ...NewsItem) error {
	var buff bytes.Buffer
	args := make([]interface{}, 0, len(news)*5)
	buff.WriteString("insert ignore into news(`title`,`href`,`spider_href`,`sign`,`create_time`) values ")
	for i,cnt:=0,len(news);i<cnt;i++{
		item := mr.transNewsItemToMysqlNewsItem(news[i])
		buff.WriteString("(?,?,?,?,?)")
		if i == cnt - 1 {
			buff.WriteString(";")
		}  else {
			buff.WriteString(",")
		}
		args = append(args, item.Title, item.Href, item.SpiderHref, item.Sign, item.CreateTime)
	}
	return mr.db.Exec(buff.String(), args...).Error
}

func (mr *mysqlRegistry) Search(ctx context.Context, keywords ...string) ([]NewsItem, error) {
	return nil, nil
}

func (mr *mysqlRegistry) Page(ctx context.Context, offset int32, limit int32) (news []NewsItem, hasMore bool, err error) {
	return nil, false, nil
}

func (mr *mysqlRegistry) transNewsItemToMysqlNewsItem(newsItem NewsItem) MysqlNewsItem {
	return MysqlNewsItem{
		Title:newsItem.Title,
		Href:newsItem.Href,
		SpiderHref:newsItem.SpiderHref,
		CreateTime:newsItem.CreateTime,
		Sign:mr.genSign(newsItem.Href),
	}
}

func (mr *mysqlRegistry) genSign(source string) string {
	w := md5.New()
	io.WriteString(w, source)
	return fmt.Sprintf("%X", w.Sum(nil))
}

