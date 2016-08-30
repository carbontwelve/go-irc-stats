package main

import (
	"flag"
	"fmt"
	"github.com/carbontwelve/go-irc-stats/ircstats"
	"log"
	"os"
	"path/filepath"
)

var (
	Version string
	Build string

	logReader ircstats.IrcLogReader

	version = flag.Bool("version", false, "Display executable version and build.")
	cleardb = flag.Bool("del", false, "Delete database and re-parse from the beginning")
	verbose = flag.Bool("v", false, "Display actual output")
	configPath = flag.String("c", "config.yaml", "Path to config.yaml")
	cwd = flag.String("d", "", "change to this directory before doing anything")
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
		logError(os.Chdir(*cwd))
	}

	// Load Configuration
	config := ircstats.Config{}
	logError(config.Load(*configPath))

	// Load Data Store (I call it Database, but its actually a binary file)
	db := ircstats.Database{}
	if *cleardb {
		logError(os.Remove(config.DatabaseLocation))
	}
	db.Load(config.DatabaseLocation)

	fmt.Println("Last Parsed: ", db.LastGenerated)
	logReader = *ircstats.NewIrcLogReader(config)

	// Load log file and parse any new lines
	// If config.Location is a directory, identify all files within and load one after the other
	var logReaderErr error
	if isDirectory(config.Location) {
		logReaderErr = filepath.Walk(config.Location, func(path string, f os.FileInfo, err error) error {
			if isDirectory(path) == true {
				return nil;
			}
			return logReader.Load(path, &db)
		});
	} else {
		logReaderErr = logReader.Load(config.Location, &db)
	}

	logError(logReaderErr)

	//
	// Save database to disk
	//
	logError(db.Save(config.DatabaseLocation))

	vd := *ircstats.NewViewData(config)
	vd.Calculate(db)
	vd.JsonData.Debug()

	//
	// Generate the template
	//
	v := *ircstats.NewView()
	logError(v.Parse("template.html", vd))
}

// Returns true if path is a directory
func isDirectory(path string) bool {
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return true;
	}
	return false;
}

// Logs an error if err is not nil
func logError(err error) {
	if err != nil {
		if *verbose {
			panic(err)
		} else {
			log.Fatal(err)
		}
	}
}
