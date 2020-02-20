package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var ExecutePipeline = func(jobs ...job) {
	in := make(chan interface{})
	wg := &sync.WaitGroup{}

	for _, job := range jobs {
		out := make(chan interface{})
		wg.Add(1)
		go worker(job, in, out, wg)
		in = out
	}
	wg.Wait()
}

// Обертка для запуска job, чтобы дождаться выполнения всех горутин
func worker(job job, in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out)
	job(in, out)
}

// Воркер для SingleHash. По аналогии с job
func workerSingleHash(in int, out chan interface{}, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	data := strconv.Itoa(in)
	mu.Lock() // чтобы не было "OverheatLock happend"
	dataMd5 := DataSignerMd5(data)
	mu.Unlock()
	dataCrc32Md5 := DataSignerCrc32(dataMd5)
	dataCrc32 := DataSignerCrc32(data)
	fmt.Printf("%s SingleHash data %s\n", data, data)
	fmt.Printf("%s SingleHash md5(data) %s\n", data, dataMd5)
	fmt.Printf("%s SingleHash crc32(md5(data)) %s\n", data, dataCrc32Md5)
	fmt.Printf("%s SingleHash crc32(data) %s\n", data, dataCrc32)
	fmt.Printf("%s SingleHash result %s\n", data, dataCrc32+"~"+dataCrc32Md5)
	out <- dataCrc32 + "~" + dataCrc32Md5
}

var SingleHash = func(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	for i := range in {
		wg.Add(1)
		go workerSingleHash(i.(int), out, mu, wg)
	}
	wg.Wait()
}

var MultiHash = func(in, out chan interface{}) {
	for i := range in {
		var sum = ""
		for j := 0; j < 6; j++ {
			dataCrc32 := DataSignerCrc32(strconv.Itoa(j) + i.(string))
			fmt.Printf("%s MultiHash: crc32(th+step1)) %d %s\n", i.(string), j, dataCrc32)
			sum += dataCrc32
		}
		fmt.Printf("%s MultiHash result:\n%s\n", i.(string), sum)
		out <- sum
	}
}

var CombineResults = func(in, out chan interface{}) {
	var array []string
	for i := range in {
		array = append(array, i.(string))
	}
	sort.Strings(array)
	result := strings.Join(array, "_")
	fmt.Printf("CombineResults \n%s\n", result)
	out <- result
}
