package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	Fetch(url string) (body string, urls []string, err error)
}

// Структура для хранения запрошенных ранее URL
type Cache struct {
	urls map[string]bool
	mux  sync.Mutex
}

func Crawl(url string, depth int, fetcher Fetcher) {
	cacheUrl := Cache{}
	cacheUrl.urls = make(map[string]bool)

	var wg sync.WaitGroup
	var crawlRecursive func(string, int)

	crawlRecursive = func(url string, depth int) {
		defer wg.Done()
		if depth <= 0 {
			return
		}
		cacheUrl.mux.Lock()
		if ok := cacheUrl.urls[url]; ok {
			cacheUrl.mux.Unlock()
			return
		} else {
			cacheUrl.urls[url] = true
			cacheUrl.mux.Unlock()
		}
		body, urls, err := fetcher.Fetch(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("found: %s %q\n", url, body)
		// Перебор всех url полученных от Fetch
		for _, u := range urls {
			wg.Add(1)
			// Распараллеливание задачи поиска
			go crawlRecursive(u, depth-1)
		}
	}
	wg.Add(1)
	crawlRecursive(url, depth)
	wg.Wait()
}

func main() {
	Crawl("http://golang.org/", 4, fetcher)
}

type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
