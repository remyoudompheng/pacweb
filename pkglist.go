package main

import (
	"fmt"
	alpm "github.com/remyoudompheng/go-alpm"
	"net/http"
)

func init() {
	http.HandleFunc("/pkglist", HandlePkglist)
}

func getPackages() ([]alpm.Package, error) {
	h, er := alpm.Init("/", "/var/lib/pacman")
	if er != nil {
		return nil, er
	}
	db, er := h.LocalDb()
	if er != nil {
		return nil, er
	}
	return db.PkgCache().Slice(), nil
}

func HandlePkglist(resp http.ResponseWriter, req *http.Request) {
	pkglist, _ := getPackages()
	er := Execute(resp, "pkglist", CommonData{}, map[string]interface{}{
		"Title":    "Hello",
		"Packages": pkglist,
	})
	if er != nil {
		resp.WriteHeader(500)
		fmt.Fprintf(resp, "internal error: %s\n", er)
	}
}
