package main

import (
	"bytes"
	"errors"
	"fmt"
	alpm "github.com/remyoudompheng/go-alpm"
	"html/template"
	"io"
	"net/http"
	"runtime/debug"
	"time"
)

var pacwebTemplate *template.Template

// TEMPLATE FUNCTIONS

// TimeFormat formats time according to a chosen format.
func TimeFormat(t time.Time) string { return t.Format(time.RFC1123) }

// HumanSize formats a file size for human readability.
func HumanSize(n int64) string {
	switch {
	case n > 1<<20:
		return fmt.Sprintf("%.2f MiB", float64(n)/float64(1<<20))
	default:
		return fmt.Sprintf("%.0f kiB", float64(n)/float64(1<<10))
	}
	panic("impossible")
}

func IsLocal(p *alpm.Package) bool {
	return p.DB().Name() == "local"
}

func InstallStatus(p *alpm.Package) string {
	alpmHandleLock.RLock()
	defer alpmHandleLock.RUnlock()
	localdb, err := getAlpm().LocalDb()
	if err != nil {
		return "could not found local DB"
	}
	localp, err := localdb.PkgByName(p.Name())
	if err == nil && localp != nil {
		switch cmp := alpm.VerCmp(p.Version(), localp.Version()); {
		case cmp > 0:
			return "Upgradable"
		case cmp == 0:
			return "Installed"
		case cmp < 0:
			return "Local version is newer"
		}
	}
	return "Not installed"
}

func init() {
	// parse templates.
	t := template.New("root")
	t.Funcs(template.FuncMap{
		"timeFmt":        TimeFormat,
		"httpStatusText": http.StatusText,
		"isLocal":        IsLocal,
		"installStatus":  InstallStatus,
		"humanSize":      HumanSize,
	})
	pacwebTemplate = template.Must(t.ParseGlob("templates/*.tpl"))
}

type TplInput struct {
	Common   CommonData
	Contents interface{}
}

type CommonData struct {
	SysMessage string
}

func Execute(tplName string, common CommonData, contents interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := pacwebTemplate.ExecuteTemplate(buf, tplName, TplInput{common, contents})
	if err != nil {
		logger.Printf("template error: %s", err)
	}
	return buf.Bytes(), err
}

type ErrorData struct {
	StatusCode int
	Error      error
	Stack      string
}

func ErrorPage(w io.Writer, common CommonData, code int, err error) error {
	if err == nil {
		err = nilError
	}
	contents := ErrorData{code, err, string(debug.Stack())}
	logger.Printf("error: %v", err)
	err = pacwebTemplate.ExecuteTemplate(w, "error", TplInput{common, contents})
	if err != nil {
		logger.Printf("template error: %s", err)
	}
	return err
}

var nilError = errors.New("nil")

func PanicPage(w io.Writer, panicData interface{}) {
	switch x := panicData.(type) {
	case error:
		if x == nil {
			x = nilError
		}
		ErrorPage(w, CommonData{}, http.StatusInternalServerError, x)
	default:
		ErrorPage(w, CommonData{}, http.StatusInternalServerError, fmt.Errorf("%+v", x))
	}
}
