package UnixDiff

import (
	"time"
)

func lastDayOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 12, 31, 0, 0, 0, 0, t.Location())
}

func firstDayOfNextYear(t time.Time) time.Time {
	return time.Date(t.Year()+1, 1, 1, 0, 0, 0, 0, t.Location())
}

// Days diff for time.Time
// a - b in days
func DaysDiff(a, b time.Time) (days int) {
	cur := b
	for cur.Year() < a.Year() {
		// add 1 to count the last day of the year too.
		days += lastDayOfYear(cur).YearDay() - cur.YearDay() + 1
		cur = firstDayOfNextYear(cur)
	}
	days += a.YearDay() - cur.YearDay()
	if b.AddDate(0, 0, days).After(a) {
		days -= 1
	}
	return days
}

// Days diff for unix time
// a - b in days
func DaysDiffUnix(unixTimeA int64, unixTimeB int64) (days int) {
	return DaysDiff(time.Unix(unixTimeA, 0), time.Unix(unixTimeB, 0))
}