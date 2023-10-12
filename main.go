package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type Result struct {
	URL   string
	Data  map[string]interface{}
	Error error
}

func main() {

	urls := []string{
		"https://dummyjson.com/products/2",
		"https://dummyjson.com/products/3",
		"https://dummyjson.com/products/1",
	}

	var wg sync.WaitGroup
	resultChannel := make(chan Result, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go fetchData(url, resultChannel, &wg)
	}

	//go func() {
	wg.Wait()
	close(resultChannel)
	//}()

	results := make(map[string]Result)

	for result := range resultChannel {
		results[result.URL] = result
		if result.Error != nil {
			fmt.Printf("Error fetching data from %s: %v\n", result.URL, result.Error)
		}
	}

	for _, url := range urls {
		result, ok := results[url]
		if ok {
			fmt.Printf("Data from %s: %v\n", url, result.Data)
		}
	}

}

func fetchData(url string, resultChan chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		resultChan <- Result{URL: url, Error: err}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		resultChan <- Result{URL: url, Error: err}
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		resultChan <- Result{URL: url, Error: err}
		return
	}

	resultChan <- Result{URL: url, Data: data}
}
