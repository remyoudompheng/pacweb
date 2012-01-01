package main

import (
	"io"
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

func init() {
	// parse templates.
	t := template.New("root")
	t.Funcs(template.FuncMap{
		"parity": Parity,
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
