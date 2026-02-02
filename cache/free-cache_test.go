package cache

import (
	"context"
	"fmt"
	"testing"
)

func Test_free_cache(t *testing.T) {
	cc := NewFreeCache[uint, string]("pp", 0)
	ctx := context.Background()

	fmt.Println(cc.Load(ctx, 11))
	fmt.Println(cc.Store(ctx, 11, "aa"))
	fmt.Println(cc.Load(ctx, 11))

}
