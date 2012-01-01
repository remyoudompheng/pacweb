package main

import (
	"flag"
	alpm "github.com/remyoudompheng/go-alpm"
	"log"
	"net/http"
	"os"
)

var logger = log.New(os.Stderr, "pacweb ", log.LstdFlags|log.Lshortfile)

func getAlpm() (*alpm.Handle, error) {
	h, er := alpm.Init("/", "/var/lib/pacman")
	if er != nil {
		return nil, er
	}

	// TODO: read /etc/pacman.conf
	h.RegisterSyncDb("core", 0)
	h.RegisterSyncDb("community", 0)
	h.RegisterSyncDb("extra", 0)
	h.RegisterSyncDb("multilib", 0)
	h.RegisterSyncDb("testing", 0)
	h.RegisterSyncDb("multilib-testing", 0)
	h.RegisterSyncDb("community-testing", 0)
	return h, er
}

func main() {
	listen := flag.String("http", "localhost:8070", "Address to listen on")
	flag.Parse()
	logger.Printf("starting HTTP server at %s", *listen)
	http.ListenAndServe(*listen, nil)
}
