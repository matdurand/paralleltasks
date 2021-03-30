package paralleltasks

import (
	"context"
	"sync"
)

type tasks struct {
	wg      *sync.WaitGroup
	ctx     context.Context
	waitCh  chan struct{}
	errChan chan error
}
type Tasks interface {
	Run(task func(ctx context.Context, errChan chan error))
	Wait() error
}

func (pt tasks) Run(task func(ctx context.Context, errChan chan error)) {
	go func() {
		defer pt.wg.Done()
		task(pt.ctx, pt.errChan)
	}()
}

func (pt tasks) Wait() error {
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

func New(ctx context.Context, taskCount int) Tasks {

	wg := sync.WaitGroup{}
	wg.Add(taskCount)

	waitCh := make(chan struct{})
	errChan := make(chan error, taskCount)

	return tasks{
		wg:      &wg,
		ctx:     ctx,
		waitCh:  waitCh,
		errChan: errChan,
	}
}
