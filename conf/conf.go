package conf

import (
	"errors"
	"fmt"

	utils "github.com/liushuojia/open"
)

type Conf interface {
	Mysql(field string) (*Mysql, error)
	Redis(field string) (*Redis, error)
	Token(field string) (*Token, error)
	Minio(field string) (*Minio, error)
	Email(field string) ([]*Email, error)
}

type config struct {
	m     map[string]any
	mysql map[string]*Mysql
	redis map[string]*Redis
	token map[string]*Token
	minio map[string]*Minio
	email map[string][]*Email
}

func New(options ...Option) (Conf, error) {
	opts := loadOptions(options...)

	c := &config{
		m:     make(map[string]any),
		mysql: make(map[string]*Mysql),
		redis: make(map[string]*Redis),
		token: make(map[string]*Token),
		minio: make(map[string]*Minio),
		email: make(map[string][]*Email),
	}
	for _, filePath := range opts.filePath {
		mapValue := make(map[string]any)
		if err := utils.Read(filePath, &mapValue); err != nil {
			return nil, errors.Join(errors.New(fmt.Sprintf("read file %s fail", filePath)), err)
		}
		for k, v := range mapValue {
			c.m[k] = v
		}

		var dataValue struct {
			Mysql map[string]*Mysql   `toml:"mysql"`
			Redis map[string]*Redis   `toml:"redis"`
			Token map[string]*Token   `toml:"token"`
			Minio map[string]*Minio   `toml:"minio"`
			Email map[string][]*Email `toml:"email"`
		}
		if err := utils.Read(filePath, &dataValue); err != nil {
			return nil, errors.Join(errors.New(fmt.Sprintf("read file %s fail", filePath)), err)
		}
		for k, v := range dataValue.Mysql {
			c.mysql[k] = v
		}
		for k, v := range dataValue.Redis {
			c.redis[k] = v
		}
		for k, v := range dataValue.Token {
			c.token[k] = v
		}
		for k, v := range dataValue.Minio {
			c.minio[k] = v
		}
		for k, v := range dataValue.Email {
			c.email[k] = v
		}
	}

	return c, nil
}

func (c *config) Mysql(field string) (*Mysql, error) {
	if v, ok := c.mysql[field]; ok {
		return v, nil
	}
	return nil, errors.New(fmt.Sprintf("not found %s", field))
}
func (c *config) Redis(field string) (*Redis, error) {
	if v, ok := c.redis[field]; ok {
		return v, nil
	}
	return nil, errors.New(fmt.Sprintf("not found %s", field))
}
func (c *config) Token(field string) (*Token, error) {
	if v, ok := c.token[field]; ok {
		return v, nil
	}
	return nil, errors.New(fmt.Sprintf("not found %s", field))
}
func (c *config) Minio(field string) (*Minio, error) {
	if v, ok := c.minio[field]; ok {
		return v, nil
	}
	return nil, errors.New(fmt.Sprintf("not found %s", field))
}
func (c *config) Email(field string) ([]*Email, error) {
	if v, ok := c.email[field]; ok {
		return v, nil
	}
	return nil, errors.New(fmt.Sprintf("not found %s", field))
}
