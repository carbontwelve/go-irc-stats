package main

import (
	"html/template"
	"io/ioutil"
	"bytes"
)

type View struct {}

func (v View) Parse(filename string, data ViewData) (err error) {
	funcMap := template.FuncMap{
		"tableflip": func () string { return "(╯°□°）╯︵ ┻━┻" },
	}
	template, err := template.New(filename).Funcs(funcMap).ParseFiles(filename)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	err = template.Execute(buf, data)
	ioutil.WriteFile("stats.html", buf.Bytes(), 0600)
	return
}