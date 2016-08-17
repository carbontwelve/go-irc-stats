package ircstats

type MaxWeek struct {
	Week  int64
	Lines int64
}

type MaxDay struct {
	Day   string
	Lines int64
}

type MaxHour struct {
	Hour  int64
	Lines int64
}

type Averages struct {
	Hour    float64			// Average lines by Hour (24 hours)
	Day     float64 		// Average lines by day (365 days)
	Week    float64 		// Average lines by week (week 1 - week 52)
	WeekDay float64 		// Average lines by week day (monday - sunday)
}

//
// Channel Statistics
//
type Channel struct {
	Name        string   		// Channel Name
	ActiveUsers map[string]User	// Active users in past 30 days
	UserCount   int64    		// Total Number of users in Channel
	LineCount   int64    		// Total Number of lines in Channel
	WordCount   int64    		// Total Word count for Channel
	MaxDay      MaxDay   		// Calculated Max Day
	MaxHour     MaxHour  		// Calculated Max Hour
	MaxWeek     MaxWeek  		// Calculated Max Week
	Averages    Averages 		// Calculated Averages
	HoursAndDaysStats    		// Inherited Hour And Days methods
	Seen                 		// First & Last Seen
}

