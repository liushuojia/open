package token

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

// Claims 加密主体, 根据实际情况修改
type Claims struct {
	ID    uint64 `json:"i,omitempty"` // ID
	AppID string `json:"a,omitempty"` //
	Key   string `json:"k,omitempty"` //

	jwt.RegisteredClaims // 嵌入标准声明
}

var _ JWT = (*jsonWebToken)(nil)

type JWT interface {
	Generate(options ...ValueOption) (token string, err error)
	Parse(token string) (*Claims, error)
}

type jsonWebToken struct {
	Issuer  string
	Subject string
	expire  string
	Secret  []byte
}

func New(options ...Option) JWT {
	opt := loadOptions(options...)
	return &jsonWebToken{
		Subject: opt.Subject,
		Issuer:  opt.Issuer,
		expire:  opt.Expire,
		Secret:  opt.Secret,
	}
}

// Generate 生成token
func (j *jsonWebToken) Generate(options ...ValueOption) (token string, err error) {
	// 设置过期时间（例如：2小时后过期）
	m, err := time.ParseDuration(j.expire)
	if err != nil {
		return "", err
	}

	expirationTime := time.Now().Add(m)

	opt := loadValueOptions(options...)

	// 创建自定义 Claims
	claims := Claims{
		ID:    opt.ID,
		AppID: opt.AppID,
		Key:   opt.Key,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),     // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),     // 生效时间（立即生效）
			Issuer:    j.Issuer,                           // 签发者
			Subject:   j.Subject,                          // 主题
		},
	}

	// 使用 HS256 算法创建 token 对象
	tokenTmp := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并生成 token 字符串
	tokenString, err := tokenTmp.SignedString(j.Secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Parse 验证token
func (j *jsonWebToken) Parse(token string) (*Claims, error) {
	if token == "" {
		return nil, errors.New("token为空")
	}

	// 解析 token
	tokenTmp, err := jwt.ParseWithClaims(
		token,
		&Claims{}, // 用于接收解析后的自定义 claims
		func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return j.Secret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	// 验证 token 并提取 claims
	if claims, ok := tokenTmp.Claims.(*Claims); ok && tokenTmp.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
