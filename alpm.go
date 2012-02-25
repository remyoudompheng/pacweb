package main

import (
	"bytes"
	"errors"
	alpm "github.com/remyoudompheng/go-alpm"
	"io"
	"net/http"
	"os"
	"sort"
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

type pkgPerBuildDate []alpm.Package

func (l pkgPerBuildDate) Len() int           { return len(l) }
func (l pkgPerBuildDate) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l pkgPerBuildDate) Less(i, j int) bool { return l[i].BuildDate().After(l[j].BuildDate()) }

var _ sort.Interface = pkgPerBuildDate{}

func latestPackages() []alpm.Package {
	alpmHandleLock.RLock()
	defer alpmHandleLock.RUnlock()
	h := getAlpm()
	pkgs := make([]alpm.Package, 0)
	syncs, err := h.SyncDbs()
	if err != nil {
		return nil
	}
	syncs.ForEach(func(db alpm.Db) error {
		pkgs = append(pkgs, db.PkgCache().Slice()...)
		return nil
	})
	sort.Sort(pkgPerBuildDate(pkgs))
	if len(pkgs) > 20 {
		pkgs = pkgs[:20]
	}
	return pkgs
}

// outdatedPackages returns {pkgname: {localver, syncver, repo}}
func outdatedPackages() map[string][3]string {
	alpmHandleLock.RLock()
	defer alpmHandleLock.RUnlock()
	h := getAlpm()
	localdb, err := h.LocalDb()
	if err != nil {
		return nil
	}
	syncs, err := h.SyncDbs()
	if err != nil {
		return nil
	}
	result := make(map[string][3]string, 16)
	locals := make(map[string]string, 100)
	syncvers := make(map[string][2]string, 100)
	localdb.PkgCache().ForEach(func(pkg alpm.Package) error {
		locals[pkg.Name()] = pkg.Version()
		return nil
	})
	for _, db := range syncs.Slice() {
		db.PkgCache().ForEach(func(pkg alpm.Package) error {
			syncver, ok := syncvers[pkg.Name()]
			if !ok || alpm.VerCmp(syncver[0], pkg.Version()) < 0 {
				syncvers[pkg.Name()] = [2]string{pkg.Version(), pkg.DB().Name()}
			}
			return nil
		})
	}
	for pkgname, localver := range locals {
		if v, ok := syncvers[pkgname]; ok && alpm.VerCmp(localver, v[0]) < 0 {
			result[pkgname] = [3]string{localver, v[0], v[1]}
		}
	}
	return result
}

// HandleHome displays the homepage.
func HandleHome(resp http.ResponseWriter, req *http.Request) {
	logRequest(req)
	if req.URL.Path != "/" {
		resp.WriteHeader(http.StatusNotFound)
		ErrorPage(resp, CommonData{}, http.StatusNotFound, NoSuchPage)
		return
	}
	page, er := Execute("homepage", CommonData{}, map[string]interface{}{
		"Latest":   latestPackages(),
		"Outdated": outdatedPackages(),
	})
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
