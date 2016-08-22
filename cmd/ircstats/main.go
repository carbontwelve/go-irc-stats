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

	// Calculate
	calculate(&db, config)

	// Save database to disk
	dbSaveErr := db.Save(config.DatabaseLocation)
	if logReaderErr != nil {
		log.Fatal(dbSaveErr)
	}

	// Lol "stats"
	fmt.Printf("%d number of new lines parsed.\n\n", logReader.NewLinesFound())
	fmt.Println("= [ Dates ] ======================================")
	fmt.Printf("First line date\t\t\t\t\t%d\n", db.Channel.FirstSeen)
	fmt.Printf("Last line date\t\t\t\t\t%d\n", db.Channel.LastSeen)
	fmt.Printf("Total Days Seen:\t\t\t\t%d\n", db.Channel.TotalDaysSeen())
	fmt.Println("= [ Averages ] ===================================")
	fmt.Printf("Mean Lines/Hr:\t\t\t\t\t%f\n", db.Channel.Averages.Hour)
	fmt.Printf("Mean Lines/Day:\t\t\t\t\t%f\n", db.Channel.Averages.Day)
	fmt.Printf("Mean Lines/Week:\t\t\t\t%f\n", db.Channel.Averages.Week)
	fmt.Printf("Mean Lines/Week Day:\t\t\t%f\n", db.Channel.Averages.WeekDay)
	fmt.Println("= [ Counters ] ===================================")
	fmt.Printf("Total Lines Parsed:\t\t\t\t%d\n", db.Channel.LineCount)
	fmt.Printf("Total Unique Users:\t\t\t\t%d\n", db.Channel.UserCount)
	fmt.Printf("Users Active in past 30 days:\t%d\n", len(db.Channel.ActiveUsers))

	fmt.Printf("Peak Day Date:\t\t\t\t\t%s\n", db.Channel.MaxDay.Day)
	fmt.Printf("Peak Day Lines:\t\t\t\t\t%d\n", db.Channel.MaxDay.Lines)

	fmt.Printf("Peak Hour:\t\t\t\t\t\t%d\n", db.Channel.MaxHour.Hour)
	fmt.Printf("Peak Hour Lines:\t\t\t\t%d\n", db.Channel.MaxHour.Lines)

	fmt.Printf("Peak Week:\t\t\t\t\t\t%d\n", db.Channel.MaxWeek.Week)
	fmt.Printf("Peak Week Lines:\t\t\t\t%d\n", db.Channel.MaxWeek.Lines)

	fmt.Println("==================================================")

	vd := *ircstats.NewViewData(config)
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

func calculate(db *ircstats.Database, config ircstats.Config) {
	// Calculate UserCount
	db.Channel.UserCount = int64(len(db.Users))
	db.Channel.MaxDay.Day, db.Channel.MaxDay.Lines = db.Channel.FindPeakDay()
	db.Channel.MaxHour.Hour, db.Channel.MaxHour.Lines = db.Channel.FindPeakHour()
	db.Channel.MaxWeek.Week, db.Channel.MaxWeek.Lines = db.Channel.FindPeakWeek()

	db.Channel.Averages.Hour = db.Channel.FindHourAverage()
	db.Channel.Averages.Week = db.Channel.FindWeekAverage()
	db.Channel.Averages.Day = db.Channel.FindDayAverage()

	// Calculate Active Users
	calculateActiveUsers(db)

	calculateHeatMapDays(db, config.HeatMapInterval)
}

func calculateActiveUsers(db *ircstats.Database) {
	timePeriod := make(map[string]bool)
	db.Channel.ActiveUsers = make(map[string]ircstats.User)

	for i := 1; i < 30; i++ {
		timePeriod[time.Now().AddDate(0, 0, -i).Format("2006-02-01")] = true
	}

	for _, u := range db.Users {
		var (
			wordCount  int64
			daysActive int64
		)

		// Check to see if user has been active within our time period (default past 30 days)
		for timePeriodDate := range timePeriod {
			if _, ok := u.Days[timePeriodDate]; ok {
				daysActive++
				wordCount += u.Days[timePeriodDate]
			}
		}

		// If the user is active, copy the struct and push it to the active users hash
		if daysActive > 0 {
			uc := u
			uc.WordCount = wordCount
			uc.DaysActive = daysActive
			uc.WordsDay = wordCount / daysActive
			db.Channel.ActiveUsers[uc.Username] = uc
		}
	}
}

type WeekValue struct {
	X     int64
	Value int64
	First string
	Last  string
}

type DayValue struct {
	X     int64
	Y     int64
	Date  string
	Class string
	Value int64
}

func calculateHeatMapDays(db *ircstats.Database, heatMapInterval uint) {

	var (
		heatMapKey [6]int
		now        time.Time
		totalDays  int
		firstWeek  string
		lastWeek   string
		weekValues []WeekValue
		dayValues  []DayValue
		x          int64
		mx         int64
		y          int64
		weekLines  int64
		dayLines   int64
		cssClass   string
	)

	now = time.Now()
	totalDays = int(db.Channel.TotalDaysSeen())

	// Create heatmap key
	for i := 1; i < 6; i++ {
		heatMapKey[i] = int(heatMapInterval) * i
	}

	// Build Days Data
	for i := 0; i < totalDays; i++ {
		elementTime := now.AddDate(0, 0, int(-(totalDays - i)))

		// Work out first week
		if i == 0 {
			firstWeek = elementTime.Format("Jan-01")
		}

		// Y is bound to week day (Sunday = 0)
		y = int64(elementTime.Weekday())

		// If day is a Sunday then begin new week
		if y == 0 {
			x += 1
			weekLines = 0
			firstWeek = elementTime.Format("Jan-01")
		}

		// If this is the first day of the month
		if elementTime.Day() == 1 {
			mx++
		}

		if db.Channel.HasDay(elementTime.Format("2006-02-01")) {
			dayLines = int64(db.Channel.Days[elementTime.Format("2006-02-01")])
		} else {
			dayLines = 0
		}

		weekLines += dayLines

		// If day is Saturday then end week
		if y == 6 {
			lastWeek = elementTime.Format("Jan-01")
			weekValues = append(weekValues, WeekValue{X: x, Value: weekLines, First: firstWeek, Last: lastWeek})
		}

		// Identify class
		classSet := false
		for i := 1; i < 6; i++ {
			if int(dayLines) < heatMapKey[i] {
				cssClass = "scale-" + strconv.Itoa(i)
				classSet = true
				break
			}
		}
		if classSet == false {
			cssClass = "scale-6"
		}

		// Append days value
		dayValues = append(dayValues, DayValue{X: x, Y: y, Date: elementTime.Format("2006-02-01"), Class: cssClass, Value: dayLines})
	}
}
