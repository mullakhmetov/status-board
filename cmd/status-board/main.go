package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mullakhmetov/status-board/internal/rest"
)

func main() {
	var port, timeout, askRate int
	var metrics bool
	var sitesPath string

	flag.IntVar(&port, "port", 8080, "server listen port")
	flag.IntVar(&timeout, "timeout", 5, "service ask timeout in seconds")
	flag.IntVar(&askRate, "check_rate", 60, "service cheks rate in seconds")
	flag.BoolVar(&metrics, "metrics", false, "enable metrics")
	flag.StringVar(&sitesPath, "sites_path", "", "abs path to sites file")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// catch signal and invoke graceful termination
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Printf("[WARN] interrupt signal")
		cancel()
	}()

	server, err := rest.NewServer(rest.ServerOpts{
		Port:         port,
		Timeout:      time.Second * time.Duration(timeout),
		ChecksRate:   time.Second * time.Duration(askRate),
		StoreMetrics: metrics,
		SitesPath:    sitesPath,
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	err = server.Run(ctx)
	if err != nil {
		log.Printf("[ERROR] terminated with error %+v", err)
		os.Exit(1)
	}

	log.Printf("[INFO] terminated")
}
