# Web Crawler

The simple web crawler function. It accepts the `url` and `dept` as params.  Based on `url` provided, it fetches the webpage and collects
all the urls on that page. The `dept` params used to decide no of recursion for each url. 
Ex. if `google.com` is provided as `url` and `dept` is set to 2, Then this function will  collect all urls from the `google.com` and it follows the same producer to all the `urls` found on `google.com`. If `dept` was set to 3, it would have further crawled all urls. 

## Simple Solution
In the Simple directory, web crawler function have been implemented using go concurrency. It spins go routine for every urls. 
This approach works well when there are small number of urls on the page and dept is low. As the dept increases and there are many url on pages then it will spin lot-of go routines. although it takes very less time to switch the context between go routines but still huge number of go routines decreases the performance or may result into memory leak( mostly because of the task it is doing ie I/0 or Memory or CPU intensive task). The web crawler have i/o and memory(RAM) intensive operations. In production or ideal scenario we might want to limit the number of such operations due to hardware limit.

## Semaphore
To address above issue, semaphore pattern have been used here. The objective is to keep the I/O and memory operation count in control i.e. limit the I/O and memory operation running concurrently. It have a buffer channel of `n` where `n` is the number of I/O and memory operations. Semaphore version have fixed races conditions found on simple solution version as well.

Note: We are not trying to limit the number of go routines. Our program still can spin as many go routine as needed. The idea is to use the semaphore when we are ready to use it. not when we expect to use it. spawning a go routing doesn't guarantee that CPU will schedule it immediately. we let the go routine fight them self when they are ready to use it.