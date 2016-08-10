package main

import (
	"html/template"
	"io/ioutil"
	"bytes"
)

type View struct {
	template *template.Template
	functionMapping template.FuncMap
}
type ViewData struct {
	PageTitle string
	PageDescription string
	Channel Channel
	ActiveUsers map[string]User
}

func (v *View) Load(filename string) (err error) {
	v.template, err = template.ParseFiles(filename)
	if err != nil {
		return err
	}

	v.functionMapping = template.FuncMap{
		"eq": func(a, b interface{}) bool {
			return a == b
		},
	}

	v.template.Funcs(v.functionMapping)
	return
}

func (v View) Parse(data ViewData) (err error) {
	buf := &bytes.Buffer{}
	err = v.template.Execute(buf, data)
	ioutil.WriteFile("stats.html", buf.Bytes(), 0600)
	return
}

func (d ViewData) TotalDays() int {
	return 123
}