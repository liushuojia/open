package conf

import utils "github.com/liushuojia/open"

type Conf interface {
}

type config struct {
	Mysql map[string]Mysql   `toml:"mysql"`
	Redis map[string]Redis   `toml:"redis"`
	Token map[string]Token   `toml:"token"`
	Minio map[string]Minio   `toml:"minio"`
	Email map[string][]Email `toml:"email"`
}

func New(options ...Option) (Conf, error) {
	opts := loadOptions(options...)

	c := &config{}
	if err := utils.Read(opts.filePath, c); err != nil {
		return nil, err
	}

	return c, nil
}
