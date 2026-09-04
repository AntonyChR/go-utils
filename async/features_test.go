package async

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestPromise_MultipleAwait_Sequential(t *testing.T) {
	t.Parallel()

	expectedVal := 123
	p := Async(func() (int, error) {
		time.Sleep(10 * time.Millisecond)
		return expectedVal, nil
	})

	for i := range 5 {
		val, err := p.Await(context.Background())
		if err != nil {
			t.Fatalf("call %d: expected no error, got %v", i, err)
		}
		if val != expectedVal {
			t.Fatalf("call %d: expected %d, got %d", i, expectedVal, val)
		}
	}
}

func TestPromise_MultipleAwait_Concurrent(t *testing.T) {
	t.Parallel()

	expectedVal := "concurrent-test"
	p := Async(func() (string, error) {
		time.Sleep(20 * time.Millisecond)
		return expectedVal, nil
	})

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := range goroutines {
		go func(idx int) {
			defer wg.Done()
			val, err := p.Await(context.Background())
			if err != nil {
				t.Errorf("goroutine %d: unexpected error: %v", idx, err)
				return
			}
			if val != expectedVal {
				t.Errorf("goroutine %d: expected %q, got %q", idx, expectedVal, val)
			}
		}(i)
	}

	wg.Wait()
}

func TestPromise_PanicRecovery(t *testing.T) {
	t.Parallel()

	panicMessage := "something exploded fatally"
	p := Async(func() (int, error) {
		panic(panicMessage)
	})

	val, err := p.Await(context.Background())
	if err == nil {
		t.Fatalf("expected error from panic, got nil")
	}
	if val != 0 {
		t.Fatalf("expected zero value, got %d", val)
	}

	var panicErr *PanicError
	if !errors.As(err, &panicErr) {
		t.Fatalf("expected error to be *PanicError, got %T: %v", err, err)
	}
	if panicErr.Value != panicMessage {
		t.Fatalf("expected panic value %q, got %v", panicMessage, panicErr.Value)
	}
	if len(panicErr.Stack) == 0 {
		t.Fatalf("expected non-empty stack trace")
	}
}

func TestPromise_StateInspection(t *testing.T) {
	t.Parallel()

	blocker := make(chan struct{})
	p := Async(func() (string, error) {
		<-blocker
		return "ok", nil
	})

	// Check pending state
	if p.State() != StatePending {
		t.Fatalf("expected StatePending, got %v", p.State())
	}
	if p.IsSettled() {
		t.Fatalf("expected IsSettled=false")
	}
	if _, _, ok := p.Result(); ok {
		t.Fatalf("expected Result() ok=false while pending")
	}

	close(blocker)
	_, _ = p.Await(context.Background())

	// Check fulfilled state
	if p.State() != StateFulfilled {
		t.Fatalf("expected StateFulfilled, got %v", p.State())
	}
	if !p.IsSettled() || !p.IsFulfilled() || p.IsRejected() {
		t.Fatalf("state helpers mismatch for fulfilled promise")
	}
	val, err, ok := p.Result()
	if !ok || err != nil || val != "ok" {
		t.Fatalf("unexpected Result(): val=%q, err=%v, ok=%v", val, err, ok)
	}
}

func TestPromise_Resolved_Rejected(t *testing.T) {
	t.Parallel()

	p1 := Resolved(999)
	if !p1.IsFulfilled() {
		t.Fatalf("expected p1 to be fulfilled")
	}
	v1, err1 := p1.Await(context.Background())
	if err1 != nil || v1 != 999 {
		t.Fatalf("expected 999 and nil error, got %d, %v", v1, err1)
	}

	testErr := errors.New("static rejection")
	p2 := Rejected[int](testErr)
	if !p2.IsRejected() {
		t.Fatalf("expected p2 to be rejected")
	}
	_, err2 := p2.Await(context.Background())
	if !errors.Is(err2, testErr) {
		t.Fatalf("expected error %v, got %v", testErr, err2)
	}
}

