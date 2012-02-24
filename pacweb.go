// pacweb is a Web interface for ArchLinux package management.
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
)

var logger = log.New(os.Stderr, "pacweb ", log.LstdFlags|log.Lshortfile)

func logRequest(req *http.Request) {
	logger.Printf("%s %s from %s", req.Method, req.URL, req.RemoteAddr)
}

var (
	profile = flag.String("pprof", "", "Dump CPU profile at that path")
)

func sigHandle() {
	for sig := range signal.Incoming {
		logger.Printf("caught signal %v", sig)
		switch sig {
		case os.SIGTERM, os.SIGINT:
			if *profile != "" {
				pprof.StopCPUProfile()
			}
			logger.Printf("exiting.")
			os.Exit(1)
		}
	}
}

func main() {
	listen := flag.String("http", "localhost:8070", "Address to listen on")
	flag.Parse()

	if *profile != "" {
		f, er := os.Create(*profile)
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
	http.ListenAndServe(*listen, nil)
}
