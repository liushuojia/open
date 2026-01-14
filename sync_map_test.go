package utils

import (
	"fmt"
	"testing"
)

func Test_SyncMap_int2String(t *testing.T) {
	// 创建带监听功能的sync.Map
	m := NewSyncMap[int, string](
		// 加载时的回调
		func(key int, value string) {
			fmt.Printf("加载键值对: key=%v, value=%v\n", key, value)
		},
		// 存储时的回调
		func(key int, value string, isNew bool) {
			if isNew {
				fmt.Printf("新增键值对: key=%v, value=%v\n", key, value)
			} else {
				fmt.Printf("修改键值对: key=%v, 新value=%v\n", key, value)
			}
		},
		// 删除时的回调
		func(key int) {
			fmt.Printf("删除键值对: key=%v\n", key)
		},
	)

	// 执行一些操作
	m.Store(2, "张三")
	fmt.Println(m.Load(2))

	m.Store(2, "李四")
	m.Load(2)

	m.Store(3, "31") // 修改已存在的键
	m.Delete(3)
}
func Test_SyncMap_string2Int(t *testing.T) {
	// 创建带监听功能的sync.Map
	m := NewSyncMap[string, int](
		// 加载时的回调
		func(key string, value int) {
			fmt.Printf("加载键值对: key=%v, value=%v\n", key, value)
		},
		// 存储时的回调
		func(key string, value int, isNew bool) {
			if isNew {
				fmt.Printf("新增键值对: key=%v, value=%v\n", key, value)
			} else {
				fmt.Printf("修改键值对: key=%v, 新value=%v\n", key, value)
			}
		},
		// 删除时的回调
		func(key string) {
			fmt.Printf("删除键值对: key=%v\n", key)
		},
	)

	// 执行一些操作
	m.Store("aaa", 123)
	fmt.Println(m.Load("aaa"))

	m.Store("aaa", 234)
	m.Load("aaa")

	m.Store("bbb", 789) // 修改已存在的键
	m.Delete("bbb")
}
func Test_SyncMap_string2Struct(t *testing.T) {
	type s struct {
		A int    `json:"a"`
		B string `json:"b"`
		C bool   `json:"c"`
	}
	m := NewSyncMap[string, s](
		// 加载时的回调
		func(key string, value s) {
			fmt.Printf("加载键值对: key=%v, value=%+v\n", key, value)
		},
		// 存储时的回调
		func(key string, value s, isNew bool) {
			if isNew {
				fmt.Printf("新增键值对: key=%v, value=%+v\n", key, value)
			} else {
				fmt.Printf("修改键值对: key=%v, 新value=%+v\n", key, value)
			}
		},
		// 删除时的回调
		func(key string) {
			fmt.Printf("删除键值对: key=%v\n", key)
		},
	)

	d := s{
		A: 123,
		B: "123",
		C: true,
	}

	// 执行一些操作
	m.Store("aaa", d)
	m.Load("aaa")

	m.Store("aaa", s{
		A: 234,
		B: "234",
		C: false,
	})
	m.Load("aaa")

	m.Store("bbb", s{
		A: 567,
		B: "789",
		C: true,
	}) // 修改已存在的键
	m.Delete("bbb")
}
