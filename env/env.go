package env

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func Get(key string) (string, bool) {
	return os.LookupEnv(key)
}
func GetInt64(key string) (int64, error) {
	v, ok := Get(key)
	if !ok {
		return 0, errors.New("env not found")
	}
	return strconv.ParseInt(v, 10, 64)
}
func GetBool(k string) (bool, error) {
	v, ok := Get(k)
	if !ok {
		return false, errors.New("not found")
	}
	return strconv.ParseBool(v)
}

func Del(key string) error {
	return os.Unsetenv(key)
}

func Set(key string, value string) error {
	return os.Setenv(key, value)
}
func SetInt64(key string, value int64) error {
	return os.Setenv(key, fmt.Sprintf("%d", value))
}
func SetBool(key string, value bool) error {
	return os.Setenv(key, fmt.Sprintf("%v", value))
}
