package main

import (
	alpm "github.com/remyoudompheng/go-alpm"
	"net/http"
)

func init() {
	http.HandleFunc("/pkglist", HandlePkglist)
	http.HandleFunc("/repolist", HandleRepolist)
	http.HandleFunc("/info", HandlePkgInfo)
}

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

func getPackages() ([]alpm.Package, error) {
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

func HandlePkglist(resp http.ResponseWriter, req *http.Request) {
	pkglist, er := getPackages()
	if er != nil {
		ErrorPage(resp, CommonData{}, http.StatusInternalServerError, er)
		return
	}
	er = Execute(resp, "pkglist", CommonData{}, map[string]interface{}{
		"Title":    "Installed packages",
		"Packages": pkglist,
	})
	if er != nil {
		ErrorPage(resp, CommonData{}, http.StatusInternalServerError, er)
	}
}

func getDbList() ([]alpm.Db, error) {
	h, er := getAlpm()
	if er != nil {
		return nil, er
	}
	dblist, er := h.SyncDbs()
	if er != nil {
		return nil, er
	}
	return dblist.Slice(), nil
}

func HandleRepolist(resp http.ResponseWriter, req *http.Request) {
	dbs, er := getDbList()
	if er != nil {
		ErrorPage(resp, CommonData{}, http.StatusInternalServerError, er)
		return
	}

	er = Execute(resp, "repolist", CommonData{}, map[string]interface{}{
		"Repos": dbs,
	})
	if er != nil {
		ErrorPage(resp, CommonData{}, http.StatusInternalServerError, er)
	}
}

func HandlePkgInfo(resp http.ResponseWriter, req *http.Request) {
}
