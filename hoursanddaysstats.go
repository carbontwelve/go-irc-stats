package main

type HoursAndDaysStats struct {
	Hours [23]uint        // lines per 24 hours
	Days  map[string]uint // words per day
}

// Increment Hours by one
func (s *HoursAndDaysStats) IncrementHour(hour uint) {
	s.Hours[hour]++
}

// Increment Days by an input number
func (s *HoursAndDaysStats) IncrementDay(date string, increment uint) {
	if _, ok := s.Days[date]; ok {
		s.Days[date] += increment
	} else {
		s.Days[date] = 1
	}
}

func (s HoursAndDaysStats) FindPeakDay() (date string, total uint) {
	for d, t := range (s.Days) {
		if (t > total) {
			date = d
			total = t
		}
	}
	return
}

func (s *HoursAndDaysStats) Initiate() {
	for i := 0; i < 23; i++ {
		s.Hours[i] = 0
	}
	s.Days = make(map[string]uint)
}
