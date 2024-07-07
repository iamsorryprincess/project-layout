package background

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/iamsorryprincess/project-layout/internal/pkg/log"
)

var (
	ErrWorkerStopped  = errors.New("worker stopped")
	ErrWorkerNotFound = errors.New("worker not found")
)

type WorkerFunc func(ctx context.Context) error

type workerItem struct {
	ID        string
	Name      string
	CloseChan chan struct{}
}

type Worker struct {
	wg        sync.WaitGroup
	mu        sync.Mutex
	isStopped bool
	items     map[string]workerItem

	logger log.Logger
}

func NewWorker(logger log.Logger) *Worker {
	return &Worker{
		wg:        sync.WaitGroup{},
		mu:        sync.Mutex{},
		isStopped: false,
		items:     make(map[string]workerItem),
		logger:    newLogger(logger),
	}
}

func (w *Worker) Start(ctx context.Context, name string, fn WorkerFunc) (string, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.isStopped {
		return "", ErrWorkerStopped
	}

	w.wg.Add(1)
	id := uuid.New().String()
	closeChan := make(chan struct{})

	w.items[id] = workerItem{
		ID:        id,
		Name:      name,
		CloseChan: closeChan,
	}

	go func() {
		defer w.wg.Done()
		for {
			select {
			case <-closeChan:
				w.logger.Debug().Str("worker", name).Msg("worker stopped")
				return
			default:
				now := time.Now()
				w.runWorkerFunc(ctx, name, fn)
				w.logger.Debug().Str("worker", name).Msgf("worker %s finished in %s", name, time.Since(now))
			}
		}
	}()

	return id, nil
}

func (w *Worker) StartWithInterval(ctx context.Context, name string, interval time.Duration, fn WorkerFunc) (string, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.isStopped {
		return "", ErrWorkerStopped
	}

	w.wg.Add(1)
	id := uuid.New().String()
	closeChan := make(chan struct{})

	w.items[id] = workerItem{
		ID:        id,
		Name:      name,
		CloseChan: closeChan,
	}

	go func() {
		defer w.wg.Done()

		now := time.Now()
		w.runWorkerFunc(ctx, name, fn)
		w.logger.Debug().Str("worker", name).Msgf("worker %s finished in %s", name, time.Since(now))

		timer := time.NewTimer(interval)
		defer timer.Stop()

		for {
			select {
			case <-closeChan:
				w.logger.Debug().Str("worker", name).Msg("worker stopped")
				return
			case <-timer.C:
				now = time.Now()
				w.runWorkerFunc(ctx, name, fn)
				w.logger.Debug().Str("worker", name).Msgf("worker %s finished in %s", name, time.Since(now))
				timer.Reset(interval)
			}
		}
	}()

	return id, nil
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
		w.logger.Debug().Str("worker", name).Msgf("worker %s finished in %s", name, time.Since(now))
	}()

	return nil
}

func (w *Worker) Stop(id string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	item, ok := w.items[id]
	if !ok {
		return ErrWorkerNotFound
	}

	item.CloseChan <- struct{}{}
	delete(w.items, id)
	return nil
}

func (w *Worker) StopAll() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.isStopped = true

	for _, item := range w.items {
		close(item.CloseChan)
	}

	clear(w.items)
	w.wg.Wait()
}

func (w *Worker) runWorkerFunc(ctx context.Context, name string, fn WorkerFunc) {
	defer w.recovery(name)
	if err := fn(ctx); err != nil {
		w.logger.Error().Str("worker", name).Msgf("worker failed with err: %v", err)
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

		w.logger.Error().Str("worker", name).Msgf("worker recovered from panic: %v", err)
	}
}
