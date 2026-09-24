package utils

import (
	"sync"
	"time"
)

type StackDebouncer[T any] struct {
	mu       sync.Mutex
	timer    *time.Timer
	delay    time.Duration
	callback func([]T)
	argsList []T
}

func NewStackDebouncer[T any](delay time.Duration, callback func([]T)) *StackDebouncer[T] {
	return &StackDebouncer[T]{
		delay:    delay,
		callback: callback,
		argsList: make([]T, 0),
	}
}

func (d *StackDebouncer[T]) Push(arg T) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.argsList = append(d.argsList, arg)

	if d.timer != nil {
		d.timer.Stop()
	}

	d.timer = time.AfterFunc(d.delay, func() {
		d.mu.Lock()
		argsToProcess := d.argsList
		d.argsList = make([]T, 0)
		d.mu.Unlock()

		if len(argsToProcess) > 0 {
			d.callback(argsToProcess)
		}
	})
}
