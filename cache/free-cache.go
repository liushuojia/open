package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/coocood/freecache"
)

var (
	_     Cache[int, any] = (*freeCache[int, any])(nil)
	cache                 = freecache.NewCache(100 * 1024 * 1024)
)

type freeCache[K KEY, V VAL] struct {
	cache   *freecache.Cache
	prefix  string
	expire  int
	keyList []string
}

/*
NewFreeCache
// prefix 前缀
// expire 有效期 单位秒
// key2Id 结构体 关键词索引 可以多个
*/
func NewFreeCache[K KEY, T VAL](prefix string, expire int, key2Id ...string) Cache[K, T] {
	if expire <= 0 {
		expire = defaultExpire
	}
	return &freeCache[K, T]{
		cache:   cache,
		prefix:  prefix,
		expire:  expire,
		keyList: key2Id,
	}
}

func (u *freeCache[K, T]) id2Value(id K) (key2Value []byte) {
	return []byte(fmt.Sprintf("%s:ID2Value:%v", u.prefix, id))
}
func (u *freeCache[K, T]) key2Id(key string) (key2ID []byte) {
	return []byte(fmt.Sprintf("%s:key2ID:%v", u.prefix, key))
}

func (u *freeCache[K, T]) Delete(ctx context.Context, id K) {
	_ = u.cache.Del(u.id2Value(id))

	if resp, err := u.Load(ctx, id); err == nil {
		val := reflect.ValueOf(resp)
		if val.Kind() == reflect.Pointer {
			val = val.Elem()
		}
		if val.Kind() == reflect.Struct {
			for _, vv := range u.keyList {
				nameField := val.FieldByName(vv)
				if nameField.IsValid() {
					_ = u.cache.Del(u.key2Id(fmt.Sprintf("%v", nameField.Interface())))
				}
			}
		}
	}
}
func (u *freeCache[K, T]) Store(ctx context.Context, id K, value T) error {
	j, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err := u.cache.Set(u.id2Value(id), j, u.expire); err != nil {
		return err
	}

	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	if val.Kind() == reflect.Struct {
		for _, vv := range u.keyList {
			nameField := val.FieldByName(vv)
			if nameField.IsValid() {
				if err := u.cache.Set(u.key2Id(fmt.Sprintf("%v", nameField.Interface())), []byte(fmt.Sprintf("%v", id)), u.expire); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
func (u *freeCache[K, T]) Load(ctx context.Context, id K) (value T, err error) {
	j, err := u.cache.Get(u.id2Value(id))
	if err != nil {
		return value, err
	}
	if err := json.Unmarshal(j, &value); err != nil {
		return value, err
	}
	return value, nil
}
func (u *freeCache[K, T]) LoadByKey(ctx context.Context, key string) (value T, err error) {
	var k K
	b, err := u.cache.Get(u.key2Id(key))
	if err != nil {
		return value, err
	}
	if err := json.Unmarshal(b, &k); err != nil {
		return value, err
	}
	return u.Load(ctx, k)
}
