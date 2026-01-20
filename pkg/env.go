package pkg

import (
	"errors"
	"strconv"

	"github.com/liushuojia/open/env"
)

const envPrefix = "app."

func Set(k, v string) error {
	return env.Set(envPrefix+k, v)
}
func Get(k string) (string, bool) {
	return env.Get(envPrefix + k)
}
func GetInt64(k string) (int64, error) {
	v, ok := env.Get(envPrefix + k)
	if !ok {
		return 0, errors.New("not found")
	}
	return strconv.ParseInt(v, 10, 64)
}
func GetBool(k string) (bool, error) {
	v, ok := env.Get(envPrefix + k)
	if !ok {
		return false, errors.New("not found")
	}
	return strconv.ParseBool(v)
}
