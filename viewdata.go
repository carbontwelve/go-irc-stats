package main

import (
	"github.com/carbontwelve/go-irc-stats/helpers"
	"fmt"
	"time"
)

type SvgGraphLabel struct {
	X     int64
	Label string
}

type SvgGraphDay struct {
	X     int64
	Y     int64
	Date  string
	Class string
	Lines uint
}

type SvgGraphWeek struct {
	X     int64
	Y     int64
	Lines int64
	First string
	Last  string
}

type ViewData struct {
	PageTitle       string
	PageDescription string
	HeatMapInterval uint
	Database        Database
	DayHeatMapDays  []SvgGraphDay
}

func (d ViewData) TotalDays() int {
	return helpers.DaysDiffUnix(d.Database.Channel.Last, d.Database.Channel.First)
}

func (d ViewData) buildDayHeatMapDays() (Days []SvgGraphDay, Weeks []SvgGraphWeek, Labels []SvgGraphLabel, MLables []SvgGraphLabel) {
	timeNow := time.Now()
	totalDays := d.TotalDays()
	Days = make([]SvgGraphDay, totalDays)
	Weeks = make([]SvgGraphWeek, (totalDays / 7) + 1)
	Labels = make([]SvgGraphLabel, 1)
	MLables = make([]SvgGraphLabel, 1)

	var (
		weekDays [7]int
		firstWeek string
		lastWeek string
		x int64
		y int64
		mx int64
		weekLines int64
		lines uint
		cssClass string
	)

	for i := 0; i < totalDays; i++ {
		elementTime := timeNow.AddDate(0, 0, -(totalDays - i))

		// Work out first week
		if (i == 0) {
			firstWeek = elementTime.Format("Jan-01")
		}

		y = int64(elementTime.Weekday())

		// If the day is Sunday
		if (y == 0) {
			x += 1
			weekLines = 0
			firstWeek = elementTime.Format("Jan-01")
		}

		// If this is the first day of the month
		if (elementTime.Day() == 1) {
			mx ++
		}

		if d.Database.HasDay(elementTime.Format("2006-02-01")) {
			lines = d.Database.Days[elementTime.Format("2006-02-01")]
		} else {
			lines = 0
		}

		weekLines += int64(lines)
		lastWeek = elementTime.Format("Jan-01")
		weekDays[elementTime.Weekday()] += int(lines)

		Weeks[x] = SvgGraphWeek{
			X: x,
			Y: y,
			Lines: weekLines,
			First: firstWeek,
			Last: lastWeek,
		}

		if (lines < d.HeatMapInterval * 5) {
			cssClass = "scale-5"
		} else if (lines < d.HeatMapInterval * 4) {
			cssClass = "scale-4"
		} else if (lines < d.HeatMapInterval * 3) {
			cssClass = "scale-3"
		} else if (lines < d.HeatMapInterval * 2) {
			cssClass = "scale-2"
		} else if (lines < d.HeatMapInterval) {
			cssClass = "scale-1"
		} else {
			cssClass = "scale-6"
		}

		Days[i] = SvgGraphDay{
			X: x,
			Y: y,
			Date: elementTime.Format("2006-02-01"),
			Class: cssClass,
			Lines: lines,
		}

		// April, July, October
		if elementTime.YearDay() == 92 || elementTime.YearDay() == 193 || elementTime.YearDay() == 274 {
			Labels = append(Labels,SvgGraphLabel{
				X: x,
				Label: elementTime.Format("Jan"),
			})
			MLables = append(MLables, SvgGraphLabel{
				X: mx,
				Label: elementTime.Format("Jan"),
			})
		}

		// New Year
		if elementTime.YearDay() == 1 {
			Labels = append(Labels, SvgGraphLabel{
				X: x,
				Label: elementTime.Format("2006"),
			})
			MLables = append(MLables, SvgGraphLabel{
				X: mx,
				Label: elementTime.Format("2006"),
			})
		}

		fmt.Printf("%d days ago [%s] is element %d\n", (totalDays - i), elementTime.Format("2006-02-01"), i)
	}

	return
}
