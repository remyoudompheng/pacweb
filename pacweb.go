// pacweb is a Web interface for ArchLinux package management.
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
)

var logger = log.New(os.Stderr, "pacweb ", log.LstdFlags|log.Lshortfile)

func logRequest(req *http.Request) {
	logger.Printf("%s %s from %s", req.Method, req.URL, req.RemoteAddr)
}

var profile string

func init() {
	flag.StringVar(&profile, "pprof", "", "Dump CPU profile at that path")
}

func sigHandle() {
	incoming := make(chan os.Signal, 1)
	signal.Notify(incoming, syscall.SIGTERM, syscall.SIGINT)
	for sig := range incoming {
		logger.Printf("caught signal %v", sig)
		if profile != "" {
			pprof.StopCPUProfile()
		}
		logger.Printf("exiting.")
		os.Exit(1)
	}
}

func main() {
	listen := flag.String("http", "localhost:8070", "Address to listen on")
	flag.Parse()

	if profile != "" {
		f, er := os.Create(profile)
		if er != nil {
			panic(er)
		}
		er = pprof.StartCPUProfile(f)
		if er != nil {
			panic(er)
		}
	}

	go sigHandle()
	logger.Printf("starting HTTP server at %s", *listen)
	err := http.ListenAndServe(*listen, nil)
	logger.Print(err)
}
