package conf

import (
	"encoding/json"
	"errors"
	"fmt"

	utils "github.com/liushuojia/open"
	mail "github.com/liushuojia/open/email"
	"github.com/liushuojia/open/minio"
	"github.com/liushuojia/open/token"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var defaultConfig Conf

func SetDefault(c Conf) {
	defaultConfig = c
}
func GetDefault() Conf {
	return defaultConfig
}

type Conf interface {
	Mysql(field string) (*gorm.DB, error)
	Redis(field string) (*redis.Client, error)
	Token(field string) (token.JWT, error)
	Minio(field string) (*minio.Conn, error)
	Email(field string) ([]mail.Client, error)

	GetInt64ByField(fields ...string) (value int64, err error)
	GetStringByField(fields ...string) (value string, err error)
	GetByField(value any, fields ...string) (err error)
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

	if defaultConfig == nil {
		SetDefault(c)
	}

	return c, nil
}

func (c *config) Mysql(field string) (*gorm.DB, error) {
	v, ok := c.mysql[field]
	if !ok {
		return nil, errors.New(fmt.Sprintf("not found %s", field))
	}
	return utils.MysqlConnect(v.Address, v.Username, v.Password, v.Database)
}
func (c *config) Redis(field string) (*redis.Client, error) {
	v, ok := c.redis[field]
	if !ok {
		return nil, errors.New(fmt.Sprintf("not found %s", field))
	}
	return utils.RedisConnect(v.Address, v.Password, v.DB)
}
func (c *config) Token(field string) (token.JWT, error) {
	if v, ok := c.token[field]; ok {
		return token.New(token.WithSecret([]byte(v.Key)), token.WithIssuer(v.Issuer), token.WithExpire(v.Expire)), nil
	}
	return nil, errors.New(fmt.Sprintf("not found %s", field))
}
func (c *config) Minio(field string) (*minio.Conn, error) {
	v, ok := c.minio[field]
	if !ok {
		return nil, errors.New(fmt.Sprintf("not found %s", field))
	}

	conn, err := minio.New().SetAddresses(v.Address).SetAccessKey(v.Access).SetSecretKey(v.Secret).SetUseSSL(v.UseSSL).Connect()
	if err != nil {
		return nil, err
	}

	return conn, nil
}
func (c *config) Email(field string) ([]mail.Client, error) {
	l, ok := c.email[field]
	if !ok {
		return nil, errors.New(fmt.Sprintf("not found %s", field))
	}
	mailList := make([]mail.Client, 0)
	for _, v := range l {
		mailList = append(mailList, mail.New(v.Account, v.Passwd, v.Smtp, v.Port))
	}
	return mailList, nil
}

func (c *config) getByField(fields ...string) (value any, err error) {
	var (
		valueAny any
		mapTmp   = c.m
	)

	if len(fields) <= 0 {
		return mapTmp, nil
	}

	for _, f := range fields {
		v, ok := mapTmp[f]
		if !ok {
			return 0, errors.New("not found")
		}
		if m, ok := v.(map[string]any); ok {
			mapTmp = m
		}
		valueAny = v
	}

	return valueAny, nil
}
func (c *config) GetInt64ByField(fields ...string) (value int64, err error) {
	v, err := c.getByField(fields...)
	if err != nil {
		return 0, err
	}
	if value, ok := v.(int64); ok {
		return value, nil
	}
	return 0, errors.New("field value is not number")
}
func (c *config) GetStringByField(fields ...string) (value string, err error) {
	v, err := c.getByField(fields...)
	if err != nil {
		return "", err
	}
	if value, ok := v.(string); ok {
		return value, nil
	}
	return "", errors.New("field value is not number")
}
func (c *config) GetByField(value any, fields ...string) (err error) {
	v, err := c.getByField(fields...)
	if err != nil {
		return err
	}

	j, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return json.Unmarshal(j, &value)
}
