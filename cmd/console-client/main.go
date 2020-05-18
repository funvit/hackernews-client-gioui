package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/peterhellberg/hn"

	log "github.com/sirupsen/logrus"
)

type indexItem struct {
	Index int
	Item  *hn.Item
}

var (
	items    = map[int]*hn.Item{}
	messages = make(chan indexItem)
)

const itemsLimit = 10

func main() {
	log.SetLevel(log.DebugLevel)
	defer func() {
		log.Println("exit")
	}()

	cli := hn.NewClient(&http.Client{
		Timeout: time.Duration(5 * time.Second),
	})

	ids, err := cli.TopStories()
	if err != nil {
		log.Panicf("stories fetch error: %s", err)

	}
	log.Debugf("got stories: %d", len(ids))

	go func() {
		for i := range messages {
			items[i.Index] = i.Item
		}
	}()

	var wg sync.WaitGroup

	for i, id := range ids[:itemsLimit] {
		wg.Add(1)
		go func(i, id int) {
			defer wg.Done()

			item, err := cli.Item(id)
			if err != nil {
				panic(fmt.Errorf("item fetch error: %w", err))
			}

			log.Debugf("item %d fetched", item.ID)
			messages <- indexItem{i, item}
		}(i, id)
	}

	wg.Wait()

	for i := 0; i < itemsLimit; i++ {
		item, ok := items[i]
		if !ok {
			log.Warnf("item %d not fetched", i)
			fmt.Println(i, "-", "NONE")
			continue
		}
		fmt.Println(i, "–", item.Title, "\n   ", item.URL, "\n")
	}
}
