package crawler

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"golang.org/x/net/html"
)

/*
*	This is simple crawl function, fetches the given url and return the links in that url page.
 */
func Crawl(url string) ([]string, error) {
	//fmt.Println("fetching data for url", url)
	resp, err := http.Get(url)
	//fmt.Println("fetched data for url", url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	//fmt.Println("visiting link in url webpage", url)
	links := visit(nil, doc)
	//fmt.Println("visiting link in url webpage", url)
	return links, nil
}


// visit appends to links each link found in n, and returns the result.
func visit(links []string, n *html.Node) []string {

	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}


func CrawlDummy(url string) ([]string, error) {
	fmt.Println("fetching URL",url)
	fmt.Println("no of go routine now:", runtime.NumGoroutine())
	switch url {
	case "http://google.com":
		time.Sleep(10 * time.Millisecond)
		return []string{"deepak.com","arhat.com", "abc.com","xyz.com","icandoit.com", "youaregreate.com", "youhavedoneit.com","justdoit.com","deepakisgreat.com"},nil
	case "deepak.com":
		time.Sleep(25 * time.Millisecond)
		return []string{"deepak.com","arhat.com","letsdoit.com"},nil
	case "arhat.com":
		time.Sleep(10 * time.Millisecond)
		return []string{"deepak.com","arhat.com"},nil
	case "letsdoit.com":
		time.Sleep(15 * time.Millisecond)
		return []string{"deepak.com","vaibhavi.com"},nil
	case "youaregreate.com":
		time.Sleep(5 * time.Millisecond)
		return []string{"deepak.com","icandoit.com"},nil
	}

	return nil, errors.New("not found")
}