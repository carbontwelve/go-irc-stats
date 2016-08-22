package main

import (
	"flag"
	"fmt"
	"github.com/carbontwelve/go-irc-stats/ircstats"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	Version string
	Build   string

	logReader ircstats.IrcLogReader

	version    = flag.Bool("version", false, "Display executable version and build.")
	verbose    = flag.Bool("v", false, "Display actual output")
	configPath = flag.String("c", "config.yaml", "Path to config.yaml")
	cwd        = flag.String("d", "", "change to this directory before doing anything")
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

	if *version {
		fmt.Printf("Version: %s Build %s\n", Version, Build)
		os.Exit(0)
	}

	// Should we be noisy?
	// If false, all stdout messages will be muted
	if *verbose {
		fmt.Println("Being Loud!")
	}

	// Should we change the current working directory? (This is super useful when testing)
	if *cwd != "" {
		err := os.Chdir(*cwd)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Load Configuration
	config := ircstats.Config{}
	configErr := config.Load(*configPath)
	if configErr != nil {
		log.Fatal(configErr)
	}

	// Load Data Store (I call it Database, but its actually a binary file)
	db := ircstats.Database{}
	db.Load(config.DatabaseLocation)

	fmt.Println("Last Parsed: ", db.LastGenerated)
	logReader = *ircstats.NewIrcLogReader(config)

	// Load log file and parse any new lines
	logReaderErr := logReader.Load(config.Location, &db)
	if logReaderErr != nil {
		log.Fatal(logReaderErr)
	}

	// Save database to disk
	dbSaveErr := db.Save(config.DatabaseLocation)
	if logReaderErr != nil {
		log.Fatal(dbSaveErr)
	}

	vd := *ircstats.NewViewData(config)
	vd.Calculate(db);
	vd.JsonData.Debug()

	j, _ := vd.GetJsonString()
	fmt.Printf("%s\n", j)

	//
	// Generate the template
	//
	v := *ircstats.NewView()
	err := v.Parse("template.html", vd)
	if err != nil {
		panic(err)
	}
}
