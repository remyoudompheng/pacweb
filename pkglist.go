package main

import (
	"fmt"
	alpm "github.com/remyoudompheng/go-alpm"
	"net/http"
)

func init() {
	http.HandleFunc("/pkglist", HandlePkglist)
}

func getLocalPackages() ([]alpm.Package, error) {
	h, er := getAlpm()
	if er != nil {
		return nil, er
	}
	db, er := h.LocalDb()
	if er != nil {
		return nil, er
	}
	return db.PkgCache().Slice(), nil
}

func getRepoPackages(repo string) ([]alpm.Package, error) {
	h, er := getAlpm()
	if er != nil {
		return nil, er
	}
	db, er := h.RegisterSyncDb(repo, 0)
	if er != nil {
		return nil, er
	}
	return db.PkgCache().Slice(), nil
}

// HandlePkglist displays a list of packages, either from local DB or a sync DB.
func HandlePkglist(resp http.ResponseWriter, req *http.Request) {
	logger.Printf("%s %s", req.Method, req.URL)
	er := req.ParseForm()
	if er != nil {
		ErrorPage(resp, CommonData{}, http.StatusInternalServerError, er)
		return
	}

	reponame := req.Form.Get("repo")
	switch reponame {
	case "", "local":
		// Local packages.
		pkglist, er := getLocalPackages()
		if er == nil {
			er = Execute(resp, "pkglist", CommonData{}, map[string]interface{}{
				"Title":    "Installed packages",
				"Repo":     "local",
				"Packages": pkglist,
			})
		}
	default:
		// Remote packages.
		pkglist, er := getRepoPackages(reponame)
		if er == nil {
			er = Execute(resp, "pkglist", CommonData{}, map[string]interface{}{
				"Title":    fmt.Sprintf("Packages from repository [%s]", reponame),
				"Repo":     reponame,
				"Packages": pkglist,
			})
		}
	}
	if er != nil {
		ErrorPage(resp, CommonData{}, http.StatusInternalServerError, er)
	}
}
