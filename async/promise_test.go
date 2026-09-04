package async 

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"strconv"
	"time"
)

func TestAsync_Success(t *testing.T) {
	t.Parallel()

	expectedVal := "hello async"
	p := Async(func() (string, error) {
		return expectedVal, nil
	})

	val, err := p.Await(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if val != expectedVal {
		t.Fatalf("expected %q, got: %q", expectedVal, val)
	}
}

func TestAsync_Error(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("something went wrong")
	p := Async(func() (int, error) {
		return 0, expectedErr
	})

	val, err := p.Await(context.Background())
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got: %v", expectedErr, err)
	}
	if val != 0 {
		t.Fatalf("expected zero value 0, got: %d", val)
	}
}

func TestPromise_Await_ContextCanceled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())

	p := Async(func() (int, error) {
		time.Sleep(200 * time.Millisecond)
		return 42, nil
	})

	// Cancel context immediately
	cancel()

	val, err := p.Await(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got: %v", err)
	}
	if val != 0 {
		t.Fatalf("expected zero value 0, got: %d", val)
	}
}

func TestPromise_Await_ContextTimeout(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	p := Async(func() (string, error) {
		time.Sleep(200 * time.Millisecond)
		return "done", nil
	})

	val, err := p.Await(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context.DeadlineExceeded, got: %v", err)
	}
	if val != "" {
		t.Fatalf("expected empty string, got: %q", val)
	}
}

func TestPromise_Then_Success(t *testing.T) {
	t.Parallel()

	p := Async(func() (int, error) {
		return 10, nil
	}).Then(func(v int) (int, error) {
		return v * 2, nil
	}).Then(func(v int) (int, error) {
		return v + 5, nil
	})

	val, err := p.Await(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if val != 25 {
		t.Fatalf("expected 25, got: %d", val)
	}
}

func TestPromise_Then_Generic_Types(t *testing.T){
	t.Parallel()

	// convertion: string -> int -> string
	val, err := Async(func() (string, error) {
		time.Sleep(100 * time.Millisecond)
		return "42", nil
	}).Then(func(s string) (int, error) {
		return strconv.Atoi(s)
	}).Then(func(n int) (string, error) {
		return fmt.Sprintf("Double 42 is %d", n*2), nil
	}).Await(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	expected := "Double 42 is 84" 

	if val !=  "Double 42 is 84"{
		t.Fatalf("expected '%s', got: '%s'",expected, val)
	}
	
}

func TestPromise_Then_InitialErrorPropagation(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("initial failure")
	var thenExecuted atomic.Bool

	p := Async(func() (string, error) {
		return "", expectedErr
	}).Then(func(v string) (string, error) {
		thenExecuted.Store(true)
		return v + " appended", nil
	})

	val, err := p.Await(context.Background())
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got: %v", expectedErr, err)
	}
	if val != "" {
		t.Fatalf("expected empty string, got: %q", val)
	}
	if thenExecuted.Load() {
		t.Fatalf("expected Then callback to not be executed after initial error")
	}
}

func TestPromise_Then_CallbackError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("error in then callback")

	p := Async(func() (int, error) {
		return 100, nil
	}).Then(func(v int) (int, error) {
		return 0, expectedErr
	})

	val, err := p.Await(context.Background())
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got: %v", expectedErr, err)
	}
	if val != 0 {
		t.Fatalf("expected 0, got: %d", val)
	}
}

func TestAsync_NonBlockingExecution(t *testing.T) {
	t.Parallel()

	started := make(chan struct{})
	done := make(chan struct{})

	p := Async(func() (bool, error) {
		close(started)
		<-done
		return true, nil
	})

	// Wait until goroutine has started
	<-started

	// Promise should not be resolved yet
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	_, err := p.Await(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded while task is running, got: %v", err)
	}

	// Release task and await completion
	close(done)

	val, err := p.Await(context.Background())
	if err != nil {
		t.Fatalf("expected no error after task completion, got: %v", err)
	}
	if !val {
		t.Fatalf("expected true, got: %v", val)
	}
}

func TestAsync_MultipleConcurrentPromises(t *testing.T) {
	t.Parallel()

	count := 20
	promises := make([]*Promise[int], count)

	for i := range count {
		val := i
		promises[i] = Async(func() (int, error) {
			time.Sleep(10 * time.Millisecond)
			return val * 2, nil
		})
	}

	for i := range count {
		res, err := promises[i].Await(context.Background())
		if err != nil {
			t.Fatalf("promise %d returned unexpected error: %v", i, err)
		}
		if res != i*2 {
			t.Fatalf("promise %d: expected %d, got %d", i, i*2, res)
		}
	}
}
