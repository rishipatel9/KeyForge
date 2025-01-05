package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var addr = flag.String("addr", "localhost:5000", "The HTTP host port for the instance that is benchmarked")

const N = 100

func writeRand() {
	key := fmt.Sprintf("key-%d", rand.Intn(100000))
	value := fmt.Sprintf("value-%d", rand.Intn(100000))

	values := url.Values{}
	values.Set("key", key)
	values.Set("value", value)
	resp, err := http.Get("http://" + *addr + "/set?key=" + values.Encode())

	if err != nil {
		log.Fatalf("error occured during set : %v", err)
	}
	defer resp.Body.Close()

	// fmt.Printf("%s - %s \n", key, value)
}

func benchmark(name string, fn func()) {
	start := time.Now()
	var max time.Duration
	var min time.Duration = time.Duration(1<<63 - 1)

	for i := 0; i < N; i++ {
		iterStart := time.Now()
		fn()
		iterTime := time.Since(iterStart)
		if iterTime > max {
			max = iterTime
		} else if iterTime < min {
			min = iterTime
		}
	}

	qps := float64(N) / (float64(time.Since(start)) / float64(time.Second))
	fmt.Printf("Func %s took %s time on a average, max : %s , min : %s , qps : %.1f \n", name, time.Since(start)/N, max, min, qps)

}
func main() {

	flag.Parse()

	const concurrency = 4
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			benchmark("write", writeRand)
			wg.Done()
		}()
	}
	wg.Wait()
}
