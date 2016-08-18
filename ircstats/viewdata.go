package ircstats

import "encoding/json"

//
// This view data struct will contain all the data that will be injected into the view template. Ideally this will be
// done as a JSON export so that JavaScript within the view can transform it in any way it sees fit.
//
type ViewData struct {
	PageTitle       string   // Page title from configuration
	PageDescription string   // Page description from configuration
	JsonData        JsonData // Json data for exporting to page
}

type JsonData struct {
	HeatMapInterval uint // HeatMap Interval from configuration
}

func NewViewData(c Config) *ViewData {

	j := JsonData{
		HeatMapInterval: c.HeatMapInterval,
	}

	return &ViewData{
		PageTitle: c.PageTitle,
		PageDescription: c.PageDescription,
		JsonData: j,
	}
}

func (vd ViewData) GetJsonString() (j []byte, err error) {
	j, err = json.Marshal(vd.JsonData)
	return
}