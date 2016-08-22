package ircstats

import (
	"bytes"
	"github.com/dustin/go-humanize"
	"html/template"
	"io/ioutil"
	"strconv"
	"math"
)

type View struct {
	viewFunctions template.FuncMap
}

func NewView() *View {
	return &View{
		viewFunctions: template.FuncMap{
			"tableflip": func() string {
				return "(╯°□°）╯︵ ┻━┻"
			},
			"mul": func(x, y int64) int64 {
				return x * y
			},
			"labelX": func(x int64) string {
				x = x * 10
				xx := float64(x)
				xx = xx + 4.5
				return strconv.FormatFloat(xx, 'f', 6, 64)
			},
			"comma": humanize.Comma,
			"int64": func(x int) int64 {
				return int64(x)
			},
		}}
}

func (v View) Parse(filename string, data ViewData) (err error) {
	funcMap := template.FuncMap{
		"tableflip": func() string {
			return "(╯°□°）╯︵ ┻━┻"
		},
		"mul": func(x, y int64) int64 {
			return x * y
		},
		"labelX": func(x int64) string {
			x = x * 10
			xx := float64(x)
			xx = xx + 4.5
			return strconv.FormatFloat(xx, 'f', 6, 64)
		},
		"round": math.Floor,
		"comma": humanize.Comma,
		"int64": func(x int) int64 {
			return int64(x)
		},
		"floattoint64": func (x float64) int64 {
			return int64(x)
		},
	}
	t, err := template.New(filename).Funcs(funcMap).ParseFiles(filename)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	err = t.Execute(buf, data)
	ioutil.WriteFile("stats.html", buf.Bytes(), 0600)
	return
}
