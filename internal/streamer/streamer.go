package streamer

import (
	"github.com/vonuki/asyncloader/config"
	"log"
	"strconv"
)

type Streamer interface {
	RunStreamer(out chan<- string)
}

type URLStreamer struct {
	log *log.Logger
	cfg *config.Config
}

func NewStreamer(log *log.Logger, cfg *config.Config) *URLStreamer {
	return &URLStreamer{
		log: log,
		cfg: cfg,
	}
}

func (r *URLStreamer) RunStreamer(out chan<- string) {
	for i := 0; i < r.cfg.URLsCount; i++ {
		// Example of DB URLs reading from config or other
		url := "exmaple_postgres://db_" + strconv.Itoa(i)
		r.log.Println("New target:", url)

		out <- url
	}
	close(out)
}
