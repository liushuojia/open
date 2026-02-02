package cache

import "context"

type (
	KEY interface {
		int | int8 | int16 | int32 | int64 |
			uint | uint8 | uint16 | uint32 | uint64 |
			string
	}
	VAL any
)

type Cache[K KEY, V VAL] interface {
	Delete(ctx context.Context, id K)
	Store(ctx context.Context, id K, value V) error
	Load(ctx context.Context, id K) (V, error)
	LoadByKey(ctx context.Context, key string) (V, error)
}
