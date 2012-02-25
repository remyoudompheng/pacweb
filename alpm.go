package main

import (
	"bytes"
	"errors"
	alpm "github.com/remyoudompheng/go-alpm"
	"io"
	"net/http"
	"os"
	"sync"
)

func init() {
	http.HandleFunc("/", HandleHome)
	http.HandleFunc("/repolist", HandleRepolist)
	http.HandleFunc("/info", HandlePkgInfo)
}

var (
	alpmHandle     *alpm.Handle
	alpmHandleFlag = new(sync.Once)
	alpmHandleLock = new(sync.RWMutex)
)

func initAlpm() {
	fconf, err := os.Open("/etc/pacman.conf")
	if err != nil {
		panic(err)
	}
	defer fconf.Close()
	conf, err := alpm.ParseConfig(fconf)
	if err != nil {
		panic(err)
	}
	h, err := conf.CreateHandle()
	if err != nil {
		panic(err)
	}
	alpmHandle = h
}

func getAlpm() *alpm.Handle {
	if alpmHandle == nil {
		alpmHandleFlag.Do(initAlpm)
	}
	return alpmHandle
}

func getRepoList() ([]byte, error) {
	alpmHandleLock.RLock()
	defer alpmHandleLock.RUnlock()
	h := getAlpm()
	dblist, er := h.SyncDbs()
	if er != nil {
		return nil, er
	}

	s, er := Execute("repolist", CommonData{}, map[string]interface{}{
		"Repos": dblist.Slice(),
	})
	return s, er
}

// forallFilenames iterates call() over all (relative) file paths
// from packages in db.
func forallFilenames(db *alpm.Db, call func(string) error) error {
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
	logRequest(req)
	if req.URL.Path != "/" {
		resp.WriteHeader(http.StatusNotFound)
		ErrorPage(resp, CommonData{}, http.StatusNotFound, NoSuchPage)
		return
	}
	page, er := Execute("homepage", CommonData{}, nil)
	if er != nil {
		ErrorPage(resp, CommonData{}, http.StatusInternalServerError, er)
	} else {
		io.Copy(resp, bytes.NewBuffer(page))
	}

}

// HandleRepolist displays basic information about available sync DBs.
func HandleRepolist(resp http.ResponseWriter, req *http.Request) {
	logRequest(req)
	page, er := getRepoList()
	if er != nil {
		ErrorPage(resp, CommonData{}, http.StatusInternalServerError, er)
	} else {
		io.Copy(resp, bytes.NewBuffer(page))
	}
}
