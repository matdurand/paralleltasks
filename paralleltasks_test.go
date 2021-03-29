package paralleltasks_test

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/matdurand/paralleltasks"
)

func readChannel(ch chan string) []string {
	vals := []string{}
	for s := range ch {
		vals = append(vals, s)
	}
	return vals
}

var _ = Describe("Parallel tasks", func() {
	When("Using a timeout context", func() {
		When("The context timeout before the execution is complete", func() {
			It("Should return a timeout error", func() {
				timeout := time.Millisecond * 1
				ctx, _ := context.WithTimeout(context.Background(), timeout)

				pt := paralleltasks.New(ctx, 1)
				respChan := make(chan string, 1)
				pt.Run(func(ctx context.Context, errChan chan error) {
					time.Sleep(time.Millisecond * 10)
					respChan <- "result1"
				})

				err := pt.Wait()
				Expect(err).To(Equal(context.DeadlineExceeded))
			})
		})

		When("The task completes before the timeout expires", func() {
			It("Should return the task result", func() {
				timeout := time.Second * 1
				ctx, _ := context.WithTimeout(context.Background(), timeout)

				pt := paralleltasks.New(ctx, 1)
				respChan := make(chan string, 1)
				pt.Run(func(ctx context.Context, errChan chan error) {
					respChan <- "result1"
				})

				err := pt.Wait()
				Expect(err).To(BeNil())

				close(respChan)
				Expect(readChannel(respChan)).To(Equal([]string{"result1"}))
			})
		})
	})

	When("One task returns an error", func() {
		It("Should return the error", func() {
			ctx := context.Background()
			pt := paralleltasks.New(ctx, 1)
			pt.Run(func(ctx context.Context, errChan chan error) {
				errChan <- errors.New("error1")
			})

			err := pt.Wait()
			Expect(err).To(Equal(errors.New("error1")))
		})
	})

	When("All tasks return a results", func() {
		It("Should return one result per task", func() {
			ctx := context.Background()
			pt := paralleltasks.New(ctx, 2)
			respChan := make(chan string, 2)
			pt.Run(func(ctx context.Context, errChan chan error) {
				respChan <- "result1"
			})
			pt.Run(func(ctx context.Context, errChan chan error) {
				time.Sleep(time.Millisecond * 100)
				respChan <- "result2"
			})

			err := pt.Wait()
			Expect(err).To(BeNil())

			close(respChan)
			Expect(readChannel(respChan)).To(Equal([]string{"result1", "result2"}))
		})
	})
})
