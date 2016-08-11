package main

import (
	"html/template"
	"io/ioutil"
	"bytes"
	"strconv"
)

type View struct {}

func (v View) Parse(filename string, data ViewData) (err error) {
	funcMap := template.FuncMap{
		"tableflip": func () string { return "(╯°□°）╯︵ ┻━┻" },
		"mul": func (x, y int64) int64 { return x * y },
		"labelX": func (x int64) string {
			x = x * 10
			xx := float64(x)
			xx = xx + 4.5
			return strconv.FormatFloat(xx, 'f', 6, 64)
		},
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