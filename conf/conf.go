package conf

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	utils "github.com/liushuojia/open"
	mail "github.com/liushuojia/open/email"
	"github.com/liushuojia/open/minio"
	"github.com/liushuojia/open/token"
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
	GetMapByField(fields ...string) (map[string]any, error)

	Stop() error
}

type config struct {
	m map[string]any

	mysql map[string]*gorm.DB
	redis map[string]*redis.Client
	token map[string]token.JWT
	minio map[string]*minio.Conn
	email map[string][]mail.Client
}

func New(options ...Option) (Conf, error) {
	opts := loadOptions(options...)

	c := &config{
		m: make(map[string]any),

		mysql: make(map[string]*gorm.DB),
		redis: make(map[string]*redis.Client),
		token: make(map[string]token.JWT),
		minio: make(map[string]*minio.Conn),
		email: make(map[string][]mail.Client),
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
			conn, err := utils.MysqlConnect(v.Address, v.Username, v.Password, v.Database)
			if err != nil {
				return nil, err
			}
			c.mysql[k] = conn
		}
		for k, v := range dataValue.Redis {
			conn, err := utils.RedisConnect(v.Address, v.Password, v.DB)
			if err != nil {
				return nil, err
			}
			c.redis[k] = conn
		}
		for k, v := range dataValue.Token {
			c.token[k] = token.New(token.WithSecret([]byte(v.Key)), token.WithIssuer(v.Issuer), token.WithExpire(v.Expire))
		}
		for k, v := range dataValue.Minio {
			conn, err := minio.New().SetAddresses(v.Address).SetAccessKey(v.Access).SetSecretKey(v.Secret).SetUseSSL(v.UseSSL).Connect()
			if err != nil {
				return nil, err
			}
			c.minio[k] = conn
		}
		for k, v := range dataValue.Email {
			mailList := make([]mail.Client, 0)
			for _, vv := range v {
				mailList = append(mailList, mail.New(vv.Account, vv.Passwd, vv.Smtp, vv.Port))
			}
			c.email[k] = mailList
		}
	}

	if defaultConfig == nil {
		SetDefault(c)
	}

	return c, nil
}

func (c *config) Mysql(field string) (*gorm.DB, error) {
	if v, ok := c.mysql[field]; ok {
		return v, nil
	}
	return nil, errors.New(fmt.Sprintf("not found %s", field))
}
func (c *config) Redis(field string) (*redis.Client, error) {
	if v, ok := c.redis[field]; ok {
		return v, nil
	}
	return nil, errors.New(fmt.Sprintf("not found %s", field))
}
func (c *config) Token(field string) (token.JWT, error) {
	if v, ok := c.token[field]; ok {
		return v, nil
	}
	return nil, errors.New(fmt.Sprintf("not found %s", field))
}
func (c *config) Minio(field string) (*minio.Conn, error) {
	if v, ok := c.minio[field]; ok {
		return v, nil
	}
	return nil, errors.New(fmt.Sprintf("not found %s", field))
}
func (c *config) Email(field string) ([]mail.Client, error) {
	if v, ok := c.email[field]; ok && len(v) > 0 {
		return v, nil
	}
	return nil, errors.New(fmt.Sprintf("not found %s", field))
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
func (c *config) GetByField(value any, fields ...string) error {
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
func (c *config) GetMapByField(fields ...string) (map[string]any, error) {
	v, err := c.getByField(fields...)
	if err != nil {
		return nil, err
	}
	mv, ok := v.(map[string]any)
	if !ok {
		return nil, errors.New("value is not map[string]any")
	}
	return mv, nil
}

func (c *config) Stop() error {
	fmt.Println("clean service with conf")
	for _, gormDB := range c.mysql {
		if db, err := gormDB.DB(); err == nil {
			_ = db.Close()
		}
	}
	for _, redisDB := range c.redis {
		if redisDB != nil {
			_ = redisDB.Close()
		}
	}
	return nil
}
