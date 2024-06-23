package background

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/iamsorryprincess/project-layout/internal/pkg/log"
)

var ErrWorkerStopped = errors.New("worker stopped")

type WorkerFunc func(ctx context.Context) error

type Worker struct {
	wg        sync.WaitGroup
	mu        sync.Mutex
	done      chan struct{}
	isStopped bool

	logger log.Logger
}

func NewWorker(logger log.Logger) *Worker {
	return &Worker{
		wg:        sync.WaitGroup{},
		mu:        sync.Mutex{},
		done:      make(chan struct{}),
		isStopped: false,
		logger:    logger,
	}
}

func (w *Worker) Start(ctx context.Context, name string, fn WorkerFunc) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.isStopped {
		return ErrWorkerStopped
	}

	w.wg.Add(1)

	go func() {
		defer w.wg.Done()
		for {
			select {
			case <-w.done:
				return
			default:
				now := time.Now()
				w.runWorkerFunc(ctx, name, fn)
				w.logger.Debug().
					Str("type", "worker").
					Str("worker", name).
					Msgf("worker %s finished in %s", name, time.Since(now))
			}
		}
	}()

	return nil
}

func (w *Worker) StartWithInterval(ctx context.Context, name string, fn WorkerFunc, interval time.Duration) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.isStopped {
		return ErrWorkerStopped
	}

	w.wg.Add(1)

	go func() {
		defer w.wg.Done()
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-w.done:
				return
			case <-ticker.C:
				now := time.Now()
				w.runWorkerFunc(ctx, name, fn)
				w.logger.Debug().
					Str("type", "worker").
					Str("worker", name).
					Msgf("worker %s finished in %s", name, time.Since(now))
			}
		}
	}()

	return nil
}

func (w *Worker) Run(ctx context.Context, name string, fn WorkerFunc) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.isStopped {
		return ErrWorkerStopped
	}

	w.wg.Add(1)

	go func() {
		defer w.wg.Done()
		now := time.Now()
		w.runWorkerFunc(ctx, name, fn)
		w.logger.Debug().
			Str("type", "worker").
			Str("worker", name).
			Msgf("worker %s finished in %s", name, time.Since(now))
	}()

	return nil
}

func (w *Worker) StopAll() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.isStopped = true
	w.done <- struct{}{}
	w.wg.Wait()
}

func (w *Worker) runWorkerFunc(ctx context.Context, name string, fn WorkerFunc) {
	defer w.recovery(name)
	if err := fn(ctx); err != nil {
		w.logger.Error().
			Str("type", "worker").
			Str("worker", name).
			Msgf("worker failed with err: %v", err)
	}
}

func (w *Worker) recovery(name string) {
	if r := recover(); r != nil {
		var err error
		switch t := r.(type) {
		case string:
			err = errors.New(t)
		case error:
			err = t
		default:
			err = fmt.Errorf("unknown error: %v", t)
		}

		w.logger.Error().
			Str("type", "worker").
			Str("worker", name).
			Msgf("worker recovered from panic: %v", err)
	}
}
