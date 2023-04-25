package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"sync"
	"time"
)

const (
	numRequestsPerSecond = 10
	numRequests          = 100
	totalRequests        = numRequestsPerSecond * 60
	url                  = "https://api.irvyn.xyz/bugreport"
)

func main() {
	start := time.Now()
	var sentRequests int
	var wg sync.WaitGroup
	const numRequestsPerBatch = 10
	const batchInterval = time.Second

	// Store response times
	responseTimes := make([]float64, numRequests)

	var mu sync.Mutex // Mutex to protect sentRequests increment

	for i := 0; i < numRequests; i += numRequestsPerBatch {
		wg.Add(numRequestsPerBatch)

		for j := i; j < i+numRequestsPerBatch && j < numRequests; j++ {
			go func(j int) {
				defer wg.Done()
				title := RandStringRunes(50)
				bugReport := RandStringRunes(1024)
				data := map[string]string{"title": title, "bugReport": bugReport}
				var buf bytes.Buffer
				enc := json.NewEncoder(&buf)
				enc.Encode(data)
				reqStart := time.Now()
				resp, err := http.Post(url, "application/json", &buf)
				if err != nil {
					fmt.Println(err)
				} else {
					responseTimes[j] = time.Since(reqStart).Seconds() * 1000
					if resp.StatusCode == http.StatusOK {
						mu.Lock()
						sentRequests++
						mu.Unlock()
					}
				}
			}(j)
		}

		wg.Wait()
		time.Sleep(batchInterval)
	}

	elapsed := time.Since(start)
	fmt.Printf("Time taken %s\n", elapsed)
	fmt.Printf("Sent %d successful requests in %.2fs\n", sentRequests, elapsed.Seconds())
	fmt.Printf("%.2f successful requests per second\n", float64(sentRequests)/elapsed.Seconds())

	// Calculate P99, P95, P75, and P50 response times
	sort.Float64s(responseTimes[:sentRequests])
	p99Index := int(0.99 * float64(sentRequests))
	p95Index := int(0.95 * float64(sentRequests))
	p75Index := int(0.75 * float64(sentRequests))
	p50Index := int(0.50 * float64(sentRequests))
	fmt.Printf("P99 response time: %.2fms\n", responseTimes[p99Index])
	fmt.Printf("P95 response time: %.2fms\n", responseTimes[p95Index])
	fmt.Printf("P75 response time: %.2fms\n", responseTimes[p75Index])
	fmt.Printf("P50 response time: %.2fms\n", responseTimes[p50Index])

}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
