package token

import (
	"fmt"
	"testing"
)

func Test_JWT(t *testing.T) {
	tt := New(
		//WithIssuer("app"),
		WithExpire("2h"),
		WithSecret([]byte("secret")),
	)

	t1, err := tt.Generate(WithValueID(1))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("token", t1)

	o, err := tt.Parse(t1)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(o.ID)

}
