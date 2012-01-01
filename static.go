package main

import (
	"net/http"
	"path/filepath"
)

const themeDir = "themes/arch"

func init() {
	// register static resources.
	http.HandleFunc("/theme/", staticServe)
}

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
