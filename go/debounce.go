package debounce

import (
	"context"
	"time"
)

type Debouncer[T any, R any] struct {
	delay   time.Duration
	f       func(T) R
	Input   chan T
	promise chan R
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewDebouncer[T, R any](ctx context.Context, f func(params T) R, delay time.Duration) *Debouncer[T, R] {

	promise := make(chan R, 1)
	input := make(chan T, 1)

	ctx, cancel := context.WithCancel(ctx)

	d := &Debouncer[T, R]{
		delay:   delay,
		ctx:     ctx,
		cancel:  cancel,
		Input:   input,
		f:       f,
		promise: promise,
	}

	go d.run()

	return d
}

func (d *Debouncer[T, R]) run() {
	var (
		latestParams T
		timer        *time.Timer
	)
	for {
		select {
		case params := <-d.Input:
			latestParams = params

			if timer != nil {
				timer.Stop()
			}

			timer = time.AfterFunc(d.delay, func() {
				result := d.f(latestParams)
				d.promise <- result
				d.cancel()
			})

		case <-d.ctx.Done():
			if timer != nil {
				timer.Stop()
			}
			close(d.Input)
			close(d.promise)
			return
		}

	}
}

func (d *Debouncer[T, R]) Promise() <-chan R {
	return d.promise
}
