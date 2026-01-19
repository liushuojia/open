package pkg

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func Test_Register(t *testing.T) {

	RegisterInit(func() error {
		log.Println("init 1")
		return nil
	}, func() error {
		log.Println("init 2")
		return nil
	})
	RegisterInit(func() error {
		log.Println("init 3")
		return nil
	}, func() error {
		log.Println("init 4")
		return nil
	})

	RegisterDestroy(func() error {
		log.Println("Destroy 1")
		return nil
	}, func() error {
		log.Println("Destroy 2")
		return nil
	})
	RegisterDestroy(func() error {
		log.Println("Destroy 3")
		return nil
	}, func() error {
		log.Println("Destroy 4")
		return nil
	})

	Run(WithConfig(
		"../conf/conf.daemon.default.toml",
		"../conf/conf.daemon.mysql.toml",
		"../conf/conf.daemon.redis.toml",
		"../conf/conf.daemon.toml",
	))
}
