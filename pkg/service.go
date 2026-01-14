package pkg

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

func Run() {
	if err := lcInit(); err != nil {
		log.Fatalln(err.Error())
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	lcDestroy()
	time.Sleep(time.Millisecond * 500)
}
