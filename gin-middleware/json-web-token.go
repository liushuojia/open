package gin_middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liushuojia/open/token"
)

const (
	headerTokenKey = "authorization"
	claimsKey      = "claims-key"
)

type JsonWebTokenMiddleware struct {
	token token.JWT
}

func NewJsonWebTokenMiddleware(tokenItem token.JWT) *JsonWebTokenMiddleware {
	return &JsonWebTokenMiddleware{
		token: tokenItem,
	}
}

func (m *JsonWebTokenMiddleware) Handle() func(c *gin.Context) {
	return func(c *gin.Context) {
		accessToken := c.GetHeader(headerTokenKey)
		if accessToken == "" {
			c.String(http.StatusUnauthorized, "token is empty")
			c.Abort()
			return
		}

		claims, err := m.token.Parse(accessToken)
		if err != nil {
			c.String(http.StatusUnauthorized, "token is wrong")
			c.Abort()
			return
		}

		c.Set(claimsKey, claims)
		c.Next()
	}
}

func GetClaimsByGinContext(c *gin.Context) (*token.Claims, error) {
	dataInterface, ok := c.Get(claimsKey)
	if !ok {
		return nil, errors.New("claim is empty")
	}

	data, ok := dataInterface.(*token.Claims)
	if !ok || data == nil {
		return nil, errors.New("claim is empty")
	}

	return data, nil
}
func MustGetClaimsByGinContext(c *gin.Context) *token.Claims {
	resp, err := GetClaimsByGinContext(c)
	if err != nil {
		panic(err)
	}
	return resp
}
