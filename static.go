package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

const staticDir = "static/"

func init() {
	// register static resources.
	http.HandleFunc("/static/", staticServe)
	// serve local files.
	http.HandleFunc("/file/", serveLocalFile)
}

func staticServe(resp http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	relpath, er := filepath.Rel("/static", path)

	if er != nil {
		resp.WriteHeader(http.StatusNotFound)
		er = ErrorPage(resp, CommonData{}, http.StatusNotFound, er)
		if er != nil {
			logger.Printf("error: %s", er)
		}
		return
	}

	logger.Printf("%s %s -> %s", req.Method, req.URL, filepath.Join(staticDir, relpath))
	http.ServeFile(resp, req, filepath.Join(staticDir, relpath))
}

var (
	errSpecialFile    = errors.New("not a regular file")
	errNonPackageFile = errors.New("not from a package")
)

// isPackageFilepath decides whether a filename is from a package
// (or only probably from a package.
func isPackageFilepath(name string) bool {
	return true
}

func serveLocalFile(resp http.ResponseWriter, req *http.Request) {
	logger.Printf("%s %s", req.Method, req.URL)
	er := req.ParseForm()
	if er != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		ErrorPage(resp, CommonData{}, http.StatusInternalServerError, er)
		return
	}
	filename := req.Form.Get("path")

	// file path sanity checks.
	switch {
	case !filepath.IsAbs(filename),
		filepath.HasPrefix(filename, "/etc"),
		filepath.HasPrefix(filename, "/home"),
		filepath.HasPrefix(filename, "/var"),
		filepath.HasPrefix(filename, "/tmp"):
		resp.WriteHeader(http.StatusForbidden)
		er := fmt.Errorf("access to %s is forbidden", filename)
		ErrorPage(resp, CommonData{}, http.StatusForbidden, er)
		return
	}

	info, er := os.Stat(filename)
	switch {
	case er != nil:
		resp.WriteHeader(http.StatusNotFound)
		ErrorPage(resp, CommonData{}, http.StatusNotFound, er)
		return
	case info.Mode()&os.ModeType != 0:
		// not a regular file.
		resp.WriteHeader(http.StatusForbidden)
		ErrorPage(resp, CommonData{}, http.StatusForbidden, errSpecialFile)
		return
	case !isPackageFilepath(filename):
		// not a package file
		resp.WriteHeader(http.StatusForbidden)
		ErrorPage(resp, CommonData{}, http.StatusForbidden, errNonPackageFile)
		return
	}

	http.ServeFile(resp, req, filename)
}
