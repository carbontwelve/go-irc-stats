package main

import (
	"fmt"
	"regexp"
)

var (
	Version string
	Build string
)

func main() {
	fmt.Printf("Version: %s Build %s\n", Version, Build)

	// Load Configuration
	config := Config{}
	configErr := config.Load("config.yaml")
	if (configErr != nil) {
		panic(configErr)
	}

	db := Database{}
	db.Load(config.DatabaseLocation)

	fmt.Println("Last Parsed: ", db.LastGenerated)

	lr := LogReader{
		RegexAction: regexp.MustCompile(`^\[(.+)\] \* (.+)$`),
		RegexMessage: regexp.MustCompile(`^\[(.+)\] <(.+)> (.+)$`),
		RegexParseAction: regexp.MustCompile(`^\[(.+)\] \* (\S+) (.+)$`),
		RegexParseMessage: regexp.MustCompile(`^\[(.+)\] <(\S+)> (.+)$`),
		Database: db,
		Config: config,
	}

	// Load log file and parse any new lines
	lr.LoadFile()

	// Get Database to calculate stats and totals
	lr.Database.Calculate()

	fmt.Printf("Last line date [%d]\n", lr.Database.Channel.Last)
	fmt.Printf("Mean Lines/Day: %f\n", lr.Database.Channel.Mean)
	fmt.Printf("Total Lines Parsed: %d\n", lr.Database.Channel.LineCount)

	// Once we are finished dump to disk cache file
	lr.Database.Save(config.DatabaseLocation)

	vd := ViewData{
		PageTitle : config.PageTitle,
		PageDescription : config.PageDescription,
		Channel : lr.Database.Channel,
		ActiveUsers : lr.Database.ActiveUsers,
	}
	v := View{}
	v.Load("template.html")
	v.Parse(vd)

	//fmt.Printf("%v\n", lr.Database)
}
