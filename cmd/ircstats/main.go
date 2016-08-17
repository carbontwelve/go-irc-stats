package main

import (
	"github.com/carbontwelve/go-irc-stats/ircstats"
	"flag"
	"os"
	"fmt"
	"log"
)

var (
	Version string
	Build string
	version = flag.Bool("version", false, "Display executable version and build.")
	configPath = flag.String("c", "config.yaml", "Path to config.yaml")
	cwd = flag.String("d", "", "change to this directory before doing anything")
	logReader ircstats.IrcLogReader
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `ircstats is a golang port from ruby of 0x263b/Stats
Usage:
        ircstats [options] ...
The options are:
  -h    Display this help text`)
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Parse command line flags
	flag.Parse()

	if (*version) {
		fmt.Printf("Version: %s Build %s\n", Version, Build)
		os.Exit(0);
	}

	if *cwd != "" {
		err := os.Chdir(*cwd)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Load Configuration
	config := ircstats.Config{}
	configErr := config.Load(*configPath)
	if (configErr != nil) {
		log.Fatal(configErr)
	}

	db := ircstats.Database{}
	db.Load(config.DatabaseLocation)

	fmt.Println("Last Parsed: ", db.LastGenerated)
	logReader = ircstats.NewIrcLogReader(config);

	// Load log file and parse any new lines
	logReaderErr := logReader.Load(config.Location, &db);
	if (logReaderErr != nil) {
		log.Fatal(logReaderErr)
	}

	// @todo calculate stats

	// Save database to disk
	dbSaveErr := db.Save(config.DatabaseLocation)
	if (logReaderErr != nil) {
		log.Fatal(dbSaveErr)
	}

//
	//// Get Database to calculate stats and totals
	//lr.Database.Calculate()
//
	//fmt.Printf("Last line date [%d]\n", lr.Database.Channel.Last)
	//fmt.Printf("Mean Lines/Day: %f\n", lr.Database.Channel.Mean)
	//fmt.Printf("Total Lines Parsed: %d\n", lr.Database.Channel.LineCount)
//
	//// Once we are finished dump to disk cache file
	//lr.Database.Save(config.DatabaseLocation)
//
	//vd := ViewData{
	//	PageTitle : config.PageTitle,
	//	PageDescription : config.PageDescription,
	//	HeatMapInterval: config.HeatMapInterval,
	//	Database : lr.Database,
	//}
//
	//vd.buildDayHeatMapDays()
	//vd.buildWeekGraph()
//
	//v := View{}
	//err := v.Parse("template.html", vd)
//
	//if (err != nil){
	//	panic(err)
	//}
}
