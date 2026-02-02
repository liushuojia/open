package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	defaultExpire       = 5 * 60
	defaultExpireSecond = defaultExpire * time.Second
)

var (
	_ Cache[int, any] = (*redisCache[int, any])(nil)
)

type redisCache[K KEY, V VAL] struct {
	redis   *redis.Client
	prefix  string
	expire  time.Duration
	keyList []string
}

/*
NewRedisCache
// prefix 前缀
// expire 有效期
// key2Id 结构体 关键词索引 可以多个
*/
func NewRedisCache[K KEY, T VAL](rds *redis.Client, prefix string, expire time.Duration, key2Id ...string) Cache[K, T] {
	if expire <= 0 {
		expire = defaultExpireSecond
	}
	return &redisCache[K, T]{
		redis:   rds,
		prefix:  prefix,
		expire:  expire,
		keyList: key2Id,
	}
}

func (u *redisCache[K, T]) id2Value(id K) (key2Value string) {
	return fmt.Sprintf("%s:ID2Value:%v", u.prefix, id)
}
func (u *redisCache[K, T]) key2Id(key string) (key2ID string) {
	return fmt.Sprintf("%s:key2ID:%v", u.prefix, key)
}

func (u *redisCache[K, T]) Delete(ctx context.Context, id K) {
	key2ValueDeleteList := make([]string, 0)
	key2ValueDeleteList = append(key2ValueDeleteList, u.id2Value(id))

	if resp, err := u.Load(ctx, id); err == nil {
		val := reflect.ValueOf(resp)
		if val.Kind() == reflect.Pointer {
			val = val.Elem()
		}
		if val.Kind() == reflect.Struct {
			for _, vv := range u.keyList {
				nameField := val.FieldByName(vv)
				if nameField.IsValid() {
					key2ValueDeleteList = append(key2ValueDeleteList, u.key2Id(fmt.Sprintf("%v", nameField.Interface())))
				}
			}
		}
	}

	_ = u.redis.Del(ctx, key2ValueDeleteList...).Err()
}
func (u *redisCache[K, T]) Store(ctx context.Context, id K, value T) error {
	j, err := json.Marshal(value)
	if err != nil {
		return err
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	expiration := time.Duration(r.Intn(6000))*time.Millisecond + u.expire
	if err := u.redis.Set(ctx, u.id2Value(id), j, expiration).Err(); err != nil {
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
				if err := u.redis.Set(ctx, u.key2Id(fmt.Sprintf("%v", nameField.Interface())), id, expiration).Err(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
func (u *redisCache[K, T]) Load(ctx context.Context, id K) (value T, err error) {
	j, err := u.redis.Get(ctx, u.id2Value(id)).Bytes()
	if err != nil {
		return value, err
	}
	if err := json.Unmarshal(j, &value); err != nil {
		return value, err
	}
	return value, nil
}
func (u *redisCache[K, T]) LoadByKey(ctx context.Context, key string) (value T, err error) {
	var k K
	b, err := u.redis.Get(ctx, u.key2Id(key)).Bytes()
	if err != nil {
		return value, err
	}
	if err := json.Unmarshal(b, &k); err != nil {
		return value, err
	}
	return u.Load(ctx, k)
}
