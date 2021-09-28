package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"golang.frontdoorhome.com/personal-project/web-crawler/crawler"
)

type result struct {
	url   string
	urls  []string
	err   error
	depth int
}

type request struct {
	url  string
	dept int
}

func main() {
	now := time.Now()

	webCrawler("http://google.com", 2)
	fmt.Println("time taken:", time.Since(now))
}

func webCrawler(url string, dept int) {
	urlChan := make(chan request)
	resultChan := make(chan result)
	done := make(chan bool)
	defer close(urlChan)
	defer close(resultChan)

	go urlCrawler(urlChan, resultChan)
	go processResult(urlChan, resultChan, done)
	initialUrlRequest := request{
		url:  url,
		dept: dept,
	}
	urlChan <- initialUrlRequest
	<-done

}

func urlCrawler(urlChan <-chan request, resultChan chan<- result) {
	semChan := make(chan struct{}, 100)
	defer close(semChan)
	for {
		requestData, ok := <-urlChan
		if !ok {
			return
		}
		fmt.Println("no of go routine now:", runtime.NumGoroutine())

		go func(requestData request) {
			semChan <- struct{}{}
			fmt.Println("go routine added")
			fetch(resultChan, requestData.url, requestData.dept)
			<-semChan
		}(requestData)

	}

}

func fetch(resultChan chan<- result, url string, depth int) {
	fmt.Printf("fetching: %s\n", url)
	urls, err := crawler.Crawl(url)
	//urls, err := crawler.CrawlDummy(url)
	if err == nil {
		fmt.Printf("found: %s\n", url)
	} else {
		fmt.Printf("Error for Url %s : %v\n", url, err)
	}

	resultChan <- result{url, urls, err, depth}
}

func processResult(urlChan chan<- request, resultChan <-chan result, done chan<- bool) {
	fetched := make(map[string]bool)
	var fanOutUrlCount int64 = 1
	var fetchMutex sync.Mutex

	for result := range resultChan {
		fetchMutex.Lock()
		fetched[result.url] = true
		fetchMutex.Unlock()
		atomic.AddInt64(&fanOutUrlCount, -1)

		if result.depth > 0 {
			if result.err != nil && atomic.LoadInt64(&fanOutUrlCount) == 0 {
				fetchMutex.Lock()
				fmt.Println("Total urls crawled", len(fetched))
				fetchMutex.Unlock()
				done <- true
			}
			go func(urls []string, dept int) {
				for _, url := range urls {
					isFetched := false
					fetchMutex.Lock()
					isFetched = !fetched[url]
					fetchMutex.Unlock()
					if isFetched {
						atomic.AddInt64(&fanOutUrlCount, 1)

						request := request{
							url:  url,
							dept: dept - 1,
						}
						urlChan <- request
					}
				}
			}(result.urls, result.depth)

		} else if atomic.LoadInt64(&fanOutUrlCount) == 0 && result.depth <= 0 {
			fmt.Println("Total urls crawled", len(fetched))
			done <- true
		}
	}

}
