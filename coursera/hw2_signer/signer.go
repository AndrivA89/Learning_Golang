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

	for _, jobItem := range jobs {
		out := make(chan interface{})
		wg.Add(1)
		go func(jobItem job, in, out chan interface{}, wg *sync.WaitGroup) {
			defer wg.Done()
			defer close(out)
			jobItem(in, out)
		}(jobItem, in, out, wg)
		in = out
	}
	wg.Wait()
}

// Воркер для SingleHash. По аналогии с job
func workerSingleHash(in int, out chan interface{}, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	data := strconv.Itoa(in)
	dataHashChan := make(chan string)
	mu.Lock() // чтобы не было "OverheatLock happend"
	dataMd5 := DataSignerMd5(data)
	mu.Unlock()
	go func(dataHash string, outChan chan string) {
		outChan <- DataSignerCrc32(dataHash)
	}(data, dataHashChan)
	dataCrc32Md5 := DataSignerCrc32(dataMd5)
	dataCrc32 := <-dataHashChan
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
	wg := &sync.WaitGroup{}
	for dataIn := range in {
		wg.Add(1)
		go func(data string) {
			wgCrc32 := &sync.WaitGroup{}
			dataCrc32 := make([]string, 6)
			for i := 0; i < 6; i++ {
				wgCrc32.Add(1)
				go func(i int, dataCrc32 []string, wgCrc32 *sync.WaitGroup, data string) {
					dataCrc32[i] = DataSignerCrc32(strconv.Itoa(i) + data)
					wgCrc32.Done()
				}(i, dataCrc32, wgCrc32, data)
			}
			wgCrc32.Wait()
			out <- strings.Join(dataCrc32, "")
			wg.Done()
		}(dataIn.(string))
	}
	wg.Wait()
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
