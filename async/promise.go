package async

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

// State represents the current lifecycle state of a Promise.
type State int

const (
	StatePending State = iota
	StateFulfilled
	StateRejected
)

// String returns the string representation of a State.
func (s State) String() string {
	switch s {
	case StatePending:
		return "Pending"
	case StateFulfilled:
		return "Fulfilled"
	case StateRejected:
		return "Rejected"
	default:
		return "Unknown"
	}
}

// PanicError represents an error caused by a recovered panic inside an async function.
type PanicError struct {
	Value any
	Stack []byte
}

func (e *PanicError) Error() string {
	return fmt.Sprintf("promise panicked: %v", e.Value)
}

// Promise represents the eventual completion (or failure) of an asynchronous operation
// and its resulting value.
type Promise[T any] struct {
	mu    sync.RWMutex
	done  chan struct{}
	state State
	value T
	err   error
}

// Async starts a function asynchronously in a separate goroutine and returns a Promise.
// If the function panics, the panic is recovered and stored as a *PanicError.
func Async[T any](fn func() (T, error)) *Promise[T] {
	return AsyncWithContext(context.Background(), func(ctx context.Context) (T, error) {
		return fn()
	})
}

// AsyncWithContext starts a context-aware function asynchronously in a separate goroutine.
func AsyncWithContext[T any](ctx context.Context, fn func(ctx context.Context) (T, error)) *Promise[T] {
	p := &Promise[T]{
		done:  make(chan struct{}),
		state: StatePending,
	}

	go func() {
		defer close(p.done)
		defer func() {
			if r := recover(); r != nil {
				p.mu.Lock()
				p.err = &PanicError{
					Value: r,
					Stack: debug.Stack(),
				}
				p.state = StateRejected
				p.mu.Unlock()
			}
		}()

		res, err := fn(ctx)
		p.mu.Lock()
		if err != nil {
			p.err = err
			p.state = StateRejected
		} else {
			p.value = res
			p.state = StateFulfilled
		}
		p.mu.Unlock()
	}()

	return p
}

// Resolved creates an immediately fulfilled Promise with the given value.
func Resolved[T any](val T) *Promise[T] {
	p := &Promise[T]{
		done:  make(chan struct{}),
		state: StateFulfilled,
		value: val,
	}
	close(p.done)
	return p
}

// Rejected creates an immediately rejected Promise with the given error.
func Rejected[T any](err error) *Promise[T] {
	p := &Promise[T]{
		done:  make(chan struct{}),
		state: StateRejected,
		err:   err,
	}
	close(p.done)
	return p
}

// Await blocks until the Promise settles or the provided context is canceled.
// It can be called multiple times sequentially or concurrently by multiple goroutines;
// once settled, all callers receive the same result without consuming it.
func (p *Promise[T]) Await(ctx context.Context) (T, error) {
	select {
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	case <-p.done:
		p.mu.RLock()
		defer p.mu.RUnlock()
		return p.value, p.err
	}
}

// State returns the current State of the Promise without blocking.
func (p *Promise[T]) State() State {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.state
}

// IsSettled returns true if the Promise is fulfilled or rejected.
func (p *Promise[T]) IsSettled() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.state != StatePending
}

// IsFulfilled returns true if the Promise was fulfilled.
func (p *Promise[T]) IsFulfilled() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.state == StateFulfilled
}

// IsRejected returns true if the Promise was rejected.
func (p *Promise[T]) IsRejected() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.state == StateRejected
}

// Result returns the value, error, and whether the Promise is settled, without blocking.
func (p *Promise[T]) Result() (T, error, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.state == StatePending {
		var zero T
		return zero, nil, false
	}
	return p.value, p.err, true
}

// Then chains a transformation function onto the Promise.
// If the parent Promise succeeds, fn is executed with the result.
// If the parent Promise fails, the error is propagated to the new Promise without executing fn.
func (p *Promise[T]) Then[U any](fn func(T) (U, error)) *Promise[U] {
	return Async(func() (U, error) {
		var zeroU U
		res, err := p.Await(context.Background())
		if err != nil {
			return zeroU, err
		}
		return fn(res)
	})
}

// Catch attaches a failure callback to the Promise.
// If the parent Promise fails, fn is executed with the error, allowing recovery or error translation.
// If the parent Promise succeeds, its value is passed through directly.
func (p *Promise[T]) Catch(fn func(error) (T, error)) *Promise[T] {
	return Async(func() (T, error) {
		res, err := p.Await(context.Background())
		if err == nil {
			return res, nil
		}
		return fn(err)
	})
}

// Finally attaches a callback that runs regardless of whether the Promise was fulfilled or rejected.
// The resulting Promise resolves to the same value or error as the original Promise.
func (p *Promise[T]) Finally(fn func()) *Promise[T] {
	return Async(func() (T, error) {
		defer fn()
		return p.Await(context.Background())
	})
}

// WithTimeout returns a new Promise that fails with context.DeadlineExceeded if the parent Promise
// does not settle within timeout.
func (p *Promise[T]) WithTimeout(timeout time.Duration) *Promise[T] {
	return Async(func() (T, error) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		return p.Await(ctx)
	})
}
