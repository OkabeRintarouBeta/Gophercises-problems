package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"okaberintarou.quiethn/hn"
)

type application struct {
}

type templateData struct {
	Stories []item
	Time    time.Duration
}

func main() {
	var numStories, port int
	flag.IntVar(&numStories, "num_stories", 30, "Number of stories to retrieve")
	flag.IntVar(&port, "port", 4000, "Port to run the server on")
	flag.Parse()
	tmpl := template.Must(template.ParseFiles("./templates/index.gohtml"))
	app := application{}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: app.routes(numStories, tmpl),
	}
	log.Fatal(srv.ListenAndServe())
}

func (app *application) mainHandler(numStories int, tmpl *template.Template) http.HandlerFunc {
	sc := storyCache{
		numStories: numStories,
		duration:   6 * time.Second,
	}

	go func() {
		timer := time.NewTicker(3 * time.Second)
		for {
			temp := storyCache{
				numStories: numStories,
				duration:   6 * time.Second,
			}
			temp.stories()
			sc.mutex.Lock()
			sc.cache = temp.cache
			sc.expiration = temp.expiration
			sc.mutex.Unlock()
			<-timer.C
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		stories, err := sc.stories()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		data := templateData{
			Stories: stories,
			Time:    time.Since(start),
		}
		buf := new(bytes.Buffer)
		if err = tmpl.Execute(buf, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		buf.WriteTo(w)
	})
}

func getTopStories(numStories int) ([]item, error) {
	var client hn.Client
	ids, err := client.TopItems()
	if err != nil {
		return nil, errors.New("fail to retrieve new stories")
	}
	var stories []item
	at := 0
	for len(stories) < numStories {
		need := (numStories - len(stories)) * 5 / 4
		stories = append(stories, getStories(ids[at:at+need])...)
		at += need
	}
	return stories[:numStories], nil

}

func getStories(ids []int) []item {
	type result struct {
		item item
		err  error
		idx  int
	}

	// goroutine
	resultCh := make(chan result)
	var wg sync.WaitGroup

	for i, id := range ids {

		wg.Add(1)
		go func(idx, id int) {
			defer wg.Done()
			var client hn.Client
			hnItem, err := client.GetStory(id)
			if err != nil {
				resultCh <- result{idx: idx, err: err}
			}
			resultCh <- result{item: parseHNItem(hnItem), idx: idx}

		}(i, id)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var results []result
	for i := 0; i < len(ids); i++ {
		results = append(results, <-resultCh)
	}

	// sort story by index
	sort.Slice(results, func(i, j int) bool {
		return results[i].idx < results[j].idx
	})
	var stories []item
	for _, res := range results {
		if res.err != nil {
			continue
		}
		if isStoryLink(res.item) {
			stories = append(stories, res.item)
		}
	}
	return stories
}

type storyCache struct {
	numStories int
	cache      []item
	expiration time.Time
	duration   time.Duration
	mutex      sync.Mutex
}

func (sc *storyCache) stories() ([]item, error) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	if time.Since(sc.expiration) < 0 {
		return sc.cache, nil
	}
	stories, err := getTopStories(sc.numStories)
	if err != nil {
		return nil, err
	}
	sc.expiration = time.Now().Add(sc.duration)
	sc.cache = stories
	return sc.cache, nil

}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

type item struct {
	hn.Item
	Host string
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}
