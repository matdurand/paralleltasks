package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/matdurand/paralleltasks"
)

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

func httpGet(url string, ctx context.Context, respChan chan string, errChan chan error) {
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

	fmt.Println("Writing body", url)
	respChan <- string(body)
}