func TestPromise_Catch_Recovery(t *testing.T) {
	t.Parallel()

	p := Async(func() (int, error) {
		return 0, errors.New("initial failure")
	}).Catch(func(err error) (int, error) {
		// Recover with fallback value 42
		return 42, nil
	})

	val, err := p.Await(context.Background())
	if err != nil {
		t.Fatalf("expected error to be caught, got %v", err)
	}
	if val != 42 {
		t.Fatalf("expected recovered value 42, got %d", val)
	}
}

func TestPromise_Catch_Passthrough(t *testing.T) {
	t.Parallel()

	var catchExecuted atomic.Bool
	p := Resolved(100).Catch(func(err error) (int, error) {
		catchExecuted.Store(true)
		return 999, nil
	})

	val, err := p.Await(context.Background())
	if err != nil || val != 100 {
		t.Fatalf("expected 100, got %d, %v", val, err)
	}
	if catchExecuted.Load() {
		t.Fatalf("catch callback should not run on fulfilled promise")
	}
}

func TestPromise_Finally(t *testing.T) {
	t.Parallel()

	var finallyRan atomic.Bool
	p := Async(func() (string, error) {
		return "hello", nil
	}).Finally(func() {
		finallyRan.Store(true)
	})

	val, err := p.Await(context.Background())
	if err != nil || val != "hello" {
		t.Fatalf("expected hello, got %s, %v", val, err)
	}
	if !finallyRan.Load() {
		t.Fatalf("finally callback did not run")
	}
}

func TestPromise_WithTimeout(t *testing.T) {
	t.Parallel()

	pSlow := Async(func() (string, error) {
		time.Sleep(100 * time.Millisecond)
		return "slow", nil
	}).WithTimeout(20 * time.Millisecond)

	_, errSlow := pSlow.Await(context.Background())
	if !errors.Is(errSlow, context.DeadlineExceeded) {
		t.Fatalf("expected DeadlineExceeded, got %v", errSlow)
	}

	pFast := Async(func() (string, error) {
		return "fast", nil
	}).WithTimeout(200 * time.Millisecond)

	valFast, errFast := pFast.Await(context.Background())
	if errFast != nil || valFast != "fast" {
		t.Fatalf("expected fast, got %s, %v", valFast, errFast)
	}
}

func TestCombinators_All_Success(t *testing.T) {
	t.Parallel()

	p1 := Async(func() (int, error) {
		time.Sleep(20 * time.Millisecond)
		return 1, nil
	})
	p2 := Async(func() (int, error) {
		time.Sleep(10 * time.Millisecond)
		return 2, nil
	})
	p3 := Resolved(3)

	all := All(p1, p2, p3)
	vals, err := all.Await(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(vals) != 3 || vals[0] != 1 || vals[1] != 2 || vals[2] != 3 {
		t.Fatalf("unexpected results: %v", vals)
	}
}

func TestCombinators_All_FailFast(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("boom")
	p1 := Async(func() (int, error) {
		time.Sleep(500 * time.Millisecond)
		return 1, nil
	})
	p2 := Async(func() (int, error) {
		return 0, expectedErr
	})

	start := time.Now()
	all := All(p1, p2)
	_, err := all.Await(context.Background())
	elapsed := time.Since(start)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected boom error, got %v", err)
	}
	if elapsed > 200*time.Millisecond {
		t.Fatalf("expected fail-fast in <200ms, took %v", elapsed)
	}
}

