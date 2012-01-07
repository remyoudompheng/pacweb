package main

import (
	"errors"
	alpm "github.com/remyoudompheng/go-alpm"
	"net/http"
)

func init() {
	http.HandleFunc("/", HandleHome)
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

// forallFilenames iterates call() over all (relative) file paths
// from packages in db.
func forallFilenames(db alpm.Db, call func(string) error) error {
	return db.PkgCache().ForEach(func(pkg alpm.Package) error {
		for _, file := range pkg.Files() {
			er := call(file.Name)
			if er != nil {
				return er
			}
		}
		return nil
	})
}

var NoSuchPage = errors.New("undefined page")

// HandleHome displays the homepage.
func HandleHome(resp http.ResponseWriter, req *http.Request) {
	logger.Printf("%s %s", req.Method, req.URL)
	if req.URL.Path != "/" {
		resp.WriteHeader(http.StatusNotFound)
		ErrorPage(resp, CommonData{}, http.StatusNotFound, NoSuchPage)
		return
	}
	Execute(resp, "homepage", CommonData{}, nil)
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
