package main

import (
	"bytes"
	alpm "github.com/remyoudompheng/go-alpm"
	"io"
	"net/http"
)

// BuildPkgInfo produces information about a package.
func BuildPkgInfo(pkgname, dbname string) (resp []byte, er error) {
	alpmHandleLock.Lock()
	defer alpmHandleLock.Unlock()

	var (
		pkg *alpm.Package
		db  *alpm.Db
	)
	h := getAlpm()
	switch dbname {
	case "", "local":
		dbname = "local"
		db, er = h.LocalDb()
	default:
		db, er = h.SyncDbByName(dbname)
	}
	if er != nil {
		return
	}

	pkg, er = db.PkgByName(pkgname)
	if er != nil {
		return
	}

	s, er := Execute("pkginfo", CommonData{}, map[string]interface{}{
		"Package": pkg,
		"Repo":    dbname,
	})
	return s, er
}

// HandlePkgInfo displays a package metadata.
func HandlePkgInfo(resp http.ResponseWriter, req *http.Request) {
	logRequest(req)
	defer func() {
		if x := recover(); x != nil {
			PanicPage(resp, x)
		}
	}()
	er := req.ParseForm()
	if er != nil {
		ErrorPage(resp, CommonData{}, http.StatusInternalServerError, er)
		return
	}

	pkgname := req.Form.Get("pkg")
	dbname := req.Form.Get("db")
	respBytes, er := BuildPkgInfo(pkgname, dbname)
	if er == nil {
		io.Copy(resp, bytes.NewBuffer(respBytes))
	} else {
		ErrorPage(resp, CommonData{}, http.StatusInternalServerError, er)
	}
}
