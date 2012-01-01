package main

import (
	"errors"
	"fmt"
	alpm "github.com/remyoudompheng/go-alpm"
	"io"
	"net/http"
	"reflect"
	"runtime/debug"
	"text/template"
)

var pacwebTemplate *template.Template

// Parity is used to create alternating styles in tables.
func Parity(x int) string {
	switch x % 2 {
	case 0:
		return "even"
	case 1:
		return "odd"
	}
	panic("plouf")
}

func IsInstalled(p *alpm.Package) bool {
	return p.DB().Name() == "local"
}

func SplitSlice(x interface{}) (y []interface{}, er error) {
	t := reflect.ValueOf(x)
	switch t.Kind() {
	case reflect.Array, reflect.Slice:
		half := (t.Len() + 1) / 2
		first, second := t.Slice(0, half), t.Slice(half, t.Len())
		return []interface{}{first, second}, nil
	}
	return nil, errors.New("not a slice")
}

func init() {
	// parse templates.
	t := template.New("root")
	t.Funcs(template.FuncMap{
		"parity":         Parity,
		"httpStatusText": http.StatusText,
		"isInstalled":    IsInstalled,
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

func Execute(w io.Writer, tplName string, common CommonData, contents interface{}) error {
	return pacwebTemplate.ExecuteTemplate(w, tplName, TplInput{common, contents})
}

type ErrorData struct {
	StatusCode int
	Error      error
	Stack      string
}

func ErrorPage(w io.Writer, common CommonData, code int, err error) error {
	contents := ErrorData{code, err, string(debug.Stack())}
	logger.Printf("error: %s", err)
	return pacwebTemplate.ExecuteTemplate(w, "error", TplInput{common, contents})
}

func PanicPage(w io.Writer, panicData interface{}) {
	switch x := panicData.(type) {
	case error:
		ErrorPage(w, CommonData{}, http.StatusInternalServerError, x)
	default:
		ErrorPage(w, CommonData{}, http.StatusInternalServerError, fmt.Errorf("%+v", x))
	}
}
