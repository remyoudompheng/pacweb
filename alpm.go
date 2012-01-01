package main

import (
	alpm "github.com/remyoudompheng/go-alpm"
	"net/http"
)

func init() {
	http.HandleFunc("/repolist", HandleRepolist)
	http.HandleFunc("/info", HandlePkgInfo)
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

// HandleRepolist displays basic information about available sync DBs.
func HandleRepolist(resp http.ResponseWriter, req *http.Request) {
	logger.Printf("%s %s", req.Method, req.URL)
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

// HandlePkgInfo displays a package metadata.
func HandlePkgInfo(resp http.ResponseWriter, req *http.Request) {
	defer func() {
		if x := recover(); x != nil {
			PanicPage(resp, x)
		}
	}()
	logger.Printf("%s %s", req.Method, req.URL)
	er := req.ParseForm()
	if er != nil {
		ErrorPage(resp, CommonData{}, http.StatusInternalServerError, er)
		return
	}

	pkgname := req.Form.Get("pkg")
	dbname := req.Form.Get("db")
	var (
		pkg *alpm.Package
		db  *alpm.Db
	)
	h, er := getAlpm()
	if er != nil {
		panic(er)
	}
	switch dbname {
	case "", "local":
		dbname = "local"
		db, er = h.LocalDb()
	default:
		db, er = h.RegisterSyncDb(dbname, 0)
	}
	if er != nil {
		panic(er)
	}

	pkg, er = db.GetPkg(pkgname)
	if er != nil {
		panic(er)
	}

	er = Execute(resp, "pkginfo", CommonData{}, map[string]interface{}{
		"Package": pkg,
		"Repo":    dbname,
	})
	if er != nil {
		panic(er)
	}
}
