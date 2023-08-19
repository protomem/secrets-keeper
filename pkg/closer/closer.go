package closer

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

type (
	Func func(ctx context.Context) error

	Closer struct {
		mux   sync.Mutex
		funcs []Func
	}
)

func New() *Closer {
	return &Closer{}
}

func (c *Closer) Add(fn Func) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.funcs = append(c.funcs, fn)
}

func (c *Closer) Close(ctx context.Context) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	var (
		msgs     = make([]string, 0, len(c.funcs))
		complete = make(chan struct{}, 1)
	)

	go func() {
		for _, f := range c.funcs {
			err := f(ctx)
			if err != nil {
				msgs = append(msgs, fmt.Sprintf("[!] %s", err.Error()))
			}
		}

		complete <- struct{}{}
	}()

	select {
	case <-complete:
		break
	case <-ctx.Done():
		return fmt.Errorf("shutdown cancelled: %w", ctx.Err())
	}

	if len(msgs) > 0 {
		return fmt.Errorf(
			"shutdown finished with error(s): %s",
			strings.Join(msgs, "; "),
		)
	}

	return nil
}
