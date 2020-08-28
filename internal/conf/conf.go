package conf

import (
	"io/ioutil"
	"os"
	"sync"
"gopkg.in/yaml.v2"
)

type Config struct {
	Mysql MsqlConfig `yaml:"mysql"`
}

var cfg *Config
var once sync.Once
func GetConfig() *Config {
	once.Do(func(){
		InitConfig()
	})
	return cfg
}

func InitConfig() {
	configPath := os.Getenv("GONEWS_CONFIG")
	if configPath == "" {
		configPath = "configs/service.yaml"
	}
	datas, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	tmp := &Config{}
	err = yaml.Unmarshal(datas, tmp)
	if err != nil {
		panic(err)
	}
	cfg = tmp
}

