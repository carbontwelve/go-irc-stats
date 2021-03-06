package main

import (
	"time"
	"math"
	"strconv"
)

type SvgGraphItem struct {
	X      int64
	Value  int64
	Height int64
}

type SvgGraphLabel struct {
	X     int64
	Label string
}

type SvgGraphDay struct {
	X     int64
	Y     int64
	Date  string
	Class string
	Lines int64
}

type SvgGraphWeek struct {
	X      int64
	Y      int64
	Height int64
	Lines  int64
	First  string
	Last   string
}

type SvgGraphData struct {
	Days            []SvgGraphDay
	Weeks           []SvgGraphWeek
	Labels          []SvgGraphLabel
	MLables         []SvgGraphLabel
	MostActiveHours []SvgGraphItem
	MostActiveDays  [7]SvgGraphItem
	Width           int64
}

type ViewData struct {
	PageTitle       string
	PageDescription string
	HeatMapInterval uint
	HeatMapKey      [6]int
	Database        Database
	SvgGraphData    SvgGraphData
	WeeksMax        uint
	WeekDayMax      int64
}

func (d ViewData) TotalDays() int64 {
	return helpers.DaysDiffUnix(d.Database.Channel.Last, d.Database.Channel.First)
}

func (d *ViewData) buildDayHeatMapDays() () {
	timeNow := time.Now()
	totalDays := d.TotalDays()
	Days := make([]SvgGraphDay, totalDays)
	Weeks := make([]SvgGraphWeek, (totalDays / 7) + 2) // there is a n+1 error where Weeks starts at 0 by the 0 element is never filled, thus the +2
	Labels := make([]SvgGraphLabel, 1)
	MLables := make([]SvgGraphLabel, 1)

	// Create heatmap key
	for i := 1; i < 6; i ++ {
		d.HeatMapKey[i] = int(d.HeatMapInterval) * i
	}

	var (
		weekDays [7]SvgGraphItem
		firstWeek string
		lastWeek string
		x int64
		y int64
		mx int64
		weekLines int64
		lines int64
		cssClass string
		i int64
	)

	for i = 0; i < int64(len(weekDays)); i++ {
		weekDays[i].X = i
	}

	for i = 0; i < totalDays; i++ {
		elementTime := timeNow.AddDate(0, 0, int(-(totalDays - i)))

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
			lines = int64(d.Database.Days[elementTime.Format("2006-02-01")])
		} else {
			lines = 0
		}

		weekLines += int64(lines)
		lastWeek = elementTime.Format("Jan-01")
		weekDays[elementTime.Weekday()].Value += lines

		Weeks[x] = SvgGraphWeek{
			X: x,
			Y: y,
			Lines: weekLines,
			First: firstWeek,
			Last: lastWeek,
		}

		// Identify class
		classSet := false
		for i := 1; i < 6; i ++ {
			if int(lines) < d.HeatMapKey[i] {
				cssClass = "scale-" + strconv.Itoa(i)
				classSet = true
				break
			}
		}
		if classSet == false {
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
			Labels = append(Labels, SvgGraphLabel{
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

		//fmt.Printf("%d days ago [%s] is element %d\n", (totalDays - i), elementTime.Format("2006-02-01"), i)
	}
	d.SvgGraphData = SvgGraphData{
		Days: Days,
		Weeks: Weeks,
		Labels: Labels,
		MLables: MLables,
		MostActiveDays: weekDays, // This is a preliminary pass because its more efficient to do it here.
	}

	d.SvgGraphData.Width = (d.SvgGraphData.Days[len(d.SvgGraphData.Days) - 1].X * 10) + 10
	return
}

func (d *ViewData) buildWeekGraph() {
	// Get week max
	for _, w := range (d.SvgGraphData.Weeks) {
		if uint(w.Lines) > uint(d.WeeksMax) {
			d.WeeksMax = uint(w.Lines)
		}
	}

	// Get Weeks.Height
	tmpWeeks := make([]SvgGraphWeek, len(d.SvgGraphData.Weeks))
	for k, w := range (d.SvgGraphData.Weeks) {
		w.Height = int64(math.Floor(float64(w.Lines) / float64(d.WeeksMax) * 100))
		tmpWeeks[k] = w
	}
	d.SvgGraphData.Weeks = tmpWeeks

	// Get Most Active Times
	tmpMostActiveTimes := make([]SvgGraphItem, len(d.Database.Hours))
	for hour, lines := range d.Database.Hours {
		tmpMostActiveTimes[hour] = SvgGraphItem{
			X: int64(hour * 10),
			Value: lines,
			Height: int64(math.Floor(float64(lines) / float64(d.Database.Channel.MaxHour.Lines) * 100)),
		}
	}
	d.SvgGraphData.MostActiveHours = tmpMostActiveTimes

	// Get weekday max
	for _, v := range(d.SvgGraphData.MostActiveDays) {
		if v.Value > d.WeekDayMax {
			d.WeekDayMax = v.Value
		}
	}

	// Get Most Active Days of Week
	var tmpMostActiveDays [7]SvgGraphItem
	for day, obj := range d.SvgGraphData.MostActiveDays {
		tmpMostActiveDays[day] = SvgGraphItem{
			X: obj.X * 10,
			Height: int64(math.Floor(float64(obj.Value) / float64(d.WeekDayMax) * 100)),
		}
	}
	d.SvgGraphData.MostActiveDays = tmpMostActiveDays

	// Get week mean

	// Get week days max
}
