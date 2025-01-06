package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var addr = flag.String("addr", "localhost:5000", "The HTTP host port for the instance that is benchmarked")

const N = 1000
const concurrency = 16

func writeRand() (key string) {
	key = fmt.Sprintf("key-%d", rand.Intn(100000))
	value := fmt.Sprintf("value-%d", rand.Intn(100000))

	values := url.Values{}
	values.Set("key", key)
	values.Set("value", value)
	resp, err := http.Get("http://" + *addr + "/set?key=" + values.Encode())

	if err != nil {
		log.Fatalf("error occured during set : %v", err)
	}
	io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()
	return key
	// fmt.Printf("%s - %s \n", key, value)
}
func benchmarkWrite() (allkeys []string) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var totalQps float64

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			qps, strs := benchmark("write", func() string { return writeRand() })
			mu.Lock()
			totalQps += qps
			allkeys = append(allkeys, strs...)
			mu.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	log.Printf("Total Read Qps : %1.f", totalQps)
	return allkeys
}
func benchmarkRead(allkeys []string) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var totalQps float64

	wg.Wait()
	log.Printf("Total Read Qps : %1.f", totalQps)

	totalQps = 0
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			qps, _ := benchmark("read", func() string { return readRand(allkeys) })
			mu.Lock()
			totalQps += qps
			mu.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	log.Printf("Total Write Qps : %1.f", totalQps)
}
func benchmark(name string, fn func() string) (qps float64, strs []string) {
	start := time.Now()
	var max time.Duration
	var min time.Duration = time.Duration(1<<63 - 1)

	for i := 0; i < N; i++ {
		iterStart := time.Now()
		strs = append(strs, fn())
		iterTime := time.Since(iterStart)
		if iterTime > max {
			max = iterTime
		} else if iterTime < min {
			min = iterTime
		}
	}

	qps = float64(N) / (float64(time.Since(start)) / float64(time.Second))
	fmt.Printf("Func %s took %s time on a average, max : %s , min : %s , qps : %.1f \n", name, time.Since(start)/N, max, min, qps)

	return qps, strs
}

func readRand(allkeys []string) (key string) {
	key = allkeys[rand.Intn(len(allkeys))]

	values := url.Values{}
	values.Set("key", key)
	resp, err := http.Get("http://" + *addr + "/get?" + values.Encode())

	if err != nil {
		log.Fatalf("error occured during set : %v", err)
	}
	defer resp.Body.Close()
	return key
}
func main() {

	flag.Parse()

	allkeys := benchmarkWrite()
	benchmarkRead(allkeys)
}
