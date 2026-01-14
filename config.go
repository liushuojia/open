package utils

import (
	"github.com/BurntSushi/toml"
)

func Read(filePath string, data any) error {
	_, err := toml.DecodeFile(filePath, data)
	return err
}
