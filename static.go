package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const themeDir = "themes/arch"

var logger = log.New(os.Stderr, "pacweb ", log.LstdFlags|log.Lshortfile)

func staticServe(resp http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	relpath, er := filepath.Rel("/theme", path)

	if er != nil {
		resp.WriteHeader(http.StatusNotFound)
		er = ErrorPage(resp, CommonData{}, http.StatusNotFound, er)
		if er != nil {
			logger.Printf("error: %s", er)
		}
		return
	}

	logger.Printf("%s %s -> %s", req.Method, req.URL, filepath.Join(themeDir, relpath))
	http.ServeFile(resp, req, filepath.Join(themeDir, relpath))
}

func init() {
	// register static resources.
	http.HandleFunc("/theme/", staticServe)
}

func main() {
	listen := flag.String("http", "localhost:8070", "Address to listen on")
	flag.Parse()
	logger.Printf("starting HTTP server at %s", *listen)
	http.ListenAndServe(*listen, nil)
}
