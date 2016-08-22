package ircstats

type Seen struct {
	FirstSeen int64
	LastSeen  int64
}

func (s *Seen) UpdateSeen(seen int64) {
	if s.FirstSeen == 0 {
		s.FirstSeen = seen
	}
	if s.LastSeen == 0 {
		s.LastSeen = seen
	}
	if s.FirstSeen > seen {
		s.FirstSeen = seen
	}
	if s.LastSeen < seen {
		s.LastSeen = seen
	}
}

func (s Seen) TotalDaysSeen() int64 {
	return DaysDiffUnix(s.LastSeen, s.FirstSeen)
}

type HoursAndDaysStats struct {
	Weeks [54]int64        // lines per 52 weeks
	Hours [24]int64        // lines per 24 hours @todo check that all 24 elements are being filled...
	Days  map[string]int64 // words per day
}

// Increment Hours by one
func (s *HoursAndDaysStats) IncrementHour(hour uint) {
	s.Hours[hour]++
}

// Increment Days by an input number
func (s *HoursAndDaysStats) IncrementDay(date string, increment int64) {
	if _, ok := s.Days[date]; ok {
		s.Days[date] += increment
	} else {
		s.Days[date] = 1
	}
}

// Increment Weeks by an input number
func (s *HoursAndDaysStats) IncrementWeek(week int, increment int64) {
	s.Weeks[week] += increment
}

func (s HoursAndDaysStats) HasDay(day string) bool {
	if _, ok := s.Days[day]; ok {
		return true
	}
	return false
}

func (s HoursAndDaysStats) FindPeakDay() (date string, total int64) {
	for d, t := range s.Days {
		if t > total {
			date = d
			total = t
		}
	}
	return
}

func (s HoursAndDaysStats) FindPeakHour() (hour int64, total int64) {
	for h, t := range s.Hours {
		if t > total {
			hour = int64(h)
			total = t
		}
	}
	return
}

func (s HoursAndDaysStats) FindPeakWeek() (week int64, total int64) {
	for w, t := range s.Weeks {
		if t > total {
			week = int64(w)
			total = t
		}
	}
	return
}

func (s HoursAndDaysStats) FindHourAverage() (avg float64) {
	var (
		sum  int64
		size int64
	)

	avg = 0.0

	for _, t := range s.Hours {
		sum += t
		size++
	}

	if size > 0 {
		avg = float64(sum) / float64(size)
	}

	return
}

func (s HoursAndDaysStats) FindDayAverage() (avg float64) {
	var (
		sum  int64
		size int64
	)

	avg = 0.0

	for _, t := range s.Days {
		sum += t
		size++
	}

	if size > 0 {
		avg = float64(sum) / float64(size)
	}

	return
}

func (s HoursAndDaysStats) FindWeekAverage() (avg float64) {
	var (
		sum  int64
		size int64
	)

	avg = 0.0

	for _, t := range s.Weeks {
		sum += t
		size++
	}

	if size > 0 {
		avg = float64(sum) / float64(size)
	}
	return
}

func (s HoursAndDaysStats) FindWeekDayAverage() (avg float64) {
	avg = 0.0
	return
}

func (s *HoursAndDaysStats) Initiate() {
	for i := 0; i < 24; i++ {
		s.Hours[i] = 0
	}

	for i := 0; i < 54; i++ {
		s.Weeks[i] = 0
	}
	s.Days = make(map[string]int64)
}
