package main

import (
	"github.com/vonuki/asyncloader/config"
	"github.com/vonuki/asyncloader/internal/collector"
	"github.com/vonuki/asyncloader/internal/service"
	"github.com/vonuki/asyncloader/internal/streamer"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "Stat loader ", log.LstdFlags)
	logger.Println("Starting loader..")
	cfg := config.LoadConfig("config/config.yml")

	// Initialize Streamer
	// It returns list of DB urls
	urlStreamer := streamer.NewStreamer(log.New(os.Stdout, "UrlStreamer ", log.LstdFlags), cfg)
	logger.Println("urlStreamer created")

	// Initialize Collector
	// Read Statistic from DBs
	coll := collector.NewCollector(log.New(os.Stdout, "Collector ", log.LstdFlags), cfg.Concurrency)
	logger.Println("Collector created")

	// Initialize Service
	// Write statistic
	srv := service.NewService(log.New(os.Stdout, "Service ", log.LstdFlags), urlStreamer, coll, cfg.Concurrency)
	logger.Println("Service initialized")

	// Run Streamer, Collector inside service
	srv.Run()

	logger.Println("Loader stopping")
	os.Exit(0)
}
