package collector

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"
)

type (
	Collector interface {
		RunCollector(in <-chan string, out chan<- *StatData)
	}

	DefaultCollector struct {
		log                      *log.Logger
		maxConcurrentConnections int
		retriesCounter           int
	}
)

func NewCollector(log *log.Logger, maxConcurrentConnections int) *DefaultCollector {
	return &DefaultCollector{
		log:                      log,
		maxConcurrentConnections: maxConcurrentConnections,
		retriesCounter:           0,
	}
}

// DummyDBConnection - Some DB reader, handler
func (r *DefaultCollector) DummyDBConnection(ctx context.Context, db string) (*StatData, error) {
	delay := rand.Intn(100)
	time.Sleep(time.Duration(delay) * time.Millisecond)
	r.log.Println("Got stat from DB", db)
	return &StatData{
		Db:   db,
		Data: "Data from " + db,
	}, nil
}

// LoadItem - read Static data from DB
func (r *DefaultCollector) LoadItem(item string, wg *sync.WaitGroup, out chan<- *StatData) {
	defer wg.Done()

	// Add tracer into context
	ctx := context.WithValue(context.Background(), "some_value", item)

	// Send request
	statData, err := r.DummyDBConnection(ctx, item)
	if err != nil {
		//Pass
		// Improve: Add Retries
		return
	}
	// Improve: Add Validation of response
	r.log.Println("Statistics collected", item)

	out <- statData
}

// RunCollector  replies to the outgoing channel
func (r *DefaultCollector) RunCollector(dbUrls <-chan string, out chan<- *StatData) {

	//Work group to sync all Collectors
	wg := sync.WaitGroup{}
	for item := range dbUrls {
		wg.Add(1)
		go r.LoadItem(item, &wg, out)
	}
	wg.Wait()

	// all data gathered
	close(out)
}
