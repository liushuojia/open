package env

import (
	"errors"
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
func Del(key string) error {
	return os.Unsetenv(key)
}
func Set(key string, value string) error {
	return os.Setenv(key, value)
}
