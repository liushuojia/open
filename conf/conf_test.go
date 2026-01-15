package conf

import (
	"fmt"
	"testing"
)

func Test_Client(t *testing.T) {
	cc, err := New(WithFilePath(
		"./conf.daemon.toml",
		"./conf.daemon.default.toml",
		"./conf.daemon.mysql.toml",
		"./conf.daemon.redis.toml",
	))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if v, err := cc.Mysql("new_example"); err == nil {
		fmt.Println(v)
	}

	fmt.Println(cc.GetInt64ByField("aa"))
	fmt.Println(cc.GetInt64ByField("cc"))
	fmt.Println(cc.GetInt64ByField("bb", "cc"))
	fmt.Println(cc.GetInt64ByField("mysql", "example", "database"))
	fmt.Println(cc.GetInt64ByField("redis", "1", "db"))

	var ccc struct {
		Name          string `json:"name"`
		Port          int    `json:"port"`
		ApiPrefix     string `json:"apiPrefix"`
		Debug         bool   `json:"debug"`
		FileSharePath string `json:"fileSharePath"`
	}
	fmt.Println(cc.GetByField(&ccc))
	fmt.Println(ccc)

}
