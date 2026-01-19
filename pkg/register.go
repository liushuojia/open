package pkg

import (
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/liushuojia/open/conf"
	log "github.com/sirupsen/logrus"
)

const (
	LifeCycleStatusUnknown = iota
	LifeCycleStatusInit
	LifeCycleStatusStart
	LifeCycleStatusStop
	LifeCycleStatusDestroy
)

var (
	lc              = NewLifeCycle()
	RegisterInit    = lc.RegisterInit
	RegisterDestroy = lc.RegisterDestroy
	Run             = lc.Run
)

type LifeCycle interface {
	RegisterInit(fn ...func() error)    // 注册服务
	RegisterDestroy(fn ...func() error) // 注销服务
	Run(opts ...Option)                 // 启动服务
}

type lifeCycle struct {
	mu              sync.Mutex
	initFuncList    []func() error
	destroyFuncList []func() error
	status          int
}

func NewLifeCycle() LifeCycle {
	return &lifeCycle{
		initFuncList:    make([]func() error, 0),
		destroyFuncList: make([]func() error, 0),
		status:          LifeCycleStatusInit,
	}
}

func (lc *lifeCycle) setStatus(status int) {
	lc.status = status
}

func (lc *lifeCycle) RegisterInit(fn ...func() error) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	if lc.status != LifeCycleStatusInit {
		panic("failed to register init function, must be called in init")
	}
	lc.initFuncList = append(lc.initFuncList, fn...)
}
func (lc *lifeCycle) RegisterDestroy(fn ...func() error) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.destroyFuncList = append(lc.destroyFuncList, fn...)
}

// Init 执行注册的初始化函数，顺序执行直到初始化函数返回error，并将error返回
func (lc *lifeCycle) init() error {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.setStatus(LifeCycleStatusInit)

	var err error
	for _, fn := range lc.initFuncList {
		if e := fn(); e != nil {
			err = errors.Join(err, e)
		}
	}

	return err
}

// Destroy 执行注册的Destroy函数，逆序执行所有Destroy函数，收集返回error聚合为errors返回
func (lc *lifeCycle) destroy() error {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.setStatus(LifeCycleStatusDestroy)

	fns := make([]func() error, 0)
	fns = append(fns, lc.destroyFuncList...)
	for i, j := 0, len(fns)-1; i < j; i, j = i+1, j-1 {
		fns[i], fns[j] = fns[j], fns[i]
	}

	var err error
	for _, fn := range fns {
		if e := fn(); e != nil {
			err = errors.Join(err, e)
		}
	}

	return err
}

// Run 启动服务
func (lc *lifeCycle) Run(opts ...Option) {

	var (
		c   conf.Conf
		err error
	)
	opt := loadOptions(opts...)
	if len(opt.pathList) > 0 {
		log.Println("start with config file ", opt.pathList)
		c, err = conf.New(conf.WithFilePath(opt.pathList...))
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
	}

	if opt.recover {
		defer func() {
			if err := recover(); err != nil {
				log.Println("panic recover", err)
			}
		}()
	}

	if err := lc.init(); err != nil {
		log.Fatalln(err.Error())
		return
	}

	if opt.noWaiting {
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	_ = lc.destroy()
	time.Sleep(time.Millisecond * 300)

	if c != nil {
		log.Println("stop service")
		_ = c.Stop()
	}
	time.Sleep(time.Millisecond * 300)
}
