package pkg

import (
	"sync"

	utils "github.com/liushuojia/open"
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
	lcInit          = lc.Init
	lcDestroy       = lc.Destroy
)

type LifeCycle interface {
	RegisterInit(fn ...func() error)    // 注册服务
	RegisterDestroy(fn ...func() error) // 注销服务
	Init() error                        // 启动服务 初始化
	Destroy() error                     // 停止服务 注销
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
func (lc *lifeCycle) Init() error {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.setStatus(LifeCycleStatusInit)

	return utils.WaitError(lc.initFuncList...)
}

// Destroy 执行注册的Destroy函数，逆序执行所有Destroy函数，收集返回error聚合为errors返回
func (lc *lifeCycle) Destroy() error {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.setStatus(LifeCycleStatusDestroy)

	fns := make([]func() error, 0)
	fns = append(fns, lc.destroyFuncList...)
	for i, j := 0, len(fns)-1; i < j; i, j = i+1, j-1 {
		fns[i], fns[j] = fns[j], fns[i]
	}

	return utils.WaitError(fns...)
}
