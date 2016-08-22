package ircstats

import (
	"encoding/json"
	"fmt"
	"time"
)

//
// This view data struct will contain all the data that will be injected into the view template. Ideally this will be
// done as a JSON export so that JavaScript within the view can transform it in any way it sees fit.
//
type ViewData struct {
	PageTitle       string   // Page title from configuration
	PageDescription string   // Page description from configuration
	JsonData        JsonData // Json data for exporting to page
}

type UserData struct {
	Username           string
	Url                string
	Avatar             string
	FirstSpoke         int64
	LastSpoke          int64
	TotalWords         int64            // Count of words
	Averages           Averages         // Used for words/day
	Vocabulary         int64            // Number of different words used
	Words              map[string]int64 // World cloud
	DaysActiveInPeriod int64
	TotalWordsInPeriod int64
	WordsByDayInPeriod float64
	ActivityPercentage float64 // Overall % contribution to Channel.WordCount
}

type JsonData struct {
	// Configurable options
	HeatMapInterval uint // HeatMap Interval from configuration
	ActivityPeriod  uint // Activity Period from configuration

	// Dates
	FirstSeen     int64 // Timestamp of first message
	LastSeen      int64 // Timestamp of last message
	TotalDaysSeen int64 // Number of days between FirstSeen and LastSeen

	// Averages
	Averages Averages // Calculated Averages

	// Counters
	MaxDay           MaxDay  // Calculated Max Day
	MaxHour          MaxHour // Calculated Max Hour
	MaxWeek          MaxWeek // Calculated Max Week
	TotalLines       int64   // Lines parsed in total
	TotalWords       int64   // Total Words (all words multiplied by times used)
	TotalUsers       int64   // Number of unique users
	TotalActiveUsers int64   // Number of active users within activity period (default 30 days)

	// Misc
	Users       map[string]UserData // Users
	ActiveUsers map[string]int64
}

func NewViewData(c Config) *ViewData {

	j := JsonData{
		HeatMapInterval: c.HeatMapInterval,
		ActivityPeriod:  c.ActivityPeriod,
	}

	return &ViewData{
		PageTitle:       c.PageTitle,
		PageDescription: c.PageDescription,
		JsonData:        j,
	}
}

func (j JsonData) Debug() {
	fmt.Println("==================================================")
	fmt.Println("Json Data Debug:")
	fmt.Println("= [ Dates ] ======================================")
	fmt.Printf("First line date\t\t\t\t\t%d\n", j.FirstSeen)
	fmt.Printf("Last line date\t\t\t\t\t%d\n", j.LastSeen)
	fmt.Printf("Total Days Seen:\t\t\t\t%d\n", j.TotalDaysSeen)
	fmt.Println("= [ Averages ] ===================================")
	fmt.Printf("Mean Lines/Hr:\t\t\t\t\t%f\n", j.Averages.Hour)
	fmt.Printf("Mean Lines/Day:\t\t\t\t\t%f\n", j.Averages.Day)
	fmt.Printf("Mean Lines/Week:\t\t\t\t%f\n", j.Averages.Week)
	fmt.Printf("Mean Lines/Week Day:\t\t\t%f\n", j.Averages.WeekDay)
	fmt.Println("= [ Counters ] ===================================")
	fmt.Printf("Total Lines Parsed:\t\t\t\t%d\n", j.TotalLines)
	fmt.Printf("Total Unique Users:\t\t\t\t%d\n", j.TotalUsers)
	fmt.Printf("Users Active in past %d days:\t%d\n", j.ActivityPeriod, j.TotalActiveUsers)

	fmt.Printf("Peak Day Date:\t\t\t\t\t%s\n", j.MaxDay.Day)
	fmt.Printf("Peak Day Lines:\t\t\t\t\t%d\n", j.MaxDay.Lines)

	fmt.Printf("Peak Hour:\t\t\t\t\t\t%d\n", j.MaxHour.Hour)
	fmt.Printf("Peak Hour Lines:\t\t\t\t%d\n", j.MaxHour.Lines)

	fmt.Printf("Peak Week:\t\t\t\t\t\t%d\n", j.MaxWeek.Week)
	fmt.Printf("Peak Week Lines:\t\t\t\t%d\n", j.MaxWeek.Lines)

	fmt.Println("==================================================")
}

// Calculate stats for View
func (vd *ViewData) Calculate(db Database) {
	// Dates
	vd.JsonData.FirstSeen = db.Channel.FirstSeen
	vd.JsonData.LastSeen = db.Channel.LastSeen
	vd.JsonData.TotalDaysSeen = db.Channel.TotalDaysSeen()

	// Calculate Counters
	vd.JsonData.TotalUsers = db.CountUsers()
	vd.JsonData.MaxDay.Day, vd.JsonData.MaxDay.Lines = db.Channel.FindPeakDay()
	vd.JsonData.MaxHour.Hour, vd.JsonData.MaxHour.Lines = db.Channel.FindPeakHour()
	vd.JsonData.MaxWeek.Week, vd.JsonData.MaxWeek.Lines = db.Channel.FindPeakWeek()
	vd.JsonData.TotalLines = db.Channel.LineCount
	vd.JsonData.TotalWords = db.Channel.WordCount

	// Calculate Averages
	vd.JsonData.Averages.Hour = db.Channel.FindHourAverage()
	vd.JsonData.Averages.Week = db.Channel.FindWeekAverage()
	vd.JsonData.Averages.Day = db.Channel.FindDayAverage()

	// Calculate Users
	vd.calculateUsers(db)
}

func (vd *ViewData) calculateUsers(db Database) {
	var (
		timePeriod map[string]bool
		users      map[string]UserData

		userWordCount  int64
		userDaysActive int64
	)

	timePeriod = make(map[string]bool)
	vd.JsonData.ActiveUsers = make(map[string]int64)
	users = make(map[string]UserData)

	for i := 1; i < int(vd.JsonData.ActivityPeriod); i++ {
		timePeriod[time.Now().AddDate(0, 0, -i).Format("2006-02-01")] = true
	}

	for nick, u := range db.Users {
		userWordCount = 0
		userDaysActive = 0

		for timePeriodDate := range timePeriod {
			if _, ok := u.Days[timePeriodDate]; ok {
				userDaysActive++
				userWordCount += u.Days[timePeriodDate]
			}
		}

		viewUserData := UserData{
			Username:           u.Username,
			Url:                u.Url,
			Avatar:             u.Avatar,
			FirstSpoke:         u.FirstSeen,
			LastSpoke:          u.LastSeen,
			TotalWords:         u.WordCount,
			Vocabulary:         int64(len(u.Words)),
			Words:              u.Words,
			DaysActiveInPeriod: userDaysActive,
			TotalWordsInPeriod: userWordCount,
			ActivityPercentage: (float64(u.WordCount) / float64(db.Channel.WordCount)) * 100,
		}

		viewUserData.Averages.Hour = u.FindHourAverage()
		viewUserData.Averages.Week = u.FindWeekAverage()
		viewUserData.Averages.Day = u.FindDayAverage()

		if userDaysActive > 0 {
			viewUserData.WordsByDayInPeriod = float64(userWordCount) / float64(userDaysActive)
			vd.JsonData.ActiveUsers[nick] = userDaysActive
		}
		users[nick] = viewUserData
	}

	vd.JsonData.Users = users
}

func (vd ViewData) GetJsonString() (j []byte, err error) {
	j, err = json.Marshal(vd.JsonData)
	return
}
