package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

func main() {
	makeRequest := func(routineNumber int) {
		resp, err := http.Get(`http://localhost:8080/counter`)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		fmt.Printf("routineNum: %d, status: %d, body: %s", routineNumber, resp.StatusCode, body)
	}

	var wg sync.WaitGroup

	for i := 0; i < 7; i++ {
		wg.Add(1)
		go func(routineNumber int) {
			defer wg.Done()
			makeRequest(routineNumber)
		}(i)
	}

	wg.Wait()
}
