package ircstats

import "encoding/json"

//
// This view data struct will contain all the data that will be injected into the view template. Ideally this will be
// done as a JSON export so that JavaScript within the view can transform it in any way it sees fit.
//
type ViewData struct {
	PageTitle       string // Page title from configuration
	PageDescription string // Page description from configuration
	HeatMapInterval int64  // HeatMap Interval from configuration
}

func (d ViewData) export() (b []byte, err error) {
	b, err = json.Marshal(d);
	return
}

