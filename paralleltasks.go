package paralleltasks

import (
	"context"
	"sync"
)

type parallelTasks struct {
	wg      *sync.WaitGroup
	ctx     context.Context
	waitCh  chan struct{}
	errChan chan error
}

func New(ctx context.Context, taskCount int) parallelTasks {

	wg := sync.WaitGroup{}
	wg.Add(taskCount)

	waitCh := make(chan struct{})
	errChan := make(chan error, taskCount)

	return parallelTasks{
		wg:      &wg,
		ctx:     ctx,
		waitCh:  waitCh,
		errChan: errChan,
	}
}

func (pt *parallelTasks) Run(task func(ctx context.Context, errChan chan error)) {
	go func() {
		defer pt.wg.Done()
		task(pt.ctx, pt.errChan)
	}()
}

func (pt *parallelTasks) Wait() error {
	go func() {
		pt.wg.Wait()
		close(pt.waitCh)
	}()

	select {
	case err := <-pt.errChan:
		return err
	case <-(pt.ctx).Done():
		return (pt.ctx).Err()
	case <-pt.waitCh:
		return nil
	}
}
