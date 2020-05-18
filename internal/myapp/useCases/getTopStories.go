package useCases

import (
	"fmt"
	"github.com/peterhellberg/hn"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync/atomic"
	"time"
)

type IndexItem struct {
	Index int
	Item  *hn.Item
}

func GetTopStories(limit int32, items chan<- IndexItem, timeout time.Duration) error {
	cli := hn.NewClient(&http.Client{
		Timeout: timeout,
	})

	ids, err := cli.TopStories()
	if err != nil {
		return fmt.Errorf("stories fetch error: %s", err)
	}

	log.Debugf("got stories: %d", len(ids))

	var fetchedItems int32

	// todo: use pool with limited size!
	for i, id := range ids[:limit] {
		go func(i, id int) {
			defer atomic.AddInt32(&fetchedItems, 1)

			item, err := cli.Item(id)
			if err != nil {
				log.Errorf("item fetch error: %w", err)
				return
			}

			log.Debugf("item ID=%d fetched, idx=%d", item.ID, i)
			items <- IndexItem{i, item}
		}(i, id)
	}

	//fixme: how to close items chan?
	go func() {
		for {
			if atomic.LoadInt32(&fetchedItems) == limit {
				log.Debug("GetTopStories: closing output chan")
				close(items)
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()

	return nil
}
