package main

import (
	"fmt"
	"regexp"
)

var (
	Version string
	BuildTime string
)

func main() {
	fmt.Println("Version: ", Version)
	fmt.Println("Build Time: ", BuildTime)

	db := Database{}
	db.Load("db.bin")

	fmt.Println("Last Parsed: ", db.LastGenerated)

	lr := LogReader{
		RegexAction: regexp.MustCompile(`^\[(.+)\] \* (.+)$`),
		RegexMessage: regexp.MustCompile(`^\[(.+)\] <(.+)> (.+)$`),
		RegexParseAction: regexp.MustCompile(`^\[(.+)\] \* (\S+) (.+)$`),
		RegexParseMessage: regexp.MustCompile(`^\[(.+)\] <(\S+)> (.+)$`),
		Database: db,
	}

	// Load log file and parse any new lines
	lr.LoadFile("irctest.log")

	// Get Database to calculate stats and totals
	lr.Database.Calculate()

	fmt.Printf("Last line date [%d]\n", lr.Database.Channel.Last)
	fmt.Printf("Mean Lines/Day: %f\n", lr.Database.Channel.Mean)

	// Once we are finished dump to disk cache file
	lr.Database.Save("db.bin")

	//fmt.Printf("%v\n", lr.Database.Channel.MaxDay)
}
