package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const (
	numRequestsPerSecond = 10
	numRequests          = 100
	totalRequests        = numRequestsPerSecond * 60
	url                  = "https://api.irvyn.xyz/register"
)

func main() {
	start := time.Now()
	var sentRequests int
	var wg sync.WaitGroup
	const numRequestsPerBatch = 10
	const batchInterval = time.Second

	for i := 0; i < numRequests; i += numRequestsPerBatch {
		wg.Add(numRequestsPerBatch)

		for j := i; j < i+numRequestsPerBatch && j < numRequests; j++ {
			go func(j int) {
				defer wg.Done()
				name := RandStringRunes(20)
				password := RandStringRunes(10)
				email := name + "@example.com"
				data := map[string]string{"name": name, "password": password, "email": email}
				var buf bytes.Buffer
				enc := json.NewEncoder(&buf)
				enc.Encode(data)
				_, err := http.Post(url, "application/json", &buf)
				if err != nil {
					fmt.Println(err)
				}
				sentRequests++
			}(j)
		}

		wg.Wait()
		time.Sleep(batchInterval)
	}

	elapsed := time.Since(start)
	fmt.Printf("Time taken %s\n", elapsed)
	fmt.Printf("Sent %d requests in %.2fs\n", sentRequests, elapsed.Seconds())
	fmt.Printf("%.2f requests per second\n", float64(sentRequests)/elapsed.Seconds())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
