package main

import (
	"flag"
	"fmt"
	"net/http"
	"path/filepath"
)

const themeDir = "themes/arch"

func staticServe(resp http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	relpath, er := filepath.Rel("/theme", path)

	if er != nil {
		resp.WriteHeader(404)
		fmt.Fprintf(resp, "Error 404: %s\n", er)
		return
	}

	http.ServeFile(resp, req, filepath.Join(themeDir, relpath))
}

func init() {
	// register static resources.
	http.HandleFunc("/theme/", staticServe)
}

func main() {
	listen := flag.String("http", "localhost:8070", "Address to listen on")
	flag.Parse()
	http.ListenAndServe(*listen, nil)
}
