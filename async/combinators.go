package async 

import (
	"context"
	"errors"
	"strings"
	"sync"
)

// AggregateError represents a collection of errors when all promises fail in Any.
type AggregateError struct {
	Errors []error
}

func (e *AggregateError) Error() string {
	var sb strings.Builder
	sb.WriteString("all promises were rejected:")
	for _, err := range e.Errors {
		sb.WriteString(" [")
		if err != nil {
			sb.WriteString(err.Error())
		}
		sb.WriteString("]")
	}
	return sb.String()
}

// SettledResult holds the terminal state and value or error of an awaited Promise in AllSettled.
type SettledResult[T any] struct {
	State State
	Value T
	Err   error
}

// All waits for all promises to resolve. If any promise is rejected,
// the returned Promise is immediately rejected with that error (fail-fast).
// The results slice preserves the exact order of the input promises.
func All[T any](promises ...*Promise[T]) *Promise[[]T] {
	if len(promises) == 0 {
		return Resolved([]T{})
	}

	return AsyncWithContext(context.Background(), func(ctx context.Context) ([]T, error) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		results := make([]T, len(promises))
		errChan := make(chan error, 1)
		var wg sync.WaitGroup
		wg.Add(len(promises))

		for i, p := range promises {
			go func(idx int, prom *Promise[T]) {
				defer wg.Done()
				val, err := prom.Await(ctx)
				if err != nil {
					select {
					case errChan <- err:
						cancel()
					default:
					}
					return
				}
				results[idx] = val
			}(i, p)
		}

		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-ctx.Done():
			select {
			case err := <-errChan:
				return nil, err
			default:
				return nil, ctx.Err()
			}
		case err := <-errChan:
			return nil, err
		case <-done:
			select {
			case err := <-errChan:
				return nil, err
			default:
				return results, nil
			}
		}
	})
}

// Race returns a Promise that settles with the result (value or error)
// of the first promise in the slice to settle.
func Race[T any](promises ...*Promise[T]) *Promise[T] {
	if len(promises) == 0 {
		return Rejected[T](errors.New("gosync: Race called with no promises"))
	}

	return AsyncWithContext(context.Background(), func(ctx context.Context) (T, error) {
		type raceRes struct {
			val T
			err error
		}
		out := make(chan raceRes, 1)

		for _, p := range promises {
			go func(prom *Promise[T]) {
				v, err := prom.Await(ctx)
				select {
				case out <- raceRes{val: v, err: err}:
				default:
				}
			}(p)
		}

		select {
		case <-ctx.Done():
			var zero T
			return zero, ctx.Err()
		case r := <-out:
			return r.val, r.err
		}
	})
}

// Any returns a Promise that resolves as soon as any of the promises fulfills.
// If all promises reject, it rejects with an *AggregateError containing all errors.
func Any[T any](promises ...*Promise[T]) *Promise[T] {
	if len(promises) == 0 {
		return Rejected[T](&AggregateError{Errors: []error{errors.New("gosync: Any called with no promises")}})
	}

	return AsyncWithContext(context.Background(), func(ctx context.Context) (T, error) {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		type successRes struct {
			val T
		}
		successChan := make(chan successRes, 1)
		errs := make([]error, len(promises))
		var wg sync.WaitGroup
		wg.Add(len(promises))

		for i, p := range promises {
			go func(idx int, prom *Promise[T]) {
				defer wg.Done()
				v, err := prom.Await(ctx)
				if err != nil {
					errs[idx] = err
					return
				}
				select {
				case successChan <- successRes{val: v}:
					cancel()
				default:
				}
			}(i, p)
		}

		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-ctx.Done():
			select {
			case s := <-successChan:
				return s.val, nil
			default:
				return *new(T), ctx.Err()
			}
		case s := <-successChan:
			return s.val, nil
		case <-done:
			select {
			case s := <-successChan:
				return s.val, nil
			default:
				return *new(T), &AggregateError{Errors: errs}
			}
		}
	})
}

// AllSettled returns a Promise that resolves after all given promises have settled (either fulfilled or rejected).
// It resolves to a slice of SettledResult preserving the input order.
func AllSettled[T any](promises ...*Promise[T]) *Promise[[]SettledResult[T]] {
	if len(promises) == 0 {
		return Resolved([]SettledResult[T]{})
	}

	return Async(func() ([]SettledResult[T], error) {
		results := make([]SettledResult[T], len(promises))
		var wg sync.WaitGroup
		wg.Add(len(promises))

		for i, p := range promises {
			go func(idx int, prom *Promise[T]) {
				defer wg.Done()
				v, err := prom.Await(context.Background())
				if err != nil {
					results[idx] = SettledResult[T]{
						State: StateRejected,
						Err:   err,
					}
				} else {
					results[idx] = SettledResult[T]{
						State: StateFulfilled,
						Value: v,
					}
				}
			}(i, p)
		}

		wg.Wait()
		return results, nil
	})
}
