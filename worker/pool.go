package worker

import (
	"errors"

	"github.com/panjf2000/ants/v2"
)

const threadPoolSize = 1000

var pool *ants.Pool

func Run(fns ...func() error) error {
	var err error
	if pool == nil {
		pool, err = ants.NewPool(threadPoolSize)
		if err != nil {
			return errors.Join(errors.New("create pool failed"), err)
		}
	}

	for _, fn := range fns {
		_ = pool.Submit(func() {
			_ = fn()
		})
	}
	return nil
}
