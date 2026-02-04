package pkg

import (
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
	return env.GetInt64(envPrefix + k)
}
func GetBool(k string) (bool, error) {
	return env.GetBool(envPrefix + k)
}
