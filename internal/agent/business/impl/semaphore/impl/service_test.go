package impl_test

import (
	"testing"
	"time"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/agent/business/impl/semaphore/impl"
)

func TestService_AcquireRelease(t *testing.T) {
	sem := impl.New(2)
	defer sem.Close()

	sem.Acquire()
	sem.Acquire()

	done := make(chan struct{})
	go func() {
		time.Sleep(100 * time.Millisecond)
		sem.Release()
		close(done)
	}()

	sem.Acquire()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Error("acquire didn't wait for release")
	}
}

func TestService_Release(t *testing.T) {
	sem := impl.New(1)
	defer sem.Close()

	sem.Acquire()

	sem.Release()

	done := make(chan struct{})
	go func() {
		sem.Acquire()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Error("failed to acquire after release")
	}
}
