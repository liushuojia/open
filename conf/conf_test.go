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

}
