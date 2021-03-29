# paralleltasks

A library for a parallel tasks execution abstract in Golang.

This library provides a simple abstraction to be able to execute a bunch of task in parallel like a standard `sync.WaitGroup` would do, but it supports using a context for timeout and cancellation, and it returns the tasks' errors (the first one to come up).

## Example

```go

func main() {

	timeout := time.Second * 10
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	urls := []string{
		"https://jsonplaceholder.typicode.com/todos/1",
		"https://jsonplaceholder.typicode.com/todos/2",
		"https://jsonplaceholder.typicode.com/todos/3",
	}

	respChan := make(chan string, len(urls))
	pt := paralleltasks.New(ctx, len(urls))
	for _, url := range urls {
		url := url
		pt.Run(func(ctx context.Context, errChan chan error) {
			httpGet(url, ctx, respChan, errChan)
		})
	}

	err := pt.Wait()
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	close(respChan)
	for s := range respChan {
		fmt.Println(s)
	}
}

func httpGet(url string, ctx context.Context, respChan chan interface{}, errChan chan error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		errChan <- err
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		errChan <- err
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errChan <- err
		return
	}

	respChan <- string(body)
}
```

This will print:
```json
{
  "userId": 1,
  "id": 2,
  "title": "quis ut nam facilis et officia qui",
  "completed": false
} 
{
  "userId": 1,
  "id": 1,
  "title": "delectus aut autem",
  "completed": false
} 
{
  "userId": 1,
  "id": 3,
  "title": "fugiat veniam minus",
  "completed": false
}
```

## Caveats

When using a timeout context, the waitgroup is going to exit and the `wait` method is going to return, even if all the go routines are not done yet. This could lead to a leak if you go routines never returns. Be sure to use the provided context in the task arguments to detect cancellation and terminate your go routine. You should also call `cancel` on your context once the `wait` method has returned to ensure that every task is cancelled.

The  `wait` method will exit if any of the task send an error in the errors channel. You will receive the first error only, even if multiple go routines push an error in the error channel.