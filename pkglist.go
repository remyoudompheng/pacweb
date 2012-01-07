package main

import (
	"bytes"
	"fmt"
	alpm "github.com/remyoudompheng/go-alpm"
	"io"
	"net/http"
)

func init() {
	http.HandleFunc("/pkglist", HandlePkglist)
}

func getLocalPackages() ([]alpm.Package, error) {
	h := getAlpm()
	db, er := h.LocalDb()
	if er != nil {
		return nil, er
	}
	return db.PkgCache().Slice(), nil
}

func getRepoPackages(repo string) ([]alpm.Package, error) {
	h := getAlpm()
	db, er := h.RegisterSyncDb(repo, 0)
	if er != nil {
		return nil, er
	}
	return db.PkgCache().Slice(), nil
}

func BuildPkglist(reponame string) (resp []byte, err error) {
	alpmHandleLock.Lock()
	defer alpmHandleLock.Unlock()

	switch reponame {
	case "", "local":
		// Local packages.
		pkglist, err := getLocalPackages()
		if err != nil {
			return nil, err
		}
		s, err := Execute("pkglist", CommonData{}, map[string]interface{}{
			"Title":    "Installed packages",
			"Repo":     "local",
			"Packages": pkglist,
		})
		return s, err
	default:
		// Remote packages.
		pkglist, err := getRepoPackages(reponame)
		if err != nil {
			return nil, err
		}
		s, err := Execute("pkglist", CommonData{}, map[string]interface{}{
			"Title":    fmt.Sprintf("Packages from repository [%s]", reponame),
			"Repo":     reponame,
			"Packages": pkglist,
		})
		return s, err
	}
	panic("impossible")
}

// HandlePkglist displays a list of packages, either from local DB or a sync DB.
func HandlePkglist(resp http.ResponseWriter, req *http.Request) {
	logRequest(req)
	er := req.ParseForm()
	if er != nil {
		ErrorPage(resp, CommonData{}, http.StatusInternalServerError, er)
		return
	}

	reponame := req.Form.Get("repo")
	response, er := BuildPkglist(reponame)
	if er != nil {
		ErrorPage(resp, CommonData{}, http.StatusInternalServerError, er)
	} else {
		io.Copy(resp, bytes.NewBuffer(response))
	}
}
