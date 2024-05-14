package service

import (
	"github.com/vonuki/asyncloader/internal/collector"
	"github.com/vonuki/asyncloader/internal/streamer"
	"log"
	"math/rand"
	"sync"
	"time"
)

type (
	Service interface {
		Run()
	}

	DefaultService struct {
		log                      *log.Logger
		streamer                 streamer.Streamer
		collector                collector.Collector
		maxConcurrentConnections int
	}
)

func NewService(log *log.Logger, streamer streamer.Streamer, collector collector.Collector, maxConcurrentConnections int) *DefaultService {
	return &DefaultService{
		log:                      log,
		streamer:                 streamer,
		collector:                collector,
		maxConcurrentConnections: maxConcurrentConnections,
	}
}

func (r *DefaultService) DummyStorage(msg *collector.StatData, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	r.log.Println("Recording stat data: ", msg.Data)
	delay := rand.Intn(100)
	time.Sleep(time.Duration(delay) * time.Millisecond)
	r.log.Println("Recorded stat data: ", msg.Data)
}

func (r *DefaultService) SaveStat(dataChan <-chan *collector.StatData, doneChan chan<- bool) {
	wg := sync.WaitGroup{}
	for {
		msg, ok := <-dataChan
		if !ok {
			doneChan <- true
			break
		}
		go r.DummyStorage(msg, &wg)
	}
	wg.Wait()
}

func (r *DefaultService) Run() {
	urlsChan := make(chan string, r.maxConcurrentConnections)
	go r.streamer.RunStreamer(urlsChan)

	dataChan := make(chan *collector.StatData, r.maxConcurrentConnections)
	go r.collector.RunCollector(urlsChan, dataChan)

	done := make(chan bool, r.maxConcurrentConnections)
	go r.SaveStat(dataChan, done)
	<-done
}
