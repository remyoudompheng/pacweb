package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

var logger = log.New(os.Stderr, "pacweb ", log.LstdFlags|log.Lshortfile)

func logRequest(req *http.Request) {
	logger.Printf("%s %s from %s", req.Method, req.URL, req.RemoteAddr)
}

func main() {
	listen := flag.String("http", "localhost:8070", "Address to listen on")
	flag.Parse()
	logger.Printf("starting HTTP server at %s", *listen)
	http.ListenAndServe(*listen, nil)
}
