package crawler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

// MyCrawler sets base URL and depth
type MyCrawler struct {
	//the base URL of the website being crawled
	BaseURL string
	Depth   int
	host    string
	//a regular expression pointer to the RegExp that will be used to extract the
	//URLs from each request.
	Rxp *regexp.Regexp
}

// CrawledMap is crawled url map
type CrawledMap struct {
	mu      sync.Mutex
	syncmap map[string]int
	urls    []string
}

var cmap *CrawledMap

func (cm *CrawledMap) set(key string) bool {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if _, ok := cm.syncmap[key]; ok {
		return false
	}
	cm.syncmap[key] = 1
	cm.urls = append(cm.urls, strings.Split(key, "://")[1])
	return true
}

func (cm *CrawledMap) getUrls() []string {
	return cm.urls
}

func (mc *MyCrawler) isValidURL() bool {
	if strings.Contains(mc.BaseURL, "http") || strings.Contains(mc.BaseURL, "www") {
		return false
	}
	mc.host = strings.Split(mc.BaseURL, "/")[0]
	mc.BaseURL = "http://" + mc.BaseURL
	_, err := url.ParseRequestURI(mc.BaseURL)
	return err == nil
}

// crawl is function to crawl the input url
func (mc *MyCrawler) crawl(URL string, depth int) {
	// fmt.Println(URL)
	// var netTransport = &http.Transport{
	// 	Dial: (&net.Dialer{
	// 		Timeout: 50 * time.Second,
	// 	}).Dial,
	// 	TLSHandshakeTimeout: 50 * time.Second,
	// }
	// var netClient = &http.Client{
	// 	Timeout:   time.Second * 50,
	// 	Transport: netTransport,
	// }
	// resp, err := netClient.Get(URL)
	resp, err := http.Get(URL)
	if depth <= 0 {
		// fmt.Println("@@@@@@@@@@@@@@")
		return
	}
	if err != nil {
		fmt.Println(err)
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Read error has occured")
		} else {
			strBody := string(body)
			// fmt.Println("$$$$$", depth, strBody)
			exURLs := mc.extractUrls(URL, strBody)
			fetched := make(chan bool)
			fmt.Println("<<<<<<<<", len(exURLs))
			for _, exURL := range exURLs {
				go func(url string) {
					if cmap.set(url) {
						//fmt.Println("$$$$$", depth)
						mc.crawl(url, depth-1)
					}
					fetched <- true
				}(exURL)
			}
			for i := range exURLs {
				_ = i
				<-fetched
			}
		}
	}

}

func (mc *MyCrawler) extractUrls(URL, body string) []string {
	var urls []string
	newURLs := mc.Rxp.FindAllStringSubmatch(body, -1)
	fmt.Println("$$$$$$$$$", len(newURLs), mc.host)
	//baseURL, _ := url.Parse(URL)
	for _, v := range newURLs {
		newURL := v[1]
		fmt.Println(mc.host, newURL)
		urltype, err := url.Parse(newURL)
		if err == nil {
			if urltype.IsAbs() == true && strings.Contains(newURL, mc.host) {
				urls = append(urls, newURL)
			} else {
				fmt.Println("skip")
			}
			// else if urltype.IsAbs() == false {
			// 	urls = append(urls, baseURL.ResolveReference(urltype).String())
			// } else if strings.HasPrefix(newURL, "//") {
			// 	urls = append(urls, "http:"+newURL)
			// } else if strings.HasPrefix(newURL, "/") {
			// 	urls = append(urls, mc.BaseURL+newURL)
			// } else {
			// 	urls = append(urls, URL+newURL)
			// }
		}
	}
	return urls
}

// Start is to start crawler
func (mc *MyCrawler) Start() ([]string, error) {
	if mc.Rxp == nil {
		mc.Rxp = regexp.MustCompile(`(?s)<a[ t]+.*?href="(http.*?)".*?>.*?</a>`)
	}
	switch {
	case mc.Depth <= 1:
		return nil, errors.New("MyCrawler Depth should be greater than 1")
	case mc.BaseURL == "":
		return nil, errors.New("MyCrawler BaseURL is empty. Please set a base url")
	case !mc.isValidURL():
		return nil, errors.New("MyCrawler BaseURL is invalid. Please set a valid base url like 'naver.com'")
	}
	cmap = &CrawledMap{
		syncmap: make(map[string]int),
	}

	mc.crawl(mc.BaseURL, mc.Depth)
	return cmap.getUrls(), nil
}
