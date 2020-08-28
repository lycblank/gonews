package conf

import "fmt"

type MsqlConfig struct {
	Username string	`yaml:"username"`
	Password string	`yaml:"password"`
	Addr string	`yaml:"addr"`
	DBName string `yaml:"dbname"`
	Param string `yaml:"param"`
}

func (mc MsqlConfig) String() string {
	return fmt.Sprintf(`%s:%s@tcp(%s)/%s?%s`, mc.Username,mc.Password,mc.Addr,mc.DBName,mc.Param)
}