func TestCombinators_Race(t *testing.T) {
	t.Parallel()

	pSlow := Async(func() (string, error) {
		time.Sleep(150 * time.Millisecond)
		return "slow", nil
	})
	pFast := Async(func() (string, error) {
		time.Sleep(10 * time.Millisecond)
		return "fast", nil
	})

	race := Race(pSlow, pFast)
	val, err := race.Await(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if val != "fast" {
		t.Fatalf("expected fast to win race, got %s", val)
	}
}

func TestCombinators_Any(t *testing.T) {
	t.Parallel()

	err1 := errors.New("err1")
	err2 := errors.New("err2")

	pErr1 := Async(func() (string, error) {
		time.Sleep(10 * time.Millisecond)
		return "", err1
	})
	pSucc := Async(func() (string, error) {
		time.Sleep(30 * time.Millisecond)
		return "success", nil
	})
	pErr2 := Async(func() (string, error) {
		time.Sleep(10 * time.Millisecond)
		return "", err2
	})

	anyP := Any(pErr1, pSucc, pErr2)
	val, err := anyP.Await(context.Background())
	if err != nil {
		t.Fatalf("expected Any to succeed, got %v", err)
	}
	if val != "success" {
		t.Fatalf("expected success, got %s", val)
	}

	// Test all failing -> AggregateError
	allFail := Any(pErr1, pErr2)
	_, errAll := allFail.Await(context.Background())
	var aggErr *AggregateError
	if !errors.As(errAll, &aggErr) {
		t.Fatalf("expected *AggregateError, got %T: %v", errAll, errAll)
	}
	if len(aggErr.Errors) != 2 {
		t.Fatalf("expected 2 errors in AggregateError, got %d", len(aggErr.Errors))
	}
}

func TestCombinators_AllSettled(t *testing.T) {
	t.Parallel()

	errTest := errors.New("failure")
	p1 := Resolved("a")
	p2 := Rejected[string](errTest)

	settled, err := AllSettled(p1, p2).Await(context.Background())
	if err != nil {
		t.Fatalf("expected AllSettled to never error, got %v", err)
	}
	if len(settled) != 2 {
		t.Fatalf("expected 2 results, got %d", len(settled))
	}
	if settled[0].State != StateFulfilled || settled[0].Value != "a" {
		t.Fatalf("unexpected settled[0]: %+v", settled[0])
	}
	if settled[1].State != StateRejected || !errors.Is(settled[1].Err, errTest) {
		t.Fatalf("unexpected settled[1]: %+v", settled[1])
	}
}

func TestUtils_AsyncWithRetry(t *testing.T) {
	t.Parallel()

	var attempts atomic.Int32
	p := AsyncWithRetry(3, 10*time.Millisecond, func() (string, error) {
		count := attempts.Add(1)
		if count < 3 {
			return "", fmt.Errorf("fail attempt %d", count)
		}
		return "retry success", nil
	})

	val, err := p.Await(context.Background())
	if err != nil {
		t.Fatalf("expected retry to succeed on 3rd attempt, got %v", err)
	}
	if val != "retry success" {
		t.Fatalf("expected retry success, got %s", val)
	}
	if attempts.Load() != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts.Load())
	}
}

func TestUtils_MapConcurrent(t *testing.T) {
	t.Parallel()

	items := []int{1, 2, 3, 4, 5, 6, 7, 8}
	var currentActive atomic.Int32
	var maxActive atomic.Int32

	p := MapConcurrent(context.Background(), items, 3, func(ctx context.Context, item int) (int, error) {
		curr := currentActive.Add(1)
		// Update max concurrency recorded
		for {
			oldMax := maxActive.Load()
			if curr <= oldMax || maxActive.CompareAndSwap(oldMax, curr) {
				break
			}
		}

		time.Sleep(15 * time.Millisecond)
		currentActive.Add(-1)
		return item * 10, nil
	})

	results, err := p.Await(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(results) != len(items) {
		t.Fatalf("expected %d results, got %d", len(items), len(results))
	}
	for i, v := range results {
		if v != items[i]*10 {
			t.Fatalf("expected results[%d] = %d, got %d", i, items[i]*10, v)
		}
	}
	if maxActive.Load() > 3 {
		t.Fatalf("concurrency limit exceeded: max active was %d (limit was 3)", maxActive.Load())
	}
}
