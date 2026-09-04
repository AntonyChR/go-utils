package async

import (
	"context"
	"sync"
	"time"
)

// AsyncWithRetry attempts to execute fn up to attempts times before returning the final error.
// Between attempts, it waits for the specified delay.
func AsyncWithRetry[T any](attempts int, delay time.Duration, fn func() (T, error)) *Promise[T] {
	return Async(func() (T, error) {
		if attempts <= 0 {
			attempts = 1
		}
		var lastErr error
		var zero T
		for i := 0; i < attempts; i++ {
			res, err := fn()
			if err == nil {
				return res, nil
			}
			lastErr = err
			if i < attempts-1 && delay > 0 {
				time.Sleep(delay)
			}
		}
		return zero, lastErr
	})
}

// MapConcurrent applies fn to each item in items concurrently, with at most limit goroutines executing concurrently.
// If ctx is canceled or fn returns an error, the operation aborts early.
func MapConcurrent[T, U any](ctx context.Context, items []T, limit int, fn func(context.Context, T) (U, error)) *Promise[[]U] {
	return AsyncWithContext(ctx, func(ctx context.Context) ([]U, error) {
		if len(items) == 0 {
			return []U{}, nil
		}
		if limit <= 0 {
			limit = 1
		}

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		results := make([]U, len(items))
		sem := make(chan struct{}, limit)
		errChan := make(chan error, 1)

		var wg sync.WaitGroup
		for i, item := range items {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case err := <-errChan:
				return nil, err
			case sem <- struct{}{}:
			}

			wg.Add(1)
			go func(idx int, it T) {
				defer func() {
					<-sem
					wg.Done()
				}()

				defer func() {
					if r := recover(); r != nil {
						panicErr := &PanicError{
							Value: r,
						}
						select {
						case errChan <- panicErr:
							cancel()
						default:
						}
					}
				}()

				res, err := fn(ctx, it)
				if err != nil {
					select {
					case errChan <- err:
						cancel()
					default:
					}
					return
				}
				results[idx] = res
			}(i, item)
		}

		wg.Wait()
		select {
		case err := <-errChan:
			return nil, err
		default:
			return results, nil
		}
	})
}
