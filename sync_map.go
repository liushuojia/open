package utils

import (
	"sync"
)

// 定义监听回调函数类型
type (
	// LoadMapFunc 当键值对被加载时触发
	LoadMapFunc[K, T any] func(key K, value T)
	// StoreMapFunc 当键值对被存储时触发（包括新增和修改）
	StoreMapFunc[K, T any] func(key K, value T, isNew bool)
	// DeleteMapFunc 当键值对被删除时触发
	DeleteMapFunc[K any] func(key K)
)

// SyncMap 带监听功能的sync.Map封装
type SyncMap[K, T any] struct {
	m        sync.Map
	onLoad   LoadMapFunc[K, T]  // 加载时的回调
	onStore  StoreMapFunc[K, T] // 存储时的回调
	onDelete DeleteMapFunc[K]   // 删除时的回调
}

// NewSyncMap 创建一个新的带监听功能的sync.Map
func NewSyncMap[K, T any](loadFunc LoadMapFunc[K, T], storeFunc StoreMapFunc[K, T], deleteFunc DeleteMapFunc[K]) *SyncMap[K, T] {
	return &SyncMap[K, T]{
		onLoad:   loadFunc,
		onStore:  storeFunc,
		onDelete: deleteFunc,
	}
}

// Load 加载键值对并触发监听
func (w *SyncMap[K, T]) Load(key K) (T, bool) {
	var (
		ok    bool
		ok1   bool
		value any
		resp  T
	)

	value, ok = w.m.Load(key)
	if ok && w.onLoad != nil {
		resp, ok1 = value.(T)
		if ok1 {
			w.onLoad(key, resp)
		}
	}

	return resp, ok1
}

// Store 存储键值对并触发监听
func (w *SyncMap[K, T]) Store(key K, value T) {
	// 先检查是否存在
	_, exists := w.m.Load(key)
	w.m.Store(key, value)

	if w.onStore != nil {
		w.onStore(key, value, !exists)
	}
}

// Delete 删除键值对并触发监听
func (w *SyncMap[K, T]) Delete(key K) {
	// 先检查是否存在
	if _, exists := w.m.Load(key); exists {
		w.m.Delete(key)
		if w.onDelete != nil {
			w.onDelete(key)
		}
	}
}
