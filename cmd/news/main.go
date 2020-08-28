package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/lycblank/gonews/internal/conf"
	"github.com/lycblank/gonews/internal/spider"
	"github.com/lycblank/gonews/pkg/registry"
	"go.uber.org/dig"
)

var container *dig.Container
func init() {
	container = dig.New()
	if err := container.Provide(conf.GetConfig); err != nil {
		panic(err)
	}
	if err := container.Provide(ProvideDB); err != nil {
		panic(err)
	}
	if err :=  container.Provide(ProvideRegistry); err != nil {
		panic(err)
	}
	if err := container.Provide(ProvideSpider); err != nil {
		panic(err)
	}
}

func main() {
	if err := container.Invoke(Run); err != nil {
		panic(err)
	}
}

func ProvideDB(cfg *conf.Config) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", cfg.Mysql.String())
	return db, err
}

func ProvideRegistry(db *gorm.DB) registry.Registry {
	mr := registry.NewMsqlRegistry(db)
	return mr
}

func ProvideSpider(rr registry.Registry) *spider.NewsSpider {
	return spider.NewNewsSpider(rr)
}

func Run(ns *spider.NewsSpider) {
	ns.Run()
}